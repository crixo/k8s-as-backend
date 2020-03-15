package clusterdep

import (
	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	common "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	arv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (cd ClusterDependencies) ensureValidationWebhookConfiguration() (*reconcile.Result, error) {
	resUtils := &common.ResourceUtils{
		Scheme:          cd.Scheme,
		Client:          cd.Client,
		ResourceFactory: createValidationWebhookConfiguration,
	}

	nsn := types.NamespacedName{Name: ValidationWebhookConfigurationName, Namespace: ""}
	found := &arv1beta1.ValidatingWebhookConfiguration{}
	return nil, common.EnsureResource(found, nsn, nil, resUtils)
}

func createValidationWebhookConfiguration(nsName types.NamespacedName, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {
	// scope := arv1beta1.NamespacedScope
	// path := "/crd"
	// sideEffects := arv1beta1.SideEffectClassNone
	// caBundle := common.AppState.ClientConfig.CAData
	// var timeout int32 = 5
	return &arv1beta1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: nsName.Name,
			// hack to force SetControllerReference where primary resource ns must match the secondary resource ns
			//Namespace: i.Namespace,
			//Namespace: nsName.Namespace,
		},
		Webhooks: []arv1beta1.ValidatingWebhook{
			// 	{
			// 	Name: vebhookName,
			// 	Rules: []arv1beta1.RuleWithOperations{{
			// 		Rule: arv1beta1.Rule{
			// 			APIGroups:   []string{"k8sasbackend.com"},
			// 			APIVersions: []string{"v1"},
			// 			Resources:   []string{"todos"},
			// 			Scope:       &scope,
			// 		},
			// 		Operations: []arv1beta1.OperationType{"*"},
			// 	}},
			// 	ClientConfig: arv1beta1.WebhookClientConfig{
			// 		Service: &arv1beta1.ServiceReference{
			// 			Namespace: i.Namespace,
			// 			Name:      serviceName,
			// 			Path:      &path,
			// 		},
			// 		CABundle: caBundle,
			// 	},
			// 	AdmissionReviewVersions: []string{"v1beta1"},
			// 	SideEffects:             &sideEffects,
			// 	TimeoutSeconds:          &timeout,
			// }
		},
	}
}
