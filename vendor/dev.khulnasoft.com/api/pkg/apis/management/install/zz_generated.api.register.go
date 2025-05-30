// Code generated by generator. DO NOT EDIT.

package install

import (
	"dev.khulnasoft.com/api/pkg/apis/management"
	v1 "dev.khulnasoft.com/api/pkg/apis/management/v1"
	"dev.khulnasoft.com/apiserver/pkg/builders"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func init() {
	Install(builders.Scheme)
}

func Install(scheme *runtime.Scheme) {
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(management.AddToScheme(scheme))
	utilruntime.Must(addKnownTypes(scheme))
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(management.SchemeGroupVersion,
		&management.AgentAuditEvent{},
		&management.AgentAuditEventList{},
		&management.Announcement{},
		&management.AnnouncementList{},
		&management.App{},
		&management.AppList{},
		&management.AppCredentials{},
		&management.Backup{},
		&management.BackupList{},
		&management.BackupApply{},
		&management.Cluster{},
		&management.ClusterList{},
		&management.ClusterAccessKey{},
		&management.ClusterAgentConfig{},
		&management.ClusterCharts{},
		&management.ClusterDomain{},
		&management.ClusterMemberAccess{},
		&management.ClusterMembers{},
		&management.ClusterReset{},
		&management.ClusterVirtualClusterDefaults{},
		&management.ClusterAccess{},
		&management.ClusterAccessList{},
		&management.ClusterRoleTemplate{},
		&management.ClusterRoleTemplateList{},
		&management.Config{},
		&management.ConfigList{},
		&management.ConvertVirtualClusterConfig{},
		&management.ConvertVirtualClusterConfigList{},
		&management.DevSpaceEnvironmentTemplate{},
		&management.DevSpaceEnvironmentTemplateList{},
		&management.DevSpaceWorkspaceInstance{},
		&management.DevSpaceWorkspaceInstanceList{},
		&management.DevSpaceDeleteOptions{},
		&management.DevSpaceStatusOptions{},
		&management.DevSpaceSshOptions{},
		&management.DevSpaceWorkspaceInstanceState{},
		&management.DevSpaceStopOptions{},
		&management.DevSpaceWorkspaceInstanceTroubleshoot{},
		&management.DevSpaceUpOptions{},
		&management.DevSpaceWorkspacePreset{},
		&management.DevSpaceWorkspacePresetList{},
		&management.DevSpaceWorkspaceTemplate{},
		&management.DevSpaceWorkspaceTemplateList{},
		&management.DirectClusterEndpointToken{},
		&management.DirectClusterEndpointTokenList{},
		&management.Event{},
		&management.EventList{},
		&management.Feature{},
		&management.FeatureList{},
		&management.IngressAuthToken{},
		&management.IngressAuthTokenList{},
		&management.Kiosk{},
		&management.KioskList{},
		&management.License{},
		&management.LicenseList{},
		&management.LicenseRequest{},
		&management.LicenseToken{},
		&management.LicenseTokenList{},
		&management.KhulnasoftUpgrade{},
		&management.KhulnasoftUpgradeList{},
		&management.OIDCClient{},
		&management.OIDCClientList{},
		&management.OwnedAccessKey{},
		&management.OwnedAccessKeyList{},
		&management.Project{},
		&management.ProjectList{},
		&management.ProjectChartInfo{},
		&management.ProjectCharts{},
		&management.ProjectClusters{},
		&management.ProjectImportSpace{},
		&management.ProjectMembers{},
		&management.ProjectMigrateSpaceInstance{},
		&management.ProjectMigrateVirtualClusterInstance{},
		&management.ProjectRunners{},
		&management.ProjectTemplates{},
		&management.ProjectSecret{},
		&management.ProjectSecretList{},
		&management.RedirectToken{},
		&management.RedirectTokenList{},
		&management.RegisterVirtualCluster{},
		&management.RegisterVirtualClusterList{},
		&management.ResetAccessKey{},
		&management.ResetAccessKeyList{},
		&management.Runner{},
		&management.RunnerList{},
		&management.RunnerAccessKey{},
		&management.RunnerConfig{},
		&management.Self{},
		&management.SelfList{},
		&management.SelfSubjectAccessReview{},
		&management.SelfSubjectAccessReviewList{},
		&management.SharedSecret{},
		&management.SharedSecretList{},
		&management.SpaceInstance{},
		&management.SpaceInstanceList{},
		&management.SpaceTemplate{},
		&management.SpaceTemplateList{},
		&management.SubjectAccessReview{},
		&management.SubjectAccessReviewList{},
		&management.Task{},
		&management.TaskList{},
		&management.TaskLog{},
		&management.Team{},
		&management.TeamList{},
		&management.TeamAccessKeys{},
		&management.TeamClusters{},
		&management.TranslateVClusterResourceName{},
		&management.TranslateVClusterResourceNameList{},
		&management.User{},
		&management.UserList{},
		&management.UserAccessKeys{},
		&management.UserClusters{},
		&management.UserPermissions{},
		&management.UserProfile{},
		&management.VirtualClusterInstance{},
		&management.VirtualClusterInstanceList{},
		&management.VirtualClusterAccessKey{},
		&management.VirtualClusterExternalDatabase{},
		&management.VirtualClusterInstanceKubeConfig{},
		&management.VirtualClusterInstanceLog{},
		&management.VirtualClusterTemplate{},
		&management.VirtualClusterTemplateList{},
	)
	return nil
}
