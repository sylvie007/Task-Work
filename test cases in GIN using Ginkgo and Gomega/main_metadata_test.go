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
		"url": "/tosca-models/csars/dcaf-cmts-argo-events.csar",
		"resolve": true,
		"coerce": false,
		"quirks": [
			"data_types.string.permissive"
		],
		"output": "dcaf.json",
		"inputs": "",
		"inputsUrl": "",
		"force": true
	}`
	_, err := ApiCall("POST", apiURL, apiBody)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	apiURL := "http://localhost:10010/compiler/v1/model/db/dcaf_service"
	apiBody := `{
		"namespace": "zip:file:d:/tosca-models/csars/dcaf-cmts-argo-events.csar!/dcaf_service.yaml",
		"version": "tick_profile_1_0",
		"includeTypes": true
	}`
	var err error
	apiResponse, err = ApiCall("DELETE", apiURL, apiBody)
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("API Test GET", func() {

	It("Re-test GET API after saving model", func() {
		apiURL := "http://localhost:10010/compiler/v1/db/models/metadata"
		apiResponse, err := ApiCall("GET", apiURL, ``)
		Expect(err).NotTo(HaveOccurred())

		// Asserting metadata count
		Expect(len(apiResponse.Data.Models[0].Metadata)).To(Equal(3))
	})
})

func ApiCall(apiType string, apiURL string, apiBody string) (APIResponse, error) {
	body := bytes.NewReader([]byte(apiBody))
	request, err := http.NewRequest(apiType, apiURL, body)
	if err != nil {
		log.Printf("Error while creating API request: %s", err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error during API call:", err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
	}
	var apiResponse APIResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	if err != nil {
		log.Println("Error unmarshaling JSON:", err)
	}
	return apiResponse, nil
}

type APIResponse struct {
	Result string `json:"result"`
	Data   struct {
		Models []struct {
			Metadata map[string]string `json:"metadata"`
		} `json:"models"`
	} `json:"data"`
}
