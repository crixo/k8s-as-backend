package webhookserver

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path"

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
	sharedResourceName := common.CreateUniqueSecondaryResourceName(i, baseName)
	cerFilePath := getPemName(i, "cert")
	keyFilePath := getPemName(i, "key")
	requeue, err := ws.createCertIfNeeded(cerFilePath, keyFilePath, sharedResourceName, i.Namespace)
	if err != nil {
		return &reconcile.Result{}, err
	} else if requeue {
		common.Log.Info("createCertIfNeeded request a requeue")
		return &reconcile.Result{Requeue: true}, err
	}

	i.Status.AdmissionWebhookPems = []string{
		cerFilePath,
		keyFilePath,
	}
	err = ws.Client.Status().Update(context.TODO(), i)
	if err != nil {
		return &reconcile.Result{}, err
	}

	//TODO: if certs have been recreated and secret already exists, delete it

	resUtils := &common.ResourceUtils{
		Scheme:          ws.Scheme,
		Client:          ws.Client,
		ResourceFactory: ws.createSecret,
	}

	secretName := common.CreateUniqueSecondaryResourceName(i, baseName)
	nsn := types.NamespacedName{Name: secretName, Namespace: i.Namespace}
	found := &corev1.Secret{}
	err = common.EnsureResource(found, nsn, i, resUtils)
	common.Log.Info("ensureSecret returns a null result")
	return nil, err
}

func getPemName(i *k8sasbackendv1alpha1.K8sAsBackend, pemKind string) string {
	return path.Join(pemFolder, fmt.Sprintf("%s_%s_%s.pem", i.Name, i.Namespace, pemKind))
}

func (ws WebhookServer) shouldCreatePems(cerFilePath string) bool {

	if common.FileNotExists(cerFilePath) { //TODO: add cert vs cluster CA verification
		return true
	}

	return false
}

func (ws WebhookServer) createCertIfNeeded(cerFilePath, keyFilePath, sharedResourceName, namespace string) (requeue bool, err error) {

	if ws.shouldCreatePems(cerFilePath) {

		if common.FileNotExists(keyFilePath) {
			createKeyAsPem(keyFilePath)
		}

		//found := &v1beta1.CertificateSigningRequest{}
		//nsName := types.NamespacedName{Name: csrName, Namespace: ""}
		//err := ws.Client.Get(context.TODO(), nsName, found)

		found, err := ws.CertClient.Get(sharedResourceName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			certRequest := ws.createKeyAndCertRequestAsPem(keyFilePath, sharedResourceName, namespace)
			err := ws.createSigningRequest(sharedResourceName, certRequest)
			if err != nil {
				panic(err)
			}
			err = ws.approveCertificate(sharedResourceName)
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

		err = ws.storeCertificateAsPem(cerFilePath, found.Status.Certificate)
		if err != nil {
			return false, err
		}
		ws.CertClient.Delete(found.Name, &metav1.DeleteOptions{})
	}

	return false, nil

}

func (ws WebhookServer) createSecret(nsName types.NamespacedName, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {
	cerFilePath := getPemName(i, "cert")
	keyFilePath := getPemName(i, "key")
	cert, _ := ioutil.ReadFile(cerFilePath)
	key, _ := ioutil.ReadFile(keyFilePath)
	return common.CreateSecret(nsName, map[string][]byte{
		"key.pem":  key,
		"cert.pem": cert,
	})
}

func (ws WebhookServer) approveCertificate(csrName string) error {

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

func (ws WebhookServer) storeCertificateAsPem(cerFilePath string, cert []byte) error {
	err := ioutil.WriteFile(cerFilePath, cert, 0644)
	if err != nil {
		panic(err)
	}
	return err
}

func createKeyAsPem(keyFilePath string) {
	rsaKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	keyfile, _ := os.Create(keyFilePath)
	defer keyfile.Close()
	pem.Encode(keyfile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(rsaKey)})

	return
}

func (ws WebhookServer) createKeyAndCertRequestAsPem(keyFilePath, webhookServerServiceName, namespace string) (csrPem []byte) {

	fileBytes, _ := ioutil.ReadFile(keyFilePath)
	privPem, _ := pem.Decode(fileBytes)
	parsedKey, _ := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	//keyBytes := parsedKey

	subj := pkix.Name{
		CommonName: fmt.Sprintf("%s.%s.svc", webhookServerServiceName, namespace),
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
			webhookServerServiceName,
			fmt.Sprintf("%s.%s", webhookServerServiceName, namespace),
			fmt.Sprintf("%s.%s.svc", webhookServerServiceName, namespace),
		},
	}

	//csr, _ := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)
	csr, _ := x509.CreateCertificateRequest(rand.Reader, &template, parsedKey)
	//csrFile, _ := os.Create("/Users/cristiano/Coding/golang/k8s-as-backend/operator/certs/server.csr")
	//openssl req -in ~/Coding/golang/k8s-as-backend/operator/server-test.csr -noout -text
	//defer csrFile.Close()
	csrPem = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csr}) //os.Stdout
	return
}

func (ws WebhookServer) createSigningRequest(csrName string, request []byte) error {
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

func verifyCert(rootPEM, certPEM []byte, name string) error {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(rootPEM)
	if !ok {
		return fmt.Errorf("failed to parse root certificate")
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		return fmt.Errorf("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %v", err.Error())
	}

	opts := x509.VerifyOptions{
		//DNSName: name,
		Roots: roots,
	}

	if _, err := cert.Verify(opts); err != nil {
		return fmt.Errorf("failed to verify certificate: %v", err.Error())
	}

	return nil
}
