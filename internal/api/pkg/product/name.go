package product

import (
	"fmt"
	"os"
	"sync"

	"dev.khulnasoft.com/admin-apis/pkg/licenseapi"
	"k8s.io/klog/v2"
)

// Product is the global variable to be set at build time
var productName string = string(licenseapi.Khulnasoft)
var once sync.Once

func loadProductVar() {
	productEnv := os.Getenv("PRODUCT")
	if productEnv == string(licenseapi.DevSpacePro) {
		productName = string(licenseapi.DevSpacePro)
	} else if productEnv == string(licenseapi.VClusterPro) {
		productName = string(licenseapi.VClusterPro)
	} else if productEnv == string(licenseapi.Khulnasoft) {
		productName = string(licenseapi.Khulnasoft)
	} else if productEnv != "" {
		klog.TODO().Error(fmt.Errorf("unrecognized product %s", productEnv), "error parsing product", "product", productEnv)
	}
}

func Name() licenseapi.ProductName {
	once.Do(loadProductVar)
	return licenseapi.ProductName(productName)
}

// Name returns the name of the product
func DisplayName() string {
	khulnasoftDisplayName := "Khulnasoft"

	switch Name() {
	case licenseapi.DevSpacePro:
		return "DevSpace Pro"
	case licenseapi.VClusterPro:
		return "vCluster Platform"
	case licenseapi.Khulnasoft:
	}

	return khulnasoftDisplayName
}
