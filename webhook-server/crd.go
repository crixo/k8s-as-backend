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
	"os"

	"github.com/crixo/k8s-as-backend/webhook-server/restclient"
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
	reviewResponse.Allowed = false

	raw := ar.Request.Object.Raw
	klog.V(2).Info(raw)

	rawJson := string(raw)
	logger.With(zap.String("raw", rawJson)).Info("ar.Request.Object.Raw")

	envKey := "TODO_APP_SVC"
	host := ""
	if value, exists := os.LookupEnv(envKey); exists {
		host = value
	} else {
		panic(fmt.Sprintf("The %s env variable is mandatory"))
	}
	//host := "http://todo-app-svc:5000"
	message := map[string]interface{}{
		"raw": rawJson,
	}

	// bytesRepresentation, err := json.Marshal(message)
	// if err != nil {
	// 	logger.Error(err.Error())
	// }

	headers := make(map[string][]string)
	headers["content-type"] = append(headers["content-type"], "application/json")
	// resp, err := http.Post(host+"/api/todo/validate", "application/json", bytes.NewBuffer(bytesRepresentation))
	resp, err := restclient.Post(host+"/api/todo/validate", message, headers)
	if err != nil {
		logger.Error(err.Error())
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	logger.With(zap.Any("raw", result)).Info("raw response")

	if v, ok := result["valid"].(bool); ok {
		reviewResponse.Allowed = v
	}

	return &reviewResponse

}
