package k8sasbackend

import (
	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type AccountFactory struct{}

func (f AccountFactory) createEmpty() runtime.Object {
	return &corev1.ServiceAccount{}
}

func (f AccountFactory) getNames() []string {
	return []string{"todo-crd"}
}

func (f AccountFactory) ensure(r *ReconcileK8sAsBackend,
	request reconcile.Request,
	i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {
	return r.ensureResource(request, i, createNamespacedName(f.getNames()[0], i.Namespace), f)
}

func (f AccountFactory) create(name string, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {
	return &corev1.ServiceAccount{
		ObjectMeta: createMeta(name, i.Namespace),
	}
}
