package authz

import (
	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	common "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (a Authz) ensureAccount(i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {
	resUtils := &common.ResourceUtils{
		Scheme:          a.Manager.GetScheme(),
		Client:          a.Manager.GetClient(),
		ResourceFactory: createAccount,
	}

	serviceAccountName := common.CreateUniqueSecondaryResourceName(i, BaseName)
	nsn := types.NamespacedName{Name: serviceAccountName, Namespace: i.Namespace}
	found := &corev1.ServiceAccount{}
	return nil, common.EnsureResource(found, nsn, i, resUtils)
}

func createAccount(nsn types.NamespacedName, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {

	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nsn.Name,
			Namespace: nsn.Namespace,
		},
	}
}
