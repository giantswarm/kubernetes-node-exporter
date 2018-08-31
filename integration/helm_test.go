// +build k8srequired

package integration

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/framework"
	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	resourceNamespace = "kube-system"
)

var (
	f *framework.Host
)

func TestHelm(t *testing.T) {
	channel := os.Getenv("CIRCLE_SHA1")

	err := framework.HelmCmd(fmt.Sprintf("registry install --wait quay.io/giantswarm/kubernetes-node-exporter-chart:%s -n test-deploy", channel))
	if err != nil {
		t.Errorf("unexpected error during installation of the chart: %v", err)
	}
	defer framework.HelmCmd("delete test-deploy --purge")

	err = checkDaemonSet()
	if err != nil {
		t.Fatalf("daemonset manifest is incorrect: %v", err)
	}

	err = framework.HelmCmd("test --debug --cleanup test-deploy")
	if err != nil {
		t.Errorf("unexpected error during test of the chart: %v", err)
	}
}

// TestMigration ensures that previously deployed resources are properly removed.
// It installs a chart with the same resources as node-exporter and apprpriate
// labels so that we can query for them. Then installs node-operator chart and
// checks that the previous resources are removed and the ones from node-exporter
// are in place.
func TestMigration(t *testing.T) {
	// install legacy resources
	err := framework.HelmCmd("install /e2e/fixtures/resources-chart -n resources")
	if err != nil {
		t.Fatalf("could not install resources chart: %v", err)
	}
	defer framework.HelmCmd("delete resources --purge")

	// check legacy resources are present
	err = checkResourcesPresent("app=node-exporter,kind=legacy")
	if err != nil {
		t.Fatalf("could check legacy resources present: %v", err)
	}
	// check managed resources are not present
	err = checkResourcesNotPresent("app=node-exporter,giantswarm.io/service-type=managed")
	if err != nil {
		t.Fatalf("could check managed resources not present: %v", err)
	}

	// install kubernetes-node-exporter-chart
	channel := os.Getenv("CIRCLE_SHA1")
	err = framework.HelmCmd(fmt.Sprintf("registry install --wait quay.io/giantswarm/kubernetes-node-exporter-chart:%s -n test-deploy", channel))
	if err != nil {
		t.Fatalf("could not install kubernetes-node-exporter-chart: %v", err)
	}
	defer framework.HelmCmd("delete test-deploy --purge")

	// check legacy resources are not present
	err = checkResourcesNotPresent("app=node-exporter,kind=legacy")
	if err != nil {
		t.Fatalf("could check legacy resources not present: %v", err)
	}
	// check managed resources are present
	err = checkResourcesPresent("app=node-exporter,giantswarm.io/service-type=managed")
	if err != nil {
		t.Fatalf("could check managed resources present: %v", err)
	}
}

func checkResourcesPresent(labelSelector string) error {
	c := f.K8sClient()
	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
	}

	ds, err := c.Extensions().DaemonSets(resourceNamespace).List(listOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(ds.Items) != 1 {
		return microerror.Newf("unexpected number of daemonsets, want 1, got %d", len(ds.Items))
	}

	r, err := c.Rbac().Roles(resourceNamespace).List(listOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(r.Items) != 1 {
		return microerror.Newf("unexpected number of roles, want 1, got %d", len(r.Items))
	}

	rb, err := c.Rbac().RoleBindings(resourceNamespace).List(listOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(rb.Items) != 1 {
		return microerror.Newf("unexpected number of rolebindings, want 1, got %d", len(rb.Items))
	}

	s, err := c.Core().Services(resourceNamespace).List(listOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(s.Items) != 1 {
		return microerror.Newf("unexpected number of services, want 1, got %d", len(s.Items))
	}

	sa, err := c.Core().ServiceAccounts(resourceNamespace).List(listOptions)
	if err != nil {
		return microerror.Mask(err)
	}
	if len(sa.Items) != 1 {
		return microerror.Newf("unexpected number of serviceaccountss, want 1, got %d", len(sa.Items))
	}

	return nil
}

func checkResourcesNotPresent(labelSelector string) error {
	c := f.K8sClient()
	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
	}

	ds, err := c.Extensions().DaemonSets(resourceNamespace).List(listOptions)
	if err == nil && len(ds.Items) > 0 {
		return microerror.New("expected error querying for daemonsets didn't happen")
	}
	if !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
	}

	r, err := c.Rbac().Roles(resourceNamespace).List(listOptions)
	if err == nil && len(r.Items) > 0 {
		return microerror.New("expected error querying for roles didn't happen")
	}
	if !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
	}

	rb, err := c.Rbac().RoleBindings(resourceNamespace).List(listOptions)
	if err == nil && len(rb.Items) > 0 {
		return microerror.New("expected error querying for rolebindings didn't happen")
	}
	if !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
	}

	s, err := c.Core().Services(resourceNamespace).List(listOptions)
	if err == nil && len(s.Items) > 0 {
		return microerror.New("expected error querying for services didn't happen")
	}
	if !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
	}

	sa, err := c.Core().ServiceAccounts(resourceNamespace).List(listOptions)
	if err == nil && len(sa.Items) > 0 {
		return microerror.New("expected error querying for serviceaccounts didn't happen")
	}
	if !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
	}

	return nil
}

// checkDaemonSet ensures that key properties of the node-exporter daemonset are
// correct.
func checkDaemonSet() error {
	name := "node-exporter"
	expectedLabels := map[string]string{
		"app": "node-exporter",
		"giantswarm.io/service-type": "managed",
	}
	expectedMatchLabels := map[string]string{
		"app": "node-exporter",
	}

	c := f.K8sClient()
	ds, err := c.Apps().DaemonSets(resourceNamespace).Get(name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return microerror.Newf("could not find daemonset: '%s' %v", name, err)
	} else if err != nil {
		return microerror.Newf("unexpected error getting daemonset: %v", err)
	}

	// Check daemonset labels.
	if !reflect.DeepEqual(expectedLabels, ds.ObjectMeta.Labels) {
		return microerror.Newf("expected labels: %v got: %v", expectedLabels, ds.ObjectMeta.Labels)
	}

	// Check selector match labels.
	if !reflect.DeepEqual(expectedMatchLabels, ds.Spec.Selector.MatchLabels) {
		return microerror.Newf("expected match labels: %v got: %v", expectedMatchLabels, ds.Spec.Selector.MatchLabels)
	}

	// Check pod labels.
	if !reflect.DeepEqual(expectedLabels, ds.Spec.Template.ObjectMeta.Labels) {
		return microerror.Newf("expected pod labels: %v got: %v", expectedLabels, ds.Spec.Template.ObjectMeta.Labels)
	}

	return nil
}
