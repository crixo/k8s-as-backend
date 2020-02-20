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
	"bytes"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	v1 "k8s.io/api/admission/v1"
	"k8s.io/klog"
)

// var (
// 	logger = zap.NewExample()
// )


// This function expects all CRDs submitted to it to be apiextensions.k8s.io/v1beta1 or apiextensions.k8s.io/v1.
func admitCRD(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.V(2).Info("admitting crd")

	//resource := "customresourcedefinitions"
	// v1beta1GVR := metav1.GroupVersionResource{Group: apiextensionsv1beta1.GroupName, Version: "v1beta1", Resource: resource}
	// v1GVR := metav1.GroupVersionResource{Group: apiextensionsv1.GroupName, Version: "v1", Resource: resource}

	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true

	raw := ar.Request.Object.Raw
	klog.V(2).Info(raw)

	rawJson := string(raw)
	logger.With(zap.String("raw", rawJson)).Info("ar.Request.Object.Raw")

	host := "http://todo-app-svc:5000"
	message := map[string]interface{}{
		"raw": rawJson,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		logger.Error(err.Error())
	}

	resp, err := http.Post(host+"/api/todo/validate", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		logger.Error(err.Error())
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	//log.Println(result)

	logger.With(zap.Any("raw", result)).Info("raw response")

	// var labels map[string]string

	// switch ar.Request.Resource {
	// case v1beta1GVR:
	// 	crd := apiextensionsv1beta1.CustomResourceDefinition{}
	// 	deserializer := codecs.UniversalDeserializer()
	// 	if _, _, err := deserializer.Decode(raw, nil, &crd); err != nil {
	// 		klog.Error(err)
	// 		return toV1AdmissionResponse(err)
	// 	}
	// 	labels = crd.Labels
	// case v1GVR:
	// 	crd := apiextensionsv1.CustomResourceDefinition{}
	// 	deserializer := codecs.UniversalDeserializer()
	// 	if _, _, err := deserializer.Decode(raw, nil, &crd); err != nil {
	// 		klog.Error(err)
	// 		return toV1AdmissionResponse(err)
	// 	}
	// 	labels = crd.Labels
	// default:
	// 	err := fmt.Errorf("expect resource to be one of [%v, %v] but got %v", v1beta1GVR, v1GVR, ar.Request.Resource)
	// 	klog.Error(err)
	// 	return toV1AdmissionResponse(err)
	// }

	// if v, ok := labels["webhook-e2e-test"]; ok {
	// 	if v == "webhook-disallow" {
	// 		reviewResponse.Allowed = false
	// 		reviewResponse.Result = &metav1.Status{Message: "the crd contains unwanted label"}
	// 	}
	// }
	return &reviewResponse

}