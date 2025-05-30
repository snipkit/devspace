// Code generated by generator. DO NOT EDIT.

package install

import (
	"dev.khulnasoft.com/agentapi/pkg/apis/khulnasoft/cluster"
	v1 "dev.khulnasoft.com/agentapi/pkg/apis/khulnasoft/cluster/v1"
	"dev.khulnasoft.com/apiserver/pkg/builders"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func init() {
	Install(builders.Scheme)
}

func Install(scheme *runtime.Scheme) {
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(cluster.AddToScheme(scheme))
	utilruntime.Must(addKnownTypes(scheme))
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(cluster.SchemeGroupVersion,
		&cluster.ChartInfo{},
		&cluster.ChartInfoList{},
		&cluster.Feature{},
		&cluster.FeatureList{},
		&cluster.HelmRelease{},
		&cluster.HelmReleaseList{},
	)
	return nil
}
