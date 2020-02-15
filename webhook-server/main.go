/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	//"github.com/spf13/cobra"

	v1 "k8s.io/api/admission/v1"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"

	// TODO: try this library to see if it generates correct json patch
	// https://github.com/mattbaird/jsonpatch

	//crd.go
	// apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	// apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	admissionv1 "k8s.io/api/admission/v1"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	//"github.com/crixo/k8s-as-backend/webhook-server/main"

	"go.uber.org/zap"
)

var (
	certFile string
	keyFile  string
	port     int
)

var (
	logger = zap.NewExample()
)

// CmdWebhook is used by agnhost Cobra.
// var CmdWebhook = &cobra.Command{
// 	Use:   "webhook",
// 	Short: "Starts a HTTP server, useful for testing MutatingAdmissionWebhook and ValidatingAdmissionWebhook",
// 	Long: `Starts a HTTP server, useful for testing MutatingAdmissionWebhook and ValidatingAdmissionWebhook.
// After deploying it to Kubernetes cluster, the Administrator needs to create a ValidatingWebhookConfiguration
// in the Kubernetes cluster to register remote webhook admission controllers.`,
// 	Args: cobra.MaximumNArgs(0),
// 	Run:  main,
// }

func init() {
	// CmdWebhook.Flags().StringVar(&certFile, "tls-cert-file", "",
	// 	"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated after server cert).")
	// CmdWebhook.Flags().StringVar(&keyFile, "tls-private-key-file", "",
	// 	"File containing the default x509 private key matching --tls-cert-file.")
	// CmdWebhook.Flags().IntVar(&port, "port", 443,
	// 	"Secure port that the webhook listens on")
	port = 443
	logger.Info("klog.SetOutput(os.Stdout)")
	klog.SetOutput(os.Stdout)
}

// admitv1beta1Func handles a v1beta1 admission
type admitv1beta1Func func(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

// admitv1beta1Func handles a v1 admission
type admitv1Func func(v1.AdmissionReview) *v1.AdmissionResponse

// admitHandler is a handler, for both validators and mutators, that supports multiple admission review versions
type admitHandler struct {
	v1beta1 admitv1beta1Func
	v1      admitv1Func
}

func newDelegateToV1AdmitHandler(f admitv1Func) admitHandler {
	return admitHandler{
		v1beta1: delegateV1beta1AdmitToV1(f),
		v1:      f,
	}
}

func delegateV1beta1AdmitToV1(f admitv1Func) admitv1beta1Func {
	return func(review v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
		in := v1.AdmissionReview{Request: convertAdmissionRequestToV1(review.Request)}
		out := f(in)
		return convertAdmissionResponseToV1beta1(out)
	}
}

// serve handles the http portion of a request prior to handing to an admit
// function
func serve(w http.ResponseWriter, r *http.Request, admit admitHandler) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		klog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	klog.V(2).Info(fmt.Sprintf("handling request: %s", body))

	deserializer := codecs.UniversalDeserializer()
	obj, gvk, err := deserializer.Decode(body, nil, nil)
	if err != nil {
		msg := fmt.Sprintf("Request could not be decoded: %v", err)
		klog.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var responseObj runtime.Object
	switch *gvk {
	case v1beta1.SchemeGroupVersion.WithKind("AdmissionReview"):
		requestedAdmissionReview, ok := obj.(*v1beta1.AdmissionReview)
		if !ok {
			klog.Errorf("Expected v1beta1.AdmissionReview but got: %T", obj)
			return
		}
		responseAdmissionReview := &v1beta1.AdmissionReview{}
		responseAdmissionReview.SetGroupVersionKind(*gvk)
		responseAdmissionReview.Response = admit.v1beta1(*requestedAdmissionReview)
		responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID
		responseObj = responseAdmissionReview
	case v1.SchemeGroupVersion.WithKind("AdmissionReview"):
		requestedAdmissionReview, ok := obj.(*v1.AdmissionReview)
		if !ok {
			klog.Errorf("Expected v1.AdmissionReview but got: %T", obj)
			return
		}
		responseAdmissionReview := &v1.AdmissionReview{}
		responseAdmissionReview.SetGroupVersionKind(*gvk)
		responseAdmissionReview.Response = admit.v1(*requestedAdmissionReview)
		responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID
		responseObj = responseAdmissionReview
	default:
		msg := fmt.Sprintf("Unsupported group version kind: %v", gvk)
		klog.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	klog.V(2).Info(fmt.Sprintf("sending response: %v", responseObj))
	respBytes, err := json.Marshal(responseObj)
	if err != nil {
		klog.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respBytes); err != nil {
		klog.Error(err)
	}
}

func serveCRD(w http.ResponseWriter, r *http.Request) {
	serve(w, r, newDelegateToV1AdmitHandler(admitCRD))
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func main() {//cmd *cobra.Command, args []string
	config := Config{
		CertFile: "/etc/webhook/certs/cert.pem",//certFile,
		KeyFile:  "/etc/webhook/certs/key.pem",//keyFile,
	}

	http.HandleFunc("/crd", serveCRD)
	http.HandleFunc("/test", test)
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		TLSConfig: configTLS(config),
	}
	err := server.ListenAndServeTLS("", "")
	if err != nil {
		panic(err)
	}
}

//crd.go
// This function expects all CRDs submitted to it to be apiextensions.k8s.io/v1beta1 or apiextensions.k8s.io/v1.
// func admitCRD(ar v1.AdmissionReview) *v1.AdmissionResponse {
// 	klog.V(2).Info("admitting crd")

// 	resource := "customresourcedefinitions"
// 	v1beta1GVR := metav1.GroupVersionResource{Group: apiextensionsv1beta1.GroupName, Version: "v1beta1", Resource: resource}
// 	v1GVR := metav1.GroupVersionResource{Group: apiextensionsv1.GroupName, Version: "v1", Resource: resource}

// 	reviewResponse := v1.AdmissionResponse{}
// 	reviewResponse.Allowed = true

// 	raw := ar.Request.Object.Raw
// 	var labels map[string]string

// 	switch ar.Request.Resource {
// 	case v1beta1GVR:
// 		crd := apiextensionsv1beta1.CustomResourceDefinition{}
// 		deserializer := codecs.UniversalDeserializer()
// 		if _, _, err := deserializer.Decode(raw, nil, &crd); err != nil {
// 			klog.Error(err)
// 			return toV1AdmissionResponse(err)
// 		}
// 		labels = crd.Labels
// 	case v1GVR:
// 		crd := apiextensionsv1.CustomResourceDefinition{}
// 		deserializer := codecs.UniversalDeserializer()
// 		if _, _, err := deserializer.Decode(raw, nil, &crd); err != nil {
// 			klog.Error(err)
// 			return toV1AdmissionResponse(err)
// 		}
// 		labels = crd.Labels
// 	default:
// 		err := fmt.Errorf("expect resource to be one of [%v, %v] but got %v", v1beta1GVR, v1GVR, ar.Request.Resource)
// 		klog.Error(err)
// 		return toV1AdmissionResponse(err)
// 	}

// 	if v, ok := labels["webhook-e2e-test"]; ok {
// 		if v == "webhook-disallow" {
// 			reviewResponse.Allowed = false
// 			reviewResponse.Result = &metav1.Status{Message: "the crd contains unwanted label"}
// 		}
// 	}
// 	return &reviewResponse

// }

// convert.go
func convertAdmissionRequestToV1(r *v1beta1.AdmissionRequest) *v1.AdmissionRequest {
	return &v1.AdmissionRequest{
		Kind:               r.Kind,
		Namespace:          r.Namespace,
		Name:               r.Name,
		Object:             r.Object,
		Resource:           r.Resource,
		Operation:          v1.Operation(r.Operation),
		UID:                r.UID,
		DryRun:             r.DryRun,
		OldObject:          r.OldObject,
		Options:            r.Options,
		RequestKind:        r.RequestKind,
		RequestResource:    r.RequestResource,
		RequestSubResource: r.RequestSubResource,
		SubResource:        r.SubResource,
		UserInfo:           r.UserInfo,
	}
}

func convertAdmissionRequestToV1beta1(r *v1.AdmissionRequest) *v1beta1.AdmissionRequest {
	return &v1beta1.AdmissionRequest{
		Kind:               r.Kind,
		Namespace:          r.Namespace,
		Name:               r.Name,
		Object:             r.Object,
		Resource:           r.Resource,
		Operation:          v1beta1.Operation(r.Operation),
		UID:                r.UID,
		DryRun:             r.DryRun,
		OldObject:          r.OldObject,
		Options:            r.Options,
		RequestKind:        r.RequestKind,
		RequestResource:    r.RequestResource,
		RequestSubResource: r.RequestSubResource,
		SubResource:        r.SubResource,
		UserInfo:           r.UserInfo,
	}
}

func convertAdmissionResponseToV1(r *v1beta1.AdmissionResponse) *v1.AdmissionResponse {
	var pt *v1.PatchType
	if r.PatchType != nil {
		t := v1.PatchType(*r.PatchType)
		pt = &t
	}
	return &v1.AdmissionResponse{
		UID:              r.UID,
		Allowed:          r.Allowed,
		AuditAnnotations: r.AuditAnnotations,
		Patch:            r.Patch,
		PatchType:        pt,
		Result:           r.Result,
	}
}

func convertAdmissionResponseToV1beta1(r *v1.AdmissionResponse) *v1beta1.AdmissionResponse {
	var pt *v1beta1.PatchType
	if r.PatchType != nil {
		t := v1beta1.PatchType(*r.PatchType)
		pt = &t
	}
	return &v1beta1.AdmissionResponse{
		UID:              r.UID,
		Allowed:          r.Allowed,
		AuditAnnotations: r.AuditAnnotations,
		Patch:            r.Patch,
		PatchType:        pt,
		Result:           r.Result,
	}
}

func toV1AdmissionResponse(err error) *v1.AdmissionResponse {
	return &v1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

//scheme.go
var scheme = runtime.NewScheme()
var codecs = serializer.NewCodecFactory(scheme)

func init() {
	addToScheme(scheme)
}

func addToScheme(scheme *runtime.Scheme) {
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(admissionv1beta1.AddToScheme(scheme))
	utilruntime.Must(admissionregistrationv1beta1.AddToScheme(scheme))
	utilruntime.Must(admissionv1.AddToScheme(scheme))
	utilruntime.Must(admissionregistrationv1.AddToScheme(scheme))
}