package k8sasbackend

// import (
// 	"crypto/x509"
// 	"encoding/pem"
// 	"testing"

// 	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/types"
// 	"sigs.k8s.io/controller-runtime/pkg/client/fake"
// 	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

// 	"k8s.io/apimachinery/pkg/runtime"
// 	"k8s.io/client-go/kubernetes/scheme"

// 	"sigs.k8s.io/controller-runtime/pkg/reconcile"

// 	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

// 	certfake "k8s.io/client-go/kubernetes/typed/certificates/v1beta1/fake"

// 	"github.com/stretchr/testify/assert"

// 	clientgotesting "k8s.io/client-go/testing"

// 	apicert "k8s.io/api/certificates/v1beta1"

// 	b64 "encoding/base64"
// )

// const certBase64String = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURxVENDQXBHZ0F3SUJBZ0lVWUxLWFVlRDcrRjRURnE4dkJ4VXpvano1UUhVd0RRWUpLb1pJaHZjTkFRRUwKQlFBd0ZURVRNQkVHQTFVRUF4TUthM1ZpWlhKdVpYUmxjekFlRncweU1EQXpNVEV4T1RRNE1EQmFGdzB5TVRBegpNVEV4T1RRNE1EQmFNRFF4TWpBd0JnTlZCQU1US1dGa2JXbHpjMmx2YmkxM1pXSm9iMjlyTFdWNFlXMXdiR1V0CmMzWmpMbVJsWm1GMWJIUXVjM1pqTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEKNlFNR1c0aFZFd3dGN2NCdEVqaE02aGUzSThQNDc5ZVQraXhCQWpXcjR5L2dkWURXeFprNG5kWnY4SHVKY0JWTQpDNUEwSG5jZXNlaFMrTmJ5a2dxN1lDd01LM3FDUlJPMVpKa1Y3d1JvSDhIOU91ZmxsczdWY0hQTDMyeU5jVnQxCnVGWFlDUUpKVUYxUjRXczJVNnFvUlVFN040T1plQXRWcnM2YmpVZzZuSFR0ZWtUUUt5ekxHMWV6UG5rNWRtNFEKdTlXVFVBa0lSTm8vemQrTkJkaDYwdy96UHVTZkRDcFV0U1dkYWFQUnVjb3hTV0ZVWDROWEYzNUhoVW91OTk3QgpoZWl6TEhiVndiK2Fqc3I1UEVJL2ZOd1RwaWtoNnZiVGkvRDZOMGRVdW94enBDSlVtSnkwSjFKc0pDSnVZeHAxCjNjeHpUbjdMUEdBcHVnR3JnWEZ1NlFJREFRQUJvNEhSTUlIT01BNEdBMVVkRHdFQi93UUVBd0lGb0RBVEJnTlYKSFNVRUREQUtCZ2dyQmdFRkJRY0RBVEFNQmdOVkhSTUJBZjhFQWpBQU1CMEdBMVVkRGdRV0JCU3VKUWJVcWhacwpYdS9KK3dGb3RpM0toOEpUSkRCNkJnTlZIUkVFY3pCeGdoMWhaRzFwYzNOcGIyNHRkMlZpYUc5dmF5MWxlR0Z0CmNHeGxMWE4yWTRJbFlXUnRhWE56YVc5dUxYZGxZbWh2YjJzdFpYaGhiWEJzWlMxemRtTXVaR1ZtWVhWc2RJSXAKWVdSdGFYTnphVzl1TFhkbFltaHZiMnN0WlhoaGJYQnNaUzF6ZG1NdVpHVm1ZWFZzZEM1emRtTXdEUVlKS29aSQpodmNOQVFFTEJRQURnZ0VCQUJBU0ZQT05JQTBkV3Z3ekF2NmJJN2Y1clhmVWtWZVkwZUIxZmx0UTlndUM0Z3JHCldsZ0gvOGxHQThCVjVNUUxTaWsxV1FiVmtackYySDVSNEFlaUVqYkxHQlRJc2doYXROSFZLZitjdEprbGhjRysKVXc5bk5HSnVhWWdvNkx5NmRBS0R5YXNuWmlHNDNKL3NOM0tHQzJPRktLUGFHK0pYWHFYbWxwU3RNbWV1OVpsagp0R2RaVDFiUktNeWlCa3cwMDd2R2pENjMzSzFBU2pLbHRDeEFwY2FZUVlkeTY5RC9kUHgweU5XcmdkUGlZd2hRClBwWThqOFVpMExaNGd1NlRtYXRPb0xGSW84OGxwUGt3RVpOZUthbENuRGNLWW9OKzByT0o2QTRzZExJVmNZNm8KSkJBc1NPSFpSa29qNkNFcG45MGF1SVdqSmtkdnM5WVNhLzNUUC9BPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="

// func TestConverCertificateStringTofile(t *testing.T) {
// 	logf.SetLogger(logf.ZapLogger(true))

// 	certPEM, err := b64.StdEncoding.DecodeString(certBase64String)
// 	if err != nil {
// 		log.Error(err, "Unable to decode cert base64 string")
// 	}

// 	log.Info("cert file content", "certPEM", string(certPEM))

// 	block, _ := pem.Decode([]byte(certPEM))
// 	if block == nil {
// 		panic("failed to parse certificate PEM")
// 	}
// 	cert, err := x509.ParseCertificate(block.Bytes)
// 	if err != nil {
// 		panic("failed to parse certificate: " + err.Error())
// 	}

// 	//t.Log(cert.Subject)
// 	log.Info("cert", "Subject", cert.Subject)
// }

// func TestReadCertificate(t *testing.T) {
// 	certFactory := &CertFactory{}

// 	certFactory.readSigningRequest()

// }

// func TestWebhookCertificate(t *testing.T) {
// 	// Set the logger to development mode for verbose logs.
// 	logf.SetLogger(logf.ZapLogger(true))

// 	var (
// 		name            = "memcached-operator"
// 		namespace       = "memcached"
// 		replicas  int32 = 3
// 	)

// 	i := &k8sasbackendv1alpha1.K8sAsBackend{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      name,
// 			Namespace: namespace,
// 		},
// 		Spec: k8sasbackendv1alpha1.K8sAsBackendSpec{
// 			Size: replicas, // Set desired number of Memcached replicas.
// 		},
// 	}

// 	// Objects to track in the fake client.
// 	objs := []runtime.Object{
// 		i,
// 	}

// 	// Register operator types with the runtime scheme.
// 	s := scheme.Scheme
// 	apiextensionsv1beta1.AddToScheme(s)
// 	s.AddKnownTypes(k8sasbackendv1alpha1.SchemeGroupVersion, i)

// 	// Create a fake client to mock API calls.
// 	cl := fake.NewFakeClient(objs...)
// 	clCert := &certfake.FakeCertificateSigningRequests{
// 		Fake: &certfake.FakeCertificatesV1beta1{
// 			Fake: &clientgotesting.Fake{
// 				ReactionChain: []clientgotesting.Reactor{
// 					CertApprovalReactor{},
// 				},
// 			},
// 		},
// 	}
// 	// Create a ReconcileMemcached object with the scheme and fake client.
// 	r := &ReconcileK8sAsBackend{client: cl, scheme: s, certClient: clCert}

// 	// todoCrd := r.createTodoCrd("todos.k8sasbackend.com", operatorInstance)
// 	// s.AddKnownTypes(apiextensionsv1beta1.SchemeGroupVersion, todoCrd)

// 	// Mock request to simulate Reconcile() being called on an event for a
// 	// watched resource .
// 	req := reconcile.Request{
// 		NamespacedName: types.NamespacedName{
// 			Name:      name,
// 			Namespace: namespace,
// 		},
// 	}

// 	certFactory := &CertFactory{}

// 	res, err := certFactory.ensure(r, req, i)
// 	assert.NoError(t, err)
// 	assert.Nil(t, res)

// }

// type CertApprovalReactor struct{}

// // Handles indicates whether or not this Reactor deals with a given
// // action.
// func (r CertApprovalReactor) Handles(action clientgotesting.Action) bool {
// 	// if action.Subresource == "approval" {
// 	// 	return true
// 	// }

// 	return true
// }

// // React handles the action and returns results.  It may choose to
// // delegate by indicated handled=false.
// func (r CertApprovalReactor) React(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
// 	cert, _ := b64.StdEncoding.DecodeString(certBase64String)
// 	updAction := action.(clientgotesting.UpdateAction)
// 	obj := updAction.GetObject()
// 	originalResource := obj.(*apicert.CertificateSigningRequest)
// 	originalResource.Status.Certificate = cert //make([]byte, 128)
// 	return true, originalResource.DeepCopy(), nil
// }

// func TestEnsureDeployment(t *testing.T) {
// 	// Set the logger to development mode for verbose logs.
// 	logf.SetLogger(logf.ZapLogger(true))

// 	var (
// 		name            = "memcached-operator"
// 		namespace       = "memcached"
// 		replicas  int32 = 3
// 	)

// 	operatorInstance := &k8sasbackendv1alpha1.K8sAsBackend{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      name,
// 			Namespace: namespace,
// 		},
// 		Spec: k8sasbackendv1alpha1.K8sAsBackendSpec{
// 			Size: replicas, // Set desired number of Memcached replicas.
// 		},
// 	}

// 	// Objects to track in the fake client.
// 	objs := []runtime.Object{
// 		operatorInstance,
// 	}

// 	// Register operator types with the runtime scheme.
// 	s := scheme.Scheme
// 	apiextensionsv1beta1.AddToScheme(s)
// 	s.AddKnownTypes(k8sasbackendv1alpha1.SchemeGroupVersion, operatorInstance)

// 	// Create a fake client to mock API calls.
// 	cl := fake.NewFakeClient(objs...)
// 	// Create a ReconcileMemcached object with the scheme and fake client.
// 	r := &ReconcileK8sAsBackend{client: cl, scheme: s}

// 	// todoCrd := r.createTodoCrd("todos.k8sasbackend.com", operatorInstance)
// 	// s.AddKnownTypes(apiextensionsv1beta1.SchemeGroupVersion, todoCrd)

// 	// Mock request to simulate Reconcile() being called on an event for a
// 	// watched resource .
// 	req := reconcile.Request{
// 		NamespacedName: types.NamespacedName{
// 			Name:      name,
// 			Namespace: namespace,
// 		},
// 	}
// 	// res, err := r.Reconcile(req)
// 	// if err != nil {
// 	// 	t.Fatalf("reconcile: (%v)", err)
// 	// }
// 	// // Check the result of reconciliation to make sure it has the desired state.
// 	// if !res.Requeue {
// 	// 	t.Error("reconcile did not requeue request as expected")
// 	// }

// 	//dep := &appsv1.Deployment{}

// 	//r.registerCrd()

// 	// res, err := r.ensureResource(req, operatorInstance, "todos.k8sasbackend.com", r.createTodoCrd)
// 	// verify(t, res, err)

// 	// res, err = r.ensureResource(req, operatorInstance, "todo-crd", r.createAccount)
// 	// verify(t, res, err)

// 	res, err := r.ensureResource(req, operatorInstance, createNamespacedName("todo-crd", operatorInstance.Namespace), &RoleFactory{})
// 	verify(t, res, err)

// 	res, err = r.ensureResource(req, operatorInstance, createNamespacedName("todo-crd", operatorInstance.Namespace), &AccountFactory{})
// 	verify(t, res, err)

// }

// func verify(t *testing.T, res *reconcile.Result, err error) {
// 	if err != nil {
// 		t.Errorf("get memcached: (%v)", err)
// 	}

// 	if res != nil {
// 		t.Error("reconcile did not return a nil Result")
// 	}
// }
