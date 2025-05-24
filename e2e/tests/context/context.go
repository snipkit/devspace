package context

import (
	"context"
	"os"

	"dev.khulnasoft.com/e2e/framework"
	"github.com/onsi/ginkgo/v2"
)

var _ = DevSpaceDescribe("devspace context test suite", func() {
	ginkgo.Context("testing context command", ginkgo.Label("context"), ginkgo.Ordered, func() {
		ctx := context.Background()
		initialDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		ginkgo.It("create a new context, switch to it and delete afterwards", func() {
			f := framework.NewDefaultFramework(initialDir + "/bin")

			err = f.DevSpaceContextCreate(ctx, "test-context")
			framework.ExpectNoError(err)

			err = f.DevSpaceContextUse(context.Background(), "test-context")
			framework.ExpectNoError(err)

			err = f.DevSpaceContextDelete(context.Background(), "test-context")
			framework.ExpectNoError(err)
		})

	})
})
