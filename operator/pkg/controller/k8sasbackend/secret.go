package k8sasbackend

import (
	"io/ioutil"

	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type SecretFactory struct{}

func (f SecretFactory) createEmpty() runtime.Object {
	return &corev1.Secret{}
}

func (f SecretFactory) getNames() []string {
	return []string{"admission-webhook-example-certs"}
}

func (f SecretFactory) ensure(r *ReconcileK8sAsBackend,
	request reconcile.Request,
	i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {

	return r.ensureResource(request, i, createNamespacedName(f.getNames()[0], i.Namespace), f)

}

func (f SecretFactory) create(name string, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {
	cert, _ := ioutil.ReadFile("/Users/cristiano/Coding/golang/k8s-as-backend/operator/server-cert.pem")
	key, _ := ioutil.ReadFile("/Users/cristiano/Coding/golang/k8s-as-backend/operator/server-key.pem")
	return &corev1.Secret{
		ObjectMeta: createMeta(name, i.Namespace),
		Type:       "Opaque",

		Data: map[string][]byte{
			"key.pem":  cert,
			"cert.pem": key,
		},
	}
}
