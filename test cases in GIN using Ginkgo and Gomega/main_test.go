package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCompilerApiOperations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Compiler Operations Suite")
}

var apiResponse APIResponse

var _ = BeforeSuite(func() {
	apiURL := "http://localhost:10010/compiler/v1/model/db/save"
	apiBody := `{
		"url": "/tosca-models/csars/cluster-resource.csar",
		"resolve": true,
		"coerce": false,
		"quirks": [
			"data_types.string.permissive"
		],
		"output": "cluster_input_service.json",
		"inputs": "",
		"inputsUrl": "",
		"force": true
	}`
	_, err := ApiCall("POST", apiURL, apiBody)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	apiURL := "http://localhost:10010/compiler/v1/model/db/cluster_input_service"
	apiBody := `{
		"namespace": "zip:file:c:/tosca-models/csars/cluster-resource.csar!/cluster_input_service.yaml",
		"version": "tick_profile_1_0",
		"includeTypes": true
	}`
	var err error
	apiResponse, err = ApiCall("DELETE", apiURL, apiBody)
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("API Test GET", func() {

	It("Re-test GET API after saving model", func() {
		apiURL := "http://localhost:10010/compiler/v1/db/models"
		apiResponse, err := ApiCall("GET", apiURL, ``)
		Expect(err).NotTo(HaveOccurred())

		serviceURLs := make([]interface{}, len(apiResponse.Data.ListOfModels))
		for i, model := range apiResponse.Data.ListOfModels {
			if modelMap, ok := model.(map[string]interface{}); ok {
				if serviceURL, exists := modelMap["service_url"]; exists {
					serviceURLs[i] = serviceURL
				}
			}
		}

		Expect(serviceURLs).To(ConsistOf(
			Equal("zip:file:c:/tosca-models/csars/cluster-resource.csar!/cluster_input_service.yaml"),
		))
	})
})

func ApiCall(apiType string, apiURL string, apiBody string) (APIResponse, error) {
	body := bytes.NewReader([]byte(apiBody))
	request, err := http.NewRequest(apiType, apiURL, body)
	if err != nil {
		log.Printf("Error while API call: %s", err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error client.Do(request):", err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error responseBody:", err)
	}
	var apiResponse APIResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	if err != nil {
		log.Println("Error unmarshaling JSON:", err)
	}
	return apiResponse, nil
}

type APIResponse struct {
	Result  string `json:"result"`
	Message string `json:"message"`
	Data    struct {
		ListOfModels []interface{} `json:"listOfModels"`
	} `json:"data"`
}
