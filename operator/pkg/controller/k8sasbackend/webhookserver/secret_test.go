package webhookserver

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"

	b64 "encoding/base64"
	"encoding/hex"

	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	common "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	"github.com/stretchr/testify/assert"
	apicert "k8s.io/api/certificates/v1beta1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	certfake "k8s.io/client-go/kubernetes/typed/certificates/v1beta1/fake"

	clientgotesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const certBase64String = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURxVENDQXBHZ0F3SUJBZ0lVWUxLWFVlRDcrRjRURnE4dkJ4VXpvano1UUhVd0RRWUpLb1pJaHZjTkFRRUwKQlFBd0ZURVRNQkVHQTFVRUF4TUthM1ZpWlhKdVpYUmxjekFlRncweU1EQXpNVEV4T1RRNE1EQmFGdzB5TVRBegpNVEV4T1RRNE1EQmFNRFF4TWpBd0JnTlZCQU1US1dGa2JXbHpjMmx2YmkxM1pXSm9iMjlyTFdWNFlXMXdiR1V0CmMzWmpMbVJsWm1GMWJIUXVjM1pqTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEKNlFNR1c0aFZFd3dGN2NCdEVqaE02aGUzSThQNDc5ZVQraXhCQWpXcjR5L2dkWURXeFprNG5kWnY4SHVKY0JWTQpDNUEwSG5jZXNlaFMrTmJ5a2dxN1lDd01LM3FDUlJPMVpKa1Y3d1JvSDhIOU91ZmxsczdWY0hQTDMyeU5jVnQxCnVGWFlDUUpKVUYxUjRXczJVNnFvUlVFN040T1plQXRWcnM2YmpVZzZuSFR0ZWtUUUt5ekxHMWV6UG5rNWRtNFEKdTlXVFVBa0lSTm8vemQrTkJkaDYwdy96UHVTZkRDcFV0U1dkYWFQUnVjb3hTV0ZVWDROWEYzNUhoVW91OTk3QgpoZWl6TEhiVndiK2Fqc3I1UEVJL2ZOd1RwaWtoNnZiVGkvRDZOMGRVdW94enBDSlVtSnkwSjFKc0pDSnVZeHAxCjNjeHpUbjdMUEdBcHVnR3JnWEZ1NlFJREFRQUJvNEhSTUlIT01BNEdBMVVkRHdFQi93UUVBd0lGb0RBVEJnTlYKSFNVRUREQUtCZ2dyQmdFRkJRY0RBVEFNQmdOVkhSTUJBZjhFQWpBQU1CMEdBMVVkRGdRV0JCU3VKUWJVcWhacwpYdS9KK3dGb3RpM0toOEpUSkRCNkJnTlZIUkVFY3pCeGdoMWhaRzFwYzNOcGIyNHRkMlZpYUc5dmF5MWxlR0Z0CmNHeGxMWE4yWTRJbFlXUnRhWE56YVc5dUxYZGxZbWh2YjJzdFpYaGhiWEJzWlMxemRtTXVaR1ZtWVhWc2RJSXAKWVdSdGFYTnphVzl1TFhkbFltaHZiMnN0WlhoaGJYQnNaUzF6ZG1NdVpHVm1ZWFZzZEM1emRtTXdEUVlKS29aSQpodmNOQVFFTEJRQURnZ0VCQUJBU0ZQT05JQTBkV3Z3ekF2NmJJN2Y1clhmVWtWZVkwZUIxZmx0UTlndUM0Z3JHCldsZ0gvOGxHQThCVjVNUUxTaWsxV1FiVmtackYySDVSNEFlaUVqYkxHQlRJc2doYXROSFZLZitjdEprbGhjRysKVXc5bk5HSnVhWWdvNkx5NmRBS0R5YXNuWmlHNDNKL3NOM0tHQzJPRktLUGFHK0pYWHFYbWxwU3RNbWV1OVpsagp0R2RaVDFiUktNeWlCa3cwMDd2R2pENjMzSzFBU2pLbHRDeEFwY2FZUVlkeTY5RC9kUHgweU5XcmdkUGlZd2hRClBwWThqOFVpMExaNGd1NlRtYXRPb0xGSW84OGxwUGt3RVpOZUthbENuRGNLWW9OKzByT0o2QTRzZExJVmNZNm8KSkJBc1NPSFpSa29qNkNFcG45MGF1SVdqSmtkdnM5WVNhLzNUUC9BPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="

var cert, _ = b64.StdEncoding.DecodeString(certBase64String)

func TestSecret(t *testing.T) {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))
	common.Log.Info("TestSecret")
	var (
		name            = "k8sasbackend-operator"
		namespace       = "default"
		replicas  int32 = 3
	)

	i := &k8sasbackendv1alpha1.K8sAsBackend{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: k8sasbackendv1alpha1.K8sAsBackendSpec{
			Size: replicas,
		},
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{
		i,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	apiextensionsv1beta1.AddToScheme(s)
	s.AddKnownTypes(k8sasbackendv1alpha1.SchemeGroupVersion, i)

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)

	// certificate fake client
	clCert := &certfake.FakeCertificateSigningRequests{
		Fake: &certfake.FakeCertificatesV1beta1{
			Fake: &clientgotesting.Fake{
				ReactionChain: []clientgotesting.Reactor{
					&CertApprovalReactor{},
				},
			},
		},
	}

	// todoCrd := r.createTodoCrd("todos.k8sasbackend.com", operatorInstance)
	// s.AddKnownTypes(apiextensionsv1beta1.SchemeGroupVersion, todoCrd)
	certFile := TempFileName("/Users/cristiano/Coding/golang/k8s-as-backend/operator/certs", "server-cert", ".pem")
	defer os.Remove(certFile)
	keyFile := TempFileName("/Users/cristiano/Coding/golang/k8s-as-backend/operator/certs", "server-key", ".pem")
	defer os.Remove(keyFile)
	webhookServer := &WebhookServer{
		CerFilePath: certFile,
		KeyFilePath: keyFile,
		Client:      cl,
		Scheme:      s,
		CertClient:  clCert,
	}

	res, err := webhookServer.ensureSecret(i)
	assert.NoError(t, err)
	assert.True(t, assert.ObjectsAreEqualValues(res, &reconcile.Result{Requeue: true}),
		"After first iteration from scratch res should not be null due to cert aproval reschedule")

	res, err = webhookServer.ensureSecret(i)
	assert.NoError(t, err)
	assert.Nil(t, res)
}

type CertApprovalReactor struct {
	resource *apicert.CertificateSigningRequest
}

// Handles indicates whether or not this Reactor deals with a given
// action.
func (r CertApprovalReactor) Handles(action clientgotesting.Action) bool {
	// if action.Subresource == "approval" {
	// 	return true
	// }

	return true
}

// React handles the action and returns results.  It may choose to
// delegate by indicated handled=false.
func (r *CertApprovalReactor) React(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
	switch v := action.(type) {
	case clientgotesting.GetActionImpl:
		err = nil
		if r.resource == nil {
			err = &metav1apierr.StatusError{ErrStatus: metav1.Status{Reason: metav1.StatusReasonNotFound}}
		}
		return true, r.resource, err
	case clientgotesting.CreateActionImpl:
		r.resource = v.GetObject().(*apicert.CertificateSigningRequest)
		return true, r.resource, nil
	case clientgotesting.UpdateActionImpl:
		r.resource = v.GetObject().(*apicert.CertificateSigningRequest)
		r.resource.Status.Certificate = cert //make([]byte, 128)
		return true, r.resource, nil
	case clientgotesting.DeleteActionImpl:
		r.resource = nil
		return true, r.resource, nil
	default:
		return false, nil, nil

	}

	// updAction := action.(clientgotesting.UpdateAction)
	// obj := updAction.GetObject()
	// originalResource := obj.(*apicert.CertificateSigningRequest)
	// originalResource.Status.Certificate = cert //make([]byte, 128)
	// return true, originalResource.DeepCopy(), nil
}

func TempFileName(dir, prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	if len(dir) == 0 {
		dir = os.TempDir()
	}
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}
