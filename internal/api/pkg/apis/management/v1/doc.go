// Api versions allow the api contract for a resource to be changed while keeping
// backward compatibility by support multiple concurrent versions
// of the same resource

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=dev.khulnasoft.com/api/pkg/apis/management
// +k8s:defaulter-gen=TypeMeta
// +groupName=management.khulnasoft.com
package v1 // import "dev.khulnasoft.com/api/pkg/apis/management/v1"
