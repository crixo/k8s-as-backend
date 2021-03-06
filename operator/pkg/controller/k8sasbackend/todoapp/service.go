package todoapp

import (
	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	common "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (t TodoApp) ensureService(i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {
	resUtils := &common.ResourceUtils{
		Scheme:          t.Manager.GetScheme(),
		Client:          t.Manager.GetClient(),
		ResourceFactory: createService,
	}

	serviceName := common.CreateUniqueSecondaryResourceName(i, BaseName)
	nsn := types.NamespacedName{Name: serviceName, Namespace: i.Namespace}
	found := &corev1.Service{}
	return nil, common.EnsureResource(found, nsn, i, resUtils)
}

func createService(nsn types.NamespacedName, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {
	matchingLabels := common.CreateMatchingLabels(i, BaseName)
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nsn.Name,
			Namespace: nsn.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: matchingLabels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       int32(SvcPort),
				TargetPort: intstr.FromInt(todoPort),
			}},
		},
	}
}
