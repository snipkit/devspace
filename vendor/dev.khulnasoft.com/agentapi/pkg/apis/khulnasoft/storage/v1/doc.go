// Api versions allow the api contract for a resource to be changed while keeping
// backward compatibility by support multiple concurrent versions
// of the same resource

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package
// +k8s:defaulter-gen=TypeMeta
// +groupName=storage.khulnasoft.com
package v1 // import "dev.khulnasoft.com/agentapi/apis/khulnasoft/storage/v1"
