package e2e

import (
	goctx "context"
	"fmt"
	"testing"
	"time"

	"github.com/crixo/k8s-as-backend/operator/pkg/apis"
	operator "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestK8sAsBackend(t *testing.T) {
	k8sAsBackendList := &operator.K8sAsBackendList{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, k8sAsBackendList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("memcached-group", func(t *testing.T) {
		t.Run("Cluster", MemcachedCluster)
		//t.Run("Cluster2", MemcachedCluster)
	})
}

func k8sasbackendTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	// create memcached custom resource
	exampleK8sAsBackend := &operator.K8sAsBackend{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-k8sasbackend",
			Namespace: namespace,
		},
		Spec: operator.K8sAsBackendSpec{
			Size:           2,
			ProductVersion: "v0.0.0",
		},
	}
	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), exampleK8sAsBackend, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}
	// wait for example-memcached to reach 3 replicas
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-k8sasbackend-todo-app", 2, retryInterval, timeout)
	if err != nil {
		return err
	}

	return e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-k8sasbackend-todos-webhook-server", 1, retryInterval, timeout)

	// err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "example-k8sasbackend-todo-app", Namespace: namespace}, exampleK8sAsBackend)
	// if err != nil {
	// 	return err
	// }
	// exampleK8sAsBackend.Spec.Size = 4
	// err = f.Client.Update(goctx.TODO(), exampleK8sAsBackend)
	// if err != nil {
	// 	return err
	// }

	// // wait for example-memcached to reach 4 replicas
	// return e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-k8sasbackend-todo-app", 4, retryInterval, timeout)
}

func MemcachedCluster(t *testing.T) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()

	// t.Log("Initializing cluster resources")
	// err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	// if err != nil {
	// 	t.Fatalf("failed to initialize cluster resources: %v", err)
	// }
	// t.Log("Initialized cluster resources")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("namespace:" + namespace)
	// get global framework variables
	f := framework.Global
	// // wait for memcached-operator to be ready
	// //err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "memcached-operator", 1, retryInterval, timeout)
	// err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "k8sasbackend-operator", 1, retryInterval, timeout)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	if err = k8sasbackendTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
