// +build k8srequired

package integration

import (
	"log"
	"os"
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/framework"
)

var (
	f *framework.Host
)

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	var v int
	var err error
	f, err = framework.NewHost()
	if err != nil {
		log.Printf("unexpected error: %v\n", err)
		os.Exit(1)
	}

	if err := f.Setup(); err != nil {
		log.Printf("unexpected error: %v\n", err)
		v = 1
	}

	if v == 0 {
		v = m.Run()
	}

	if os.Getenv("KEEP_RESOURCES") != "true" {
		f.Teardown()
	}

	os.Exit(v)
}

func TestHelm(t *testing.T) {
	log.Printf("Testing")

	err := framework.HelmCmd("install quay.io/giantswarm/kubernetes-node-exporter-chart")
	if err != nil {
		t.Errorf("unexpected error during installation of the chart %v", err)
	}

}
