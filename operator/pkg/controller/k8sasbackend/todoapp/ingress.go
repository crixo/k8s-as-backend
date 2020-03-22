package todoapp

import (
	"fmt"

	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	common "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (t TodoApp) ensureIngress(i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {
	resUtils := &common.ResourceUtils{
		Scheme:          t.Manager.GetScheme(),
		Client:          t.Manager.GetClient(),
		ResourceFactory: createIngress,
	}

	ingressName := common.CreateUniqueSecondaryResourceName(i, BaseName)
	nsn := types.NamespacedName{Name: ingressName, Namespace: i.Namespace}
	found := &extv1beta1.Ingress{}
	return nil, common.EnsureResource(found, nsn, i, resUtils)
}

func createIngress(resNamespacedName types.NamespacedName, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {

	serviceName := common.CreateUniqueSecondaryResourceName(i, BaseName)
	path := fmt.Sprintf("/%s/%s/%s(/|$)(.*)", i.Namespace, i.Name, todoAppUrlSegmentIdentifier)
	backend := &extv1beta1.IngressBackend{
		ServiceName: serviceName,
		ServicePort: intstr.FromInt(SvcPort),
	}

	return &extv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resNamespacedName.Name,
			Namespace: resNamespacedName.Namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/$2",
			},
		},
		Spec: extv1beta1.IngressSpec{
			Rules: []extv1beta1.IngressRule{
				getRule("", path, backend),
			},
		},
	}
}

func getRule(host string, path string, backend *extv1beta1.IngressBackend) extv1beta1.IngressRule {
	rule := extv1beta1.IngressRule{}
	rule.Host = host
	rule.HTTP = &extv1beta1.HTTPIngressRuleValue{
		Paths: []extv1beta1.HTTPIngressPath{
			extv1beta1.HTTPIngressPath{
				Path:    path,
				Backend: *backend,
			},
		},
	}
	return rule
}
