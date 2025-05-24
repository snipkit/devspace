package product

import (
	"fmt"
	"strings"

	"dev.khulnasoft.com/admin-apis/pkg/licenseapi"
)

// Replace replaces the product name in the given usage string
// based on the current product.Product().
//
// It replaces "khulnasoft" with the specific product name:
//   - "devspace pro" for product.DevSpacePro
//   - "vcluster platform" for product.VClusterPro
//   - No replacement for product.Khulnasoft
//
// This handles case insensitive replaces like "khulnasoft" -> "devspace pro".
//
// It also handles case sensitive replaces:
//   - "Khulnasoft" -> "DevSpace.Pro" for product.DevSpacePro
//   - "Khulnasoft" -> "vCluster Platform" for product.VClusterPro
//
// This allows customizing command usage text for different products.
//
// Parameters:
//   - content: The string to update
//
// Returns:
//   - The updated string with product name replaced if needed.
func Replace(content string) string {
	switch Name() {
	case licenseapi.DevSpacePro:
		content = strings.Replace(content, "khulnasoft.com", "devspace.pro", -1)
		content = strings.Replace(content, "khulnasoft.host", "devspace.host", -1)

		content = strings.Replace(content, "khulnasoft", "devspace pro", -1)
		content = strings.Replace(content, "Khulnasoft", "DevSpace.Pro", -1)
	case licenseapi.VClusterPro:
		content = strings.Replace(content, "khulnasoft.com", "vcluster.pro", -1)
		content = strings.Replace(content, "khulnasoft.host", "vcluster.host", -1)

		content = strings.Replace(content, "khulnasoft", "vcluster platform", -1)
		content = strings.Replace(content, "Khulnasoft", "vCluster Platform", -1)
	case licenseapi.Khulnasoft:
	}

	return content
}

// ReplaceWithHeader replaces the "khulnasoft" product name in the given
// usage string with the specific product name based on product.Product().
// It also adds a header with padding around the product name and usage.
//
// The product name replacements are:
//
//   - "devspace pro" for product.DevSpacePro
//   - "vcluster platform" for product.VClusterPro
//   - No replacement for product.Khulnasoft
//
// This handles case insensitive replaces like "khulnasoft" -> "devspace pro".
//
// It also handles case sensitive replaces:
//   - "Khulnasoft" -> "DevSpace.Pro" for product.DevSpacePro
//   - "Khulnasoft" -> "vCluster Platform" for product.VClusterPro
//
// Parameters:
//   - use: The usage string
//   - content: The content string to run product name replacement on
//
// Returns:
//   - The content string with product name replaced and header added
func ReplaceWithHeader(use, content string) string {
	maxChar := 56

	productName := licenseapi.Khulnasoft

	switch Name() {
	case licenseapi.DevSpacePro:
		productName = "devspace pro"
	case licenseapi.VClusterPro:
		productName = "vcluster platform"
	case licenseapi.Khulnasoft:
	}

	paddingSize := (maxChar - 2 - len(productName) - len(use)) / 2

	separator := strings.Repeat("#", paddingSize*2+len(productName)+len(use)+2+1)
	padding := strings.Repeat("#", paddingSize)

	return fmt.Sprintf(`%s
%s %s %s %s
%s
%s
`, separator, padding, productName, use, padding, separator, Replace(content))
}
