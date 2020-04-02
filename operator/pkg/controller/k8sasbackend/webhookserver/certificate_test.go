package webhookserver

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"os"
	"testing"
)

//var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}

func TestCertValidity(t *testing.T) {
	rootPem, err := ioutil.ReadFile("../../../../cacert.pem")
	if err != nil {
		t.Error(err)
	}
	certPem, _ := ioutil.ReadFile("../../../../kab01-todos-webhook-server_cert.pem")
	err = verifyCert(rootPem, certPem, "")
	if err != nil {
		t.Error(err)
	}
}

func TestPemCreation(t *testing.T) {
	keyBytes, _ := rsa.GenerateKey(rand.Reader, 2048)

	pem.Encode(os.Stdout, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keyBytes)})

	//emailAddress := "test@example.com"
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
		// KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		// BasicConstraintsValid: true,
		DNSNames: []string{"service=admission-webhook-example-svc", "service=admission-webhook-example-svc.default", "service=admission-webhook-example-svc.default.svc"},
		// Extensions: []pkix.Extension{
		// 	{}
		// },
	}

	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)
	// file, _ := os.Create("/Users/cristiano/Coding/golang/k8s-as-backend/operator/server-test.csr")
	// //openssl req -in ~/Coding/golang/k8s-as-backend/operator/server-test.csr -noout -text
	// defer file.Close()
	pem.Encode(os.Stdout, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}) //os.Stdout

}
