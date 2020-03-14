package webhookserver

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"os"

	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	common "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	v1beta1 "k8s.io/api/certificates/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (ws WebhookServer) ensureSecret(i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {
	//_, err := ws.createCertIfNeeded()
	requeue, err := ws.createCertIfNeeded()
	if err != nil {
		return &reconcile.Result{}, err
	} else if requeue {
		common.Log.Info("createCertIfNeeded request a requeue")
		return &reconcile.Result{Requeue: true}, err
	}

	resUtils := &common.ResourceUtils{
		Scheme:          ws.Scheme,
		Client:          ws.Client,
		ResourceFactory: ws.createSecret,
	}

	found := &corev1.Secret{}
	err = common.EnsureResource(found, secretName, i, resUtils)
	common.Log.Info("ensureSecret returns a null result")
	return nil, err
}

func (ws WebhookServer) createCertIfNeeded() (requeue bool, err error) {
	if common.FileNotExists(ws.CerFilePath) {

		//found := &v1beta1.CertificateSigningRequest{}
		//nsName := types.NamespacedName{Name: csrName, Namespace: ""}
		//err := ws.Client.Get(context.TODO(), nsName, found)

		found, err := ws.CertClient.Get(csrName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			certRequest := ws.createKeyAndCertRequestAsPem()
			err := ws.createSigningRequest(certRequest)
			if err != nil {
				panic(err)
			}
			err = ws.approveCertificate()
			if err != nil {
				panic(err)
			}

			log.Info("Requeue - waiting for certificate provisioning")
			return true, err
		} else if err != nil {
			// Error that isn't due to the resource not existing
			//log.Error(err, fmt.Sprintf("Failed to get %s", common.GetKind(found)))
			return false, err
		}

		if len(found.Status.Certificate) == 0 {
			return true, err
		}

		err = ws.storeCertificateAsPem(found.Status.Certificate)
		if err != nil {
			return false, err
		}
		ws.CertClient.Delete(found.Name, &metav1.DeleteOptions{})
	}

	return false, nil

}

func (ws WebhookServer) createSecret(resourceName string, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {
	nsName := types.NamespacedName{Name: resourceName, Namespace: i.Namespace}
	cert, _ := ioutil.ReadFile(ws.CerFilePath)
	key, _ := ioutil.ReadFile(ws.KeyFilePath)
	return common.CreateSecret(nsName, map[string][]byte{
		"key.pem":  key,
		"cert.pem": cert,
	})
}

func (ws WebhookServer) approveCertificate() error {

	// found := &v1beta1.CertificateSigningRequest{}
	// nsName := types.NamespacedName{Name: csrName, Namespace: ""}
	// err := ws.Client.Get(context.TODO(), nsName, found)
	found, err := ws.CertClient.Get(csrName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	found.Status.Conditions = []v1beta1.CertificateSigningRequestCondition{{
		LastUpdateTime: metav1.Now(),
		Message:        "This CSR was approved by k8s-as-backend operator.",
		Reason:         "OperatorApprove",
		Type:           v1beta1.CertificateApproved,
	}}
	resp, err := ws.CertClient.UpdateApproval(found)
	if err != nil {
		panic(err)
	}

	log.Info("UpdateApproval", "ContentLength", len(resp.Status.Certificate))
	return err
}

func (ws WebhookServer) storeCertificateAsPem(cert []byte) error {
	err := ioutil.WriteFile(ws.CerFilePath, cert, 0644)
	if err != nil {
		panic(err)
	}
	return err
}

func (ws WebhookServer) createKeyAndCertRequestAsPem() (csrPem []byte) {
	keyBytes, _ := rsa.GenerateKey(rand.Reader, 2048)
	keyfile, _ := os.Create(ws.KeyFilePath)
	defer keyfile.Close()
	pem.Encode(keyfile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keyBytes)})

	subj := pkix.Name{
		CommonName: "admission-webhook-example-svc.default.svc",
		// Country:            []string{"AU"},
		// Province:           []string{"Some-State"},
		// Locality:           []string{"MyCity"},
		// Organization:       []string{"Company Ltd"},
		// OrganizationalUnit: []string{"IT"},
		// ExtraNames: []pkix.AttributeTypeAndValue{
		// 	{Type: oidEmailAddress, Value: emailAddress},
		// },
	}

	template := x509.CertificateRequest{
		Subject: subj,
		//EmailAddresses:     []string{emailAddress},
		SignatureAlgorithm: x509.SHA256WithRSA,
		DNSNames: []string{
			"service=admission-webhook-example-svc",
			"service=admission-webhook-example-svc.default",
			"service=admission-webhook-example-svc.default.svc"},
	}

	csr, _ := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)
	//csrFile, _ := os.Create("/Users/cristiano/Coding/golang/k8s-as-backend/operator/certs/server.csr")
	//openssl req -in ~/Coding/golang/k8s-as-backend/operator/server-test.csr -noout -text
	//defer csrFile.Close()
	csrPem = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csr}) //os.Stdout
	return
}

func (ws WebhookServer) createSigningRequest(request []byte) error {
	// name := "admission-webhook-example-svc.default"

	res := &v1beta1.CertificateSigningRequest{
		ObjectMeta: common.CreateMeta(csrName, ""),
		Spec: v1beta1.CertificateSigningRequestSpec{
			Groups:  []string{"system:authenticated"},
			Request: request,
			Usages:  []v1beta1.KeyUsage{"digital signature", "key encipherment", "server auth"},
		},
	}

	//return ws.Client.Create(context.TODO(), res)
	_, err := ws.CertClient.Create(res)
	return err
}
