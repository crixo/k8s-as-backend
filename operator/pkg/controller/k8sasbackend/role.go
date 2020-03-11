package k8sasbackend

import (
	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type RoleFactory struct{}

func (f RoleFactory) createEmpty() runtime.Object {
	return &rbac.Role{}
}

func (f RoleFactory) getNames() []string {
	return []string{"todo-crd"}
}

func (f RoleFactory) ensure(r *ReconcileK8sAsBackend,
	request reconcile.Request,
	i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {

	return r.ensureResource(request, i, createNamespacedName(f.getNames()[0], i.Namespace), f)

}

func (f RoleFactory) create(name string, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {
	return &rbac.Role{
		ObjectMeta: createMeta(name, i.Namespace),
		Rules: []rbac.PolicyRule{{
			APIGroups: []string{"k8sasbackend.com"},
			Resources: []string{"todo", "todos"},
			Verbs:     []string{"get", "list", "delete", "watch", "update"},
		},
		},
	}
}
