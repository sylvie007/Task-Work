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

var dcaf_resource Config

func loadConfig() {
	data, err := ioutil.ReadFile("dcaf_resource.json")
	if err != nil {
		log.Fatalf("Error reading dcaf_resource file: %v", err)
	}
	err = json.Unmarshal(data, &dcaf_resource)
	if err != nil {
		log.Fatalf("Error parsing dcaf_resourceig file: %v", err)
	}
}
func TestCompilerApiOperations(t *testing.T) {
	RegisterFailHandler(Fail)
	loadConfig()
	RunSpecs(t, "Compiler Operations Suite")
}

var _ = BeforeSuite(func() {
	apiURL := dcaf_resource.SaveModelAPI.SaveModelURL
	apiBody := dcaf_resource.SaveModelAPI.SaveModelBody
	_, err := ApiCall("POST", apiURL, apiBody)
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("Compiler APIs", func() {
	var APIResponseInputs APIResponseInputs
	var _ = BeforeEach(func() {
		apiURL := dcaf_resource.InputAPI.GetInputsURL
		apiBody := dcaf_resource.InputAPI.GetInputsBody
		responseBody, err := ApiCall("GET", apiURL, apiBody)
		Expect(err).NotTo(HaveOccurred())
		err = json.Unmarshal(responseBody, &APIResponseInputs)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should have expected count of dataTypeName Integer", func() {
		expectedCount := dcaf_resource.InputAPI.IntegerCounts
		totalCount := 0
		for _, events := range APIResponseInputs.Data {
			for _, event := range events {
				if event.DataTypeName == "integer" {
					totalCount++
				}
			}
		}
		Expect(totalCount).To(Equal(expectedCount))
	})

	It("should have expected count of dataTypeName String", func() {
		expectedCount := dcaf_resource.InputAPI.StringCounts
		totalCount := 0
		for _, events := range APIResponseInputs.Data {
			for _, event := range events {
				if event.DataTypeName == "string" {
					totalCount++
				}
			}
		}
		Expect(totalCount).To(Equal(expectedCount))
	})

	It("should have expected count of dataTypeName List", func() {
		expectedCount := dcaf_resource.InputAPI.ListCounts
		totalCount := 0
		for _, events := range APIResponseInputs.Data {
			for _, event := range events {
				if event.DataTypeName == "list" {
					totalCount++
				}
			}
		}
		Expect(totalCount).To(Equal(expectedCount))
	})

	It("should match expected name for each data object", func() {
		for _, events := range APIResponseInputs.Data {
			for _, event := range events {
				_, exists := dcaf_resource.InputAPI.ExpectedNames[event.Name]
				Expect(exists).To(BeTrue(), "Unexpected name found: %s", event.Name)
			}
		}
	})
})
var _ = AfterSuite(func() {
	apiURL := dcaf_resource.DeleteModelAPI.DeleteModelURL
	apiBody := dcaf_resource.DeleteModelAPI.DeleteModelBody
	var err error
	_, err = ApiCall("DELETE", apiURL, apiBody)
	Expect(err).NotTo(HaveOccurred())
})

func ApiCall(apiType string, apiURL string, apiBody string) ([]byte, error) {
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
	return responseBody, nil
}

type APIResponseInputs struct {
	Result  string                 `json:"result"`
	Message string                 `json:"message"`
	Data    map[string][]EventData `json:"data"`
}

type EventData struct {
	DataTypeName string      `json:"datatypename"`
	Default      interface{} `json:"default"`
	Name         string      `json:"name"`
	Namespace    struct {
		URL string `json:"url"`
	} `json:"namespace"`
}

type Config struct {
	SaveModelAPI struct {
		SaveModelURL  string `json:"saveModelURL"`
		SaveModelBody string `json:"saveModelBody"`
	} `json:"saveModelAPI"`
	DeleteModelAPI struct {
		DeleteModelURL  string `json:"deleteModelURL"`
		DeleteModelBody string `json:"deleteModelBody"`
	} `json:"deleteModelAPI"`
	InputAPI struct {
		GetInputsURL  string          `json:"getInputsURL"`
		GetInputsBody string          `json:"getInputsBody"`
		IntegerCounts int             `json:"integerCounts"`
		StringCounts  int             `json:"stringCounts"`
		ListCounts    int             `json:"listCounts"`
		ExpectedNames map[string]bool `json:"expectedNames"`
	} `json:"getInputAPI"`
}
