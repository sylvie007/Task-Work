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

// Define the constant for the clout file name
const cloutFileName = "gin/compiler/dcaf_service.json"

var dcafmultilist Dcafmultilist

func loadDcafmultilist() {
	data, err := ioutil.ReadFile("dcafmultilist.json")
	if err != nil {
		log.Fatalf("Error reading dcafmultilist file: %v", err)
	}
	err = json.Unmarshal(data, &dcafmultilist)
	if err != nil {
		log.Fatalf("Error parsing dcafmultilist file: %v", err)
	}
}

func TestCompilerApiOperations(t *testing.T) {
	RegisterFailHandler(Fail)
	loadDcafmultilist()
	RunSpecs(t, "Compiler Operations Suite")
}

var _ = BeforeSuite(func() {
	apiURL := dcafmultilist.CreateInstanceAPI.CreateInstanceURL
	apiBody := dcafmultilist.CreateInstanceAPI.CreateInstanceBody
	_, err := ApiCall("POST", apiURL, apiBody)
	Expect(err).NotTo(HaveOccurred())

})

var _ = Describe("Service Orchestrator APIs", func() {
	var APIResponseInstances []InstanceData
	var _ = BeforeEach(func() {
		apiURL := dcafmultilist.GetInstancesAPI.GetInstancesURL
		responseBody, err := ApiCall("GET", apiURL, "")
		Expect(err).NotTo(HaveOccurred())
		err = json.Unmarshal(responseBody, &APIResponseInstances)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return the expected number of instances", func() {
		Expect(len(APIResponseInstances)).To(Equal(dcafmultilist.GetInstancesAPI.ExpectedUidCount))
	})
	var _ = BeforeEach(func() {
		apiURL := dcafmultilist.DemoInstanceAPI.APIURL
		responseBody, err := ApiCall("GET", apiURL, "")
		Expect(err).NotTo(HaveOccurred())
		err = json.Unmarshal(responseBody, &demoInstanceResponse)
		Expect(err).NotTo(HaveOccurred())
	})

	var _ = Describe("GET instances by name- API", func() {

		It("should have the expected number of vertexes", func() {
			Expect(len(demoInstanceResponse.Vertexes)).To(Equal(dcafmultilist.DemoInstanceAPI.ExpectedVtx))
		})

		It("should have the correct name", func() {
			Expect(demoInstanceResponse.Name).To(Equal("demo1"))
		})

		It("should have empty dependent_instance", func() {
			Expect(len(demoInstanceResponse.DependentInstance)).To(Equal(1)) // Assuming empty dependent_instance has a single empty string
			Expect(demoInstanceResponse.DependentInstance[0]).To(Equal(""))
		})

	})
	var _ = BeforeEach(func() {
		apiURL := dcafmultilist.DeployedInstancesAPI.APIURL
		responseBody, err := ApiCall("GET", apiURL, "")
		Expect(err).NotTo(HaveOccurred())
		err = json.Unmarshal(responseBody, &deployedInstancesResponse)
		Expect(err).NotTo(HaveOccurred())
	})

	// Add a new Describe block for the new API
	var _ = Describe("Deployed Instances APIs", func() {
		It("should return the correct data", func() {
			Expect(deployedInstancesResponse.Data).To(Equal(dcafmultilist.DeployedInstancesAPI.ExpectedData))
		})

		It("should have the correct message", func() {
			Expect(deployedInstancesResponse.Message).To(Equal(dcafmultilist.DeployedInstancesAPI.ExpectedMessage))
		})

		It("should have the correct result", func() {
			Expect(deployedInstancesResponse.Result).To(Equal(dcafmultilist.DeployedInstancesAPI.ExpectedResult))
		})
	})

	var _ = BeforeEach(func() {
		// Read the content of the JSON file containing the clout file data
		jsonData, err := ioutil.ReadFile(cloutFileName)
		Expect(err).NotTo(HaveOccurred())
	
		// Unmarshal the JSON data into the existing struct variable
		err = json.Unmarshal(jsonData, &dcafmultilist)
		Expect(err).NotTo(HaveOccurred())
	
		// Get the save clout API URL from the loaded struct
		saveCloutAPIURL := dcafmultilist.SaveCloutFileAPI.SavecloutURL
	
		// Make the API call to save the clout file
		responseBody, err := ApiCall("PUT", saveCloutAPIURL, "")
		Expect(err).NotTo(HaveOccurred())
	
		// Unmarshal the response body to check the expected message and result
		var saveCloutResponse struct {
			Message string `json:"message"`
			Result  string `json:"result"`
		}
		err = json.Unmarshal(responseBody, &saveCloutResponse)
		Expect(err).NotTo(HaveOccurred())
	
		// Check if the message and result match the expected values
		Expect(saveCloutResponse.Message).To(Equal(dcafmultilist.SaveCloutFileAPI.ExpectedMessage))
		Expect(saveCloutResponse.Result).To(Equal(dcafmultilist.SaveCloutFileAPI.ExpectedResult))
	})
	
	var _ = BeforeEach(func() {
		// Read the content of the JSON file containing the clout file data
		jsonData, err := ioutil.ReadFile(cloutFileName)
		Expect(err).NotTo(HaveOccurred())
	
		// Unmarshal the JSON data into the existing struct variable
		err = json.Unmarshal(jsonData, &dcafmultilist)
		Expect(err).NotTo(HaveOccurred())
	
		// Get the read clout API URL from the loaded struct
		readCloutAPIURL := dcafmultilist.ReadCloutAPI.ReadCloutURL
	
		// Make the API call to read the clout file
		responseBody, err := ApiCall("GET", readCloutAPIURL, "")
		Expect(err).NotTo(HaveOccurred())
	
		// Unmarshal the response body to check the expected message and result
		var readCloutResponse struct {
			Message string `json:"message"`
			Result  string `json:"result"`
			Data    []string `json:"data"`
		}
		err = json.Unmarshal(responseBody, &readCloutResponse)
		Expect(err).NotTo(HaveOccurred())
	
		// Check if the message and result match the expected values
		Expect(readCloutResponse.Message).To(Equal(dcafmultilist.ReadCloutAPI.ExpectedMessage))
		Expect(readCloutResponse.Result).To(Equal(dcafmultilist.ReadCloutAPI.ExpectedResult))
	})
	
	var _ = BeforeEach(func() {
		// Read the content of the JSON file containing the clout file data
		jsonData, err := ioutil.ReadFile(cloutFileName)
		Expect(err).NotTo(HaveOccurred())
	
		// Unmarshal the JSON data into the existing struct variable
		err = json.Unmarshal(jsonData, &dcafmultilist)
		Expect(err).NotTo(HaveOccurred())
	
		// Get the parse model API URL from the loaded struct
		parseModelAPIURL := dcafmultilist.ParseModelAPI.ParseModelURL
	
		// Make the API call to parse the clout file
		responseBody, err := ApiCall("POST", parseModelAPIURL, "")
		Expect(err).NotTo(HaveOccurred())
	
		// Unmarshal the response body to check the expected message and result
		var parseModelResponse struct {
			Message string `json:"message"`
			Result  string `json:"result"`
		}
		err = json.Unmarshal(responseBody, &parseModelResponse)
		Expect(err).NotTo(HaveOccurred())
	
		// Check if the message and result match the expected values
		Expect(parseModelResponse.Message).To(Equal(dcafmultilist.ParseModelAPI.ExpectedMessage))
		Expect(parseModelResponse.Result).To(Equal(dcafmultilist.ParseModelAPI.ExpectedResult))
	})
	

var _ = AfterSuite(func() {
	apiURL := dcafmultilist.DeleteInstanceAPI.DeleteInstanceURL
	var err error
	_, err = ApiCall("DELETE", apiURL, ``)
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

type CreateInstanceAPI struct {
	CreateInstanceURL  string `json:"createInstanceURL"`
	CreateInstanceBody string `json:"createInstanceBody"`
}
type InstanceData struct {
	UID               string            `json:"uid"`
	Name              string            `json:"name"`
	DependentInstance []string          `json:"dependent_instance"`
	Version           string            `json:"version"`
	GrammarVersion    string            `json:"grammarversion"`
	Properties        map[string]string `json:"properties"`
	Vertexes          []interface{}     `json:"vertexes"`
}

type DemoInstanceData struct {
	UID               string            `json:"uid"`
	Name              string            `json:"name"`
	DependentInstance []string          `json:"dependent_instance"`
	Version           string            `json:"version"`
	GrammarVersion    string            `json:"grammarversion"`
	Properties        map[string]string `json:"properties"`
	Vertexes          []interface{}     `json:"vertexes"`
}

var demoInstanceResponse DemoInstanceData

type DeployedInstancesResponse struct {
	Result  string   `json:"result"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

var deployedInstancesResponse DeployedInstancesResponse

type Dcafmultilist struct {
	CreateInstanceAPI struct {
		CreateInstanceURL  string `json:"createInstanceURL"`
		CreateInstanceBody string `json:"createInstanceBody"`
	} `json:"createInstanceAPI"`
	GetInstancesAPI struct {
		GetInstancesURL      string `json:"getInstancesURL"`
		ExpectedResult       string `json:"expectedResult"`
		ExpectedUidCount     int    `json:"expectedUidCount"`
		ExpectedVersionCount int    `json:"expectedVersionCount"`
	} `json:"getInstancesAPI"`
	DemoInstanceAPI struct {
		APIURL       string `json:"apiURL"`
		ExpectedAttr int    `json:"expectedAttributes"`
		ExpectedVtx  int    `json:"expectedVertexes"`
	} `json:"demoInstanceAPI"`
	DeleteInstanceAPI struct {
		DeleteInstanceURL string `json:"deleteModelURL"`
	} `json:"deleteInstanceAPI"`
	DeployedInstancesAPI struct {
		APIURL          string   `json:"apiURL"`
		ExpectedData    []string `json:"expectedData"`
		ExpectedMessage string   `json:"expectedMessage"`
		ExpectedResult  string   `json:"expectedResult"`
	} `json:"deployedInstancesAPI"`
	SaveCloutFileAPI struct {
		SavecloutURL    string `json:"savecloutURL"`
		ExpectedMessage string `json:"expectedMessage"`
		ExpectedResult  string `json:"expectedResult"`
	}
	ReadCloutAPI struct {
		ReadCloutURL    string   `json:"readcloutURL"`
		ExpectedData    []string `json:"expectedData"`
		ExpectedMessage string   `json:"expectedMessage"`
	}
	ParseModelAPI struct {
		ParseModelURL string `json:"parseModelURL"`
	}
}
