package webhookserver

import (
	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"

	certv1beta1 "k8s.io/client-go/kubernetes/typed/certificates/v1beta1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	csrName    string = "admission-webhook-example-svc.default"
	secretName string = "admission-webhook-example-certs"
)

var log = logf.Log.WithName("controller_k8sasbackend_webhook")

// Users/cristiano/Coding/golang/k8s-as-backend/operator/certs/server-cert.pem
// Users/cristiano/Coding/golang/k8s-as-backend/operator/certs/server-key.pem
type WebhookServer struct {
	Client      client.Client
	Scheme      *runtime.Scheme
	CerFilePath string
	KeyFilePath string
	CertClient  certv1beta1.CertificateSigningRequestInterface
}

func (ws WebhookServer) GetWatchedResources() []runtime.Object {
	return []runtime.Object{
		&corev1.Secret{},
		&appsv1.Deployment{},
		&corev1.Service{},
		&admissionregistrationv1beta1.ValidatingWebhookConfiguration{},
	}
}

func (ws *WebhookServer) Reconcile(i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {

	err := ws.ensureSecret(i)
	if err != nil {
		return &reconcile.Result{}, err
	}

	err = ws.ensureDeployment(i)
	if err != nil {
		return &reconcile.Result{}, err
	}

	return nil, nil
}
