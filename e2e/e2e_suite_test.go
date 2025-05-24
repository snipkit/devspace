package e2e

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"

	"dev.khulnasoft.com/e2e/framework"

	// Register tests
	_ "dev.khulnasoft.com/e2e/tests/build"
	_ "dev.khulnasoft.com/e2e/tests/context"
	_ "dev.khulnasoft.com/e2e/tests/ide"
	_ "dev.khulnasoft.com/e2e/tests/integration"
	_ "dev.khulnasoft.com/e2e/tests/machine"
	_ "dev.khulnasoft.com/e2e/tests/machineprovider"
	_ "dev.khulnasoft.com/e2e/tests/provider"
	_ "dev.khulnasoft.com/e2e/tests/ssh"
	_ "dev.khulnasoft.com/e2e/tests/up"
)

// TestRunE2ETests checks configuration parameters (specified through flags) and then runs
// E2E tests using the Ginkgo runner.
// If a "report directory" is specified, one or more JUnit test reports will be
// generated in this directory, and cluster logs will also be saved.
// This function is called on each Ginkgo node in parallel mode.
func TestRunE2ETests(t *testing.T) {
	if runtime.GOOS != "linux" {
		go framework.ServeAgent()

		// wait for http server to be up and running
		for {
			time.Sleep(time.Second)
			if os.Getenv("DEVSPACE_AGENT_URL") != "" {
				break
			}
		}
	}
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DevSpace e2e suite")
}
