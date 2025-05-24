package watch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	managementv1 "dev.khulnasoft.com/api/pkg/apis/management/v1"
	storagev1 "dev.khulnasoft.com/api/pkg/apis/storage/v1"
	loftclient "dev.khulnasoft.com/api/pkg/clientset/versioned"
	informers "dev.khulnasoft.com/api/pkg/informers/externalversions"
	informermanagementv1 "dev.khulnasoft.com/api/pkg/informers/externalversions/management/v1"
	"dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/platform"
	"dev.khulnasoft.com/pkg/platform/client"
	"dev.khulnasoft.com/pkg/platform/project"
	"dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

// WorkspacesCmd holds the cmd flags
type WorkspacesCmd struct {
	*flags.GlobalFlags

	Log log.Logger
}

// NewWorkspacesCmd creates a new command
func NewWorkspacesCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &WorkspacesCmd{
		GlobalFlags: globalFlags,
		Log:         log.Default.ErrorStreamOnly(),
	}
	c := &cobra.Command{
		Use:    "workspaces",
		Short:  "Watches all workspaces for a project",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run(cobraCmd.Context(), os.Stdin, os.Stdout, os.Stderr)
		},
	}

	return c
}

type ProWorkspaceInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   managementv1.DevSpaceWorkspaceInstanceSpec `json:"spec,omitempty"`
	Status ProWorkspaceInstanceStatus               `json:"status,omitempty"`
}

type ProWorkspaceInstanceStatus struct {
	managementv1.DevSpaceWorkspaceInstanceStatus `json:",inline"`

	Source *provider.WorkspaceSource    `json:"source,omitempty"`
	IDE    *provider.WorkspaceIDEConfig `json:"ide,omitempty"`
}

func (cmd *WorkspacesCmd) Run(ctx context.Context, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	if cmd.Context == "" {
		cmd.Context = config.DefaultContext
	}

	projectName := os.Getenv(provider.LOFT_PROJECT)
	if projectName == "" {
		return fmt.Errorf("project name not found")
	}

	baseClient, err := client.InitClientFromPath(ctx, cmd.Config)
	if err != nil {
		return err
	}

	managementConfig, err := baseClient.ManagementConfig()
	if err != nil {
		return err
	}

	clientset, err := loftclient.NewForConfig(managementConfig)
	if err != nil {
		return err
	}

	factory := informers.NewSharedInformerFactoryWithOptions(clientset, time.Second*60,
		informers.WithNamespace(project.ProjectNamespace(projectName)),
	)
	workspaceInformer := factory.Management().V1().DevSpaceWorkspaceInstances()

	self := baseClient.Self()
	filterByOwner := os.Getenv(provider.LOFT_FILTER_BY_OWNER) == "true"
	instanceStore := newStore(workspaceInformer, self, cmd.Context, filterByOwner, cmd.Log)

	_, err = workspaceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			instance, ok := obj.(*managementv1.DevSpaceWorkspaceInstance)
			if !ok {
				return
			}
			instanceStore.Add(instance)
			printInstances(stdout, instanceStore.List())
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			oldInstance, ok := oldObj.(*managementv1.DevSpaceWorkspaceInstance)
			if !ok {
				return
			}
			newInstance, ok := newObj.(*managementv1.DevSpaceWorkspaceInstance)
			if !ok {
				return
			}
			instanceStore.Update(oldInstance, newInstance)
			printInstances(stdout, instanceStore.List())
		},
		DeleteFunc: func(obj interface{}) {
			instance, ok := obj.(*managementv1.DevSpaceWorkspaceInstance)
			if !ok {
				// check for DeletedFinalStateUnknown. Can happen if the informer misses the delete event
				u, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					return
				}
				instance, ok = u.Obj.(*managementv1.DevSpaceWorkspaceInstance)
				if !ok {
					return
				}
			}
			instanceStore.Delete(instance)
			printInstances(stdout, instanceStore.List())
		},
	})
	if err != nil {
		return err
	}

	stopCh := make(chan struct{})
	defer close(stopCh)
	go func() {
		factory.Start(stopCh)
		factory.WaitForCacheSync(stopCh)

		// Kick off initial message
		printInstances(stdout, instanceStore.List())
	}()

	<-stopCh

	return nil
}

type instanceStore struct {
	informer      informermanagementv1.DevSpaceWorkspaceInstanceInformer
	self          *managementv1.Self
	context       string
	filterByOwner bool

	m         sync.Mutex
	instances map[string]*ProWorkspaceInstance

	log log.Logger
}

func newStore(informer informermanagementv1.DevSpaceWorkspaceInstanceInformer, self *managementv1.Self, context string, filterByOwner bool, log log.Logger) *instanceStore {
	return &instanceStore{
		informer:      informer,
		self:          self,
		context:       context,
		filterByOwner: filterByOwner,
		instances:     map[string]*ProWorkspaceInstance{},
		log:           log,
	}
}

func (s *instanceStore) key(meta metav1.ObjectMeta) string {
	return fmt.Sprintf("%s/%s", meta.Namespace, meta.Name)
}

func (s *instanceStore) Add(instance *managementv1.DevSpaceWorkspaceInstance) {
	if s.filterByOwner && !platform.IsOwner(s.self, instance.Spec.Owner) {
		return
	}
	var source *provider.WorkspaceSource
	if instance.GetAnnotations() != nil && instance.GetAnnotations()[storagev1.DevSpaceWorkspaceSourceAnnotation] != "" {
		source = provider.ParseWorkspaceSource(instance.GetAnnotations()[storagev1.DevSpaceWorkspaceSourceAnnotation])
	}

	var ideConfig *provider.WorkspaceIDEConfig
	if instance.GetLabels() != nil && instance.GetLabels()[storagev1.DevSpaceWorkspaceIDLabel] != "" {
		id := instance.GetLabels()[storagev1.DevSpaceWorkspaceIDLabel]
		workspaceConfig, err := provider.LoadWorkspaceConfig(s.context, id)
		if err == nil {
			ideConfig = &workspaceConfig.IDE
		}
	}

	proInstance := &ProWorkspaceInstance{
		TypeMeta:   instance.TypeMeta,
		ObjectMeta: instance.ObjectMeta,
		Spec:       instance.Spec,
		Status: ProWorkspaceInstanceStatus{
			DevSpaceWorkspaceInstanceStatus: instance.Status,
			Source:                        source,
			IDE:                           ideConfig,
		},
	}

	key := s.key(instance.ObjectMeta)
	s.m.Lock()
	s.instances[key] = proInstance
	s.m.Unlock()
}

func (s *instanceStore) Update(oldInstance *managementv1.DevSpaceWorkspaceInstance, newInstance *managementv1.DevSpaceWorkspaceInstance) {
	s.Add(newInstance)
}

func (s *instanceStore) Delete(instance *managementv1.DevSpaceWorkspaceInstance) {
	if s.filterByOwner && !platform.IsOwner(s.self, instance.Spec.Owner) {
		return
	}

	s.m.Lock()
	defer s.m.Unlock()
	key := s.key(instance.ObjectMeta)
	delete(s.instances, key)
}

func (s *instanceStore) List() []*ProWorkspaceInstance {
	instanceList := []*ProWorkspaceInstance{}
	// Check local imported workspaces
	// Eventually this should be implemented by filtering based on ownership and access on the CRD, for now we're stuck with this approach...
	localWorkspaces, err := workspace.ListLocalWorkspaces(s.context, false, s.log)
	if err == nil {
		for _, workspace := range localWorkspaces {
			if workspace.Imported && workspace.Pro != nil {
				// get instance for imported workspace
				selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
					MatchLabels: map[string]string{
						storagev1.DevSpaceWorkspaceUIDLabel: workspace.UID,
					},
				})
				if err != nil {
					continue
				}

				l, err := s.informer.Lister().
					DevSpaceWorkspaceInstances(project.ProjectFromNamespace(workspace.Pro.Project)).
					List(selector)
				if err != nil {
					continue
				}
				if len(l) == 0 {
					continue
				}
				instance := l[0]
				s.m.Lock()
				if _, ok := s.instances[s.key(instance.ObjectMeta)]; ok {
					continue
				}
				s.m.Unlock()

				var source *provider.WorkspaceSource
				if instance.GetAnnotations() != nil && instance.GetAnnotations()[storagev1.DevSpaceWorkspaceSourceAnnotation] != "" {
					source = provider.ParseWorkspaceSource(instance.GetAnnotations()[storagev1.DevSpaceWorkspaceSourceAnnotation])
				}

				var ideConfig *provider.WorkspaceIDEConfig
				if instance.GetLabels() != nil && instance.GetLabels()[storagev1.DevSpaceWorkspaceIDLabel] != "" {
					id := instance.GetLabels()[storagev1.DevSpaceWorkspaceIDLabel]
					workspaceConfig, err := provider.LoadWorkspaceConfig(s.context, id)
					if err == nil {
						ideConfig = &workspaceConfig.IDE
					}
				}

				proInstance := &ProWorkspaceInstance{
					TypeMeta:   instance.TypeMeta,
					ObjectMeta: instance.ObjectMeta,
					Spec:       instance.Spec,
					Status: ProWorkspaceInstanceStatus{
						DevSpaceWorkspaceInstanceStatus: instance.Status,
						Source:                        source,
						IDE:                           ideConfig,
					},
				}
				instanceList = append(instanceList, proInstance)
			}
		}
	}

	s.m.Lock()
	for _, instance := range s.instances {
		instanceList = append(instanceList, instance)
	}
	s.m.Unlock()

	return instanceList
}

func printInstances(w io.Writer, instances []*ProWorkspaceInstance) {
	out, err := json.Marshal(instances)
	if err != nil {
		return
	}

	fmt.Fprintln(w, string(out))
}
