package k8sasbackend

import (
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ValidatingWebhookConfigurationV1Factory struct{}

func (f ValidatingWebhookConfigurationV1Factory) createEmpty() runtime.Object {
	return &admissionregistrationv1.ValidatingWebhookConfiguration{}
}

type ValidatingWebhookConfigurationV1Beta1Factory struct{}

func (f ValidatingWebhookConfigurationV1Beta1Factory) createEmpty() runtime.Object {
	return &admissionregistrationv1beta1.ValidatingWebhookConfiguration{}
}
