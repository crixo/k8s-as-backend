//https://www.thegreatcodeadventure.com/mocking-http-requests-in-golang/
package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"net/http"

	v1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crixo/k8s-as-backend/webhook-server/mocks"
	"github.com/crixo/k8s-as-backend/webhook-server/restclient"
)

func TestHealthCheckHandler(t *testing.T) {

	restclient.Client = &mocks.MockClient{}

	// build response JSON
	jsonObj := `{"valid":true,"message":""}}`

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		// create a new reader with that JSON
		r := ioutil.NopCloser(bytes.NewReader([]byte(jsonObj)))
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	requestDto := &v1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Request: &v1.AdmissionRequest{
			UID: "e911857d-c318-11e8-bbad-025000000001",
			Kind: metav1.GroupVersionKind{
				Kind: "Namespace",
			},
			Operation: "CREATE",
			Object: runtime.RawExtension{
				Raw: []byte(`{"metadata": {
													"name": "test",
													"uid": "e911857d-c318-11e8-bbad-025000000001",
													"creationTimestamp": "2018-09-28T12:20:39Z"
												}}`),
			},
		},
	}
	bytesRepresentation, _ := json.Marshal(requestDto)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/crd", bytes.NewBuffer(bytesRepresentation))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serveCRD)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// // Check the response body is what we expect.
	// expected := `{"alive": true}`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
	t.Log("here")
	//t.Log(rr.Body.String())

	responseAdmissionReview := &v1.AdmissionReview{}
	json.NewDecoder(rr.Body).Decode(&responseAdmissionReview)
	t.Log(responseAdmissionReview.String())

}
