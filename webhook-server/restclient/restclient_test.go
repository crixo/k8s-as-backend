package restclient

import (
	"fmt"
	"testing"
)

func TestRestClient(t *testing.T) {
	host := "http://localhost/todo-app"
	code := "fake-code"
	rawJson := fmt.Sprintf(
		`{"apiVersion":"k8sasbackend.com/v1","kind":"Todo","metadata":{"creationTimestamp":"2020-02-23T13:24:07Z","generation":1,"name":"%s","namespace":"default","uid":"891e8fb1-9fb2-476c-a5bf-b6a65982fb23"},"spec":{"message":"string","when":"2020-02-23T13:23:53.873Z"}}`,
		code,
	)
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
	resp, err := Post(host+"/api/todo/validate", message, headers)
	if err != nil {
		t.Error(err.Error())
	}

	expected := 200
	actual := resp.StatusCode
	if expected != actual {
		t.Errorf("remote ep retruns an unexpected status code: got %v want %v", actual, expected)
	}

}
