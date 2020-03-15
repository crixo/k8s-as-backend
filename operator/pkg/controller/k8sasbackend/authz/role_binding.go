package authz

import (
	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	common "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	rbac "k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (a Authz) ensureRoleBinding(i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {
	resUtils := &common.ResourceUtils{
		Scheme:          a.Manager.GetScheme(),
		Client:          a.Manager.GetClient(),
		ResourceFactory: createRoleBinding,
	}

	nsn := types.NamespacedName{Name: roleBindingName, Namespace: i.Namespace}
	found := &rbac.RoleBinding{}
	return nil, common.EnsureResource(found, nsn, i, resUtils)
}

func createRoleBinding(nsn types.NamespacedName, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {

	return &rbac.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nsn.Name,
			Namespace: nsn.Namespace,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
		Subjects: []rbac.Subject{{
			Kind:      "ServiceAccount",
			Name:      ServiceAccountName,
			Namespace: i.Namespace,
		}},
	}
}
