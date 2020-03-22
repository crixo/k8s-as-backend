package webhookserver

import (
	"context"
	"fmt"
	"reflect"
	"time"

	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	clusterdep "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/clusterdep"
	common "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	arv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (ws WebhookServer) ensureValidationWebhook(i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {

	nsn := types.NamespacedName{Name: clusterdep.ValidationWebhookConfigurationName, Namespace: ""}
	found := &arv1beta1.ValidatingWebhookConfiguration{}
	err := ws.Client.Get(context.TODO(), nsn, found)
	log.Info(fmt.Sprintf("Getting %T %s", found, found.Name))
	if err != nil && errors.IsNotFound(err) {
		log.Info(fmt.Sprintf("Cluster resource %T %s not found, requeuing", found, found.Name))
		return &reconcile.Result{Requeue: true, RequeueAfter: time.Second * 2}, err
	} else if err != nil {
		log.Error(err, fmt.Sprintf("Failed to get %T", found))
		return &reconcile.Result{}, err
	}

	desiredWh := createValidationWebhook(i)
	currentWebhook, foundIdx := findWebhook(found, todosWebhookName)
	if currentWebhook != nil && reflect.DeepEqual(currentWebhook, desiredWh) {
		return nil, nil
	}

	if foundIdx > -1 {
		// a := found.Webhooks
		// i := foundIdx
		// a[i] = a[len(a)-1] // Copy last element to index i.
		// a[len(a)-1] = nil   // Erase last element (write zero value).
		// a = a[:len(a)-1]   // Truncate slice.
		found.Webhooks[foundIdx] = desiredWh
	} else {
		found.Webhooks = append(found.Webhooks, desiredWh)
	}

	log.Info(fmt.Sprintf("Updating %T %s", found, found.Name))
	err = ws.Client.Update(context.TODO(), found)
	if err != nil {
		log.Error(err, fmt.Sprintf("Failed to get %T", found))
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func findWebhook(res *arv1beta1.ValidatingWebhookConfiguration, webhookName string) (*arv1beta1.ValidatingWebhook, int) {
	for i, wh := range res.Webhooks {
		if wh.Name == webhookName {
			return &wh, i
		}
	}
	return nil, -1
}

func createValidationWebhook(i *k8sasbackendv1alpha1.K8sAsBackend) arv1beta1.ValidatingWebhook {
	scope := arv1beta1.NamespacedScope
	path := "/crd"
	sideEffects := arv1beta1.SideEffectClassNone
	caBundle := common.AppState.ClientConfig.CAData
	var timeout int32 = 5
	serviceName := common.CreateUniqueSecondaryResourceName(i, baseName)
	return arv1beta1.ValidatingWebhook{
		Name: todosWebhookName,
		Rules: []arv1beta1.RuleWithOperations{{
			Rule: arv1beta1.Rule{
				APIGroups:   []string{"k8sasbackend.com"},
				APIVersions: []string{"v1"},
				Resources:   []string{"todos"},
				Scope:       &scope,
			},
			Operations: []arv1beta1.OperationType{"*"},
		}},
		ClientConfig: arv1beta1.WebhookClientConfig{
			Service: &arv1beta1.ServiceReference{
				Namespace: i.Namespace,
				Name:      serviceName,
				Path:      &path,
			},
			CABundle: caBundle,
		},
		AdmissionReviewVersions: []string{"v1beta1"},
		SideEffects:             &sideEffects,
		TimeoutSeconds:          &timeout,
	}
}
