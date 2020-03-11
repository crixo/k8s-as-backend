package k8sasbackend

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"

	v1beta1 "k8s.io/api/certificates/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type CertFactory struct{}

func (f CertFactory) createEmpty() runtime.Object {
	return &v1beta1.CertificateSigningRequest{}
}

func (f CertFactory) getNames() []string {
	return []string{"admission-webhook-example-svc"}
}

func (f CertFactory) ensure(r *ReconcileK8sAsBackend,
	request reconcile.Request,
	i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {

	resourceName := fmt.Sprintf("%s.%s", f.getNames()[0], i.Namespace)

	// certificateSigningRequest, err := r.certClient.Get(resourceName, v1.GetOptions{})
	// log.Info("certificateSigningRequest", "req", certificateSigningRequest)
	// log.Error(err, "getting certificateSigningRequest")

	result, err := r.ensureResource(request, i, createNamespacedName(resourceName, ""), f)
	if result != nil {
		return result, err
	}

	found := &v1beta1.CertificateSigningRequest{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      resourceName,
		Namespace: "",
	}, found)
	if err != nil {
		return &reconcile.Result{}, err
	}

	log.Info("check found.Status.Certificate before approving", "Certificate", found.Status.Certificate)
	//emptyByteVar := make([]byte, 128)
	if found.Status.Certificate == nil { //|| bytes.Equal(found.Status.Certificate, emptyByteVar)
		found.Status.Conditions = []v1beta1.CertificateSigningRequestCondition{{
			LastUpdateTime: metav1.Now(),
			Message:        "This CSR was approved by k8s-as-backend operator.",
			Reason:         "OperatorApprove",
			Type:           v1beta1.CertificateApproved,
		}}
		resp, err := r.certClient.UpdateApproval(found)
		//log.Info("Approval response", "CertificateSigningRequest", resp)
		if err != nil {
			return &reconcile.Result{}, err
		}

		if resp.Status.Certificate == nil {
			log.Info("Requeuing waiting for approvale", "Certificate", found.Status.Certificate)
			return &reconcile.Result{Requeue: true}, err
		}

		found = resp
	}

	cert := found.Status.Certificate
	err = ioutil.WriteFile("/Users/cristiano/Coding/golang/k8s-as-backend/operator/server-cert.pem", cert, 0644)
	if err != nil {
		log.Error(err, "Unable to write cert to file")
		return &reconcile.Result{}, err
	}

	return nil, nil

}

func (f CertFactory) create(name string, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {
	//$(cat ${tmpdir}/server.csr | base64 | tr -d '\n')
	// https://golang.org/pkg/crypto/x509/#CertificateRequest
	//x509.ParseCertificateRequest(block.Bytes)
	//csrAsBase64String := f.readSigningRequest()
	content, _ := ioutil.ReadFile("/Users/cristiano/Coding/golang/k8s-as-backend/operator/server.csr")
	return &v1beta1.CertificateSigningRequest{
		ObjectMeta: createMeta(name, ""),
		Spec: v1beta1.CertificateSigningRequestSpec{
			Groups:  []string{"system:authenticated"},
			Request: content,
			Usages:  []v1beta1.KeyUsage{"digital signature", "key encipherment", "server auth"},
		},
	}
}

func (f CertFactory) readSigningRequest() string {
	content, err := ioutil.ReadFile("/Users/cristiano/Coding/golang/k8s-as-backend/operator/server.csr")
	if err != nil {
		log.Error(err, "Unable to read signign request cert file")
	}

	// Convert []byte to string and print to screen
	base64text := b64.StdEncoding.EncodeToString([]byte(content))
	res := strings.Replace(base64text, "\n", "", -1)
	//fmt.Println(res)

	return res
}
