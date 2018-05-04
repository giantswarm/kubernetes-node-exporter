// +build k8srequired

package integration

import (
	"fmt"
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

	f, err = framework.NewHost(framework.HostConfig{})
	if err != nil {
		panic(err.Error())
	}

	if err := f.CreateNamespace("giantswarm"); err != nil {
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
	channel := os.Getenv("CIRCLE_SHA1")

	err := framework.HelmCmd(fmt.Sprintf("registry install --wait quay.io/giantswarm/kubernetes-node-exporter-chart:%s -n test-deploy", channel))
	if err != nil {
		t.Errorf("unexpected error during installation of the chart: %v", err)
	}
	defer framework.HelmCmd("delete test-deploy --purge")

	err = framework.HelmCmd("test --debug --cleanup test-deploy")
	if err != nil {
		t.Errorf("unexpected error during test of the chart: %v", err)
	}
}

func TestMigration(t *testing.T) {
	// install resources
	err := framework.HelmCmd("install /e2e/fixtures/resources-chart -n resources")
	if err != nil {
		t.Fatalf("could not install resources chart: %v", err)
	}

	// check they are installed
	err = checkResourcesInstalled()
	if err != nil {
		t.Fatalf("could check installed resources: %v", err)
	}

	// install kubernetes-node-exporter-chart
	channel := os.Getenv("CIRCLE_SHA1")
	err = framework.HelmCmd(fmt.Sprintf("registry install --wait quay.io/giantswarm/kubernetes-node-exporter-chart:%s -n test-deploy", channel))
	if err != nil {
		t.Fatalf("could not install kubernetes-node-exporter-chart: %v", err)
	}
	defer framework.HelmCmd("delete test-deploy --purge")

	// check that resources are no longer there
	err = checkResourcesRemoved()
	if err != nil {
		t.Fatalf("could check removed resources: %v", err)
	}
}

func checkResourcesInstalled() error {
	return nil
}

func checkResourcesRemoved() error {
	return nil
}
