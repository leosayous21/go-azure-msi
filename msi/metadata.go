package msi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const instanceMetaDataURL string = "http://169.254.169.254/metadata/instance?api-version=2019-08-15"

type Compute struct {
	Compute MetaData `json:"compute"`
}

type MetaData struct {
	SubscriptionId string `json:"subscriptionId"`
	VMName string `json:"name"`
	VMssName string `json:"vmScaleSetName"`
	ResourceGroupName string `json:"resourceGroupName"`
}

/*GetInstanceMetadata ()
 *Calls the Azure in-VM Instance Metadata service and returns the results to the caller*/
func GetInstanceMetadata() (MetaData, error) {
	var metadata MetaData

	// Build a request to call the instance Azure in-VM metadata service
	req, err := http.NewRequest("GET", instanceMetaDataURL, nil)
	if err != nil {
		return metadata, errors.New(fmt.Sprintf("failed creating http request object to retrieve instance metadata - %s", err))
	}

	// Set the required header for the HTTP request
	req.Header.Add("Metadata", "true")

	// Create the HTTP client and call the instance metadata service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return metadata, errors.New(fmt.Sprintf("failed calling the in-VM instance metadata service  %s", err))
	}
	// Complete reading the body
	defer resp.Body.Close()

	// Now return the instance metadata JSON or another error if the status code is not in 2xx range
	if (resp.StatusCode >= 200) && (resp.StatusCode <= 299) {
		bodyContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return metadata, errors.New(fmt.Sprintf("failed reading resp.Body from metadata endpoint - %s", err))
		}
		var compute Compute
		json.Unmarshal(bodyContent, &compute)
		metadata = compute.Compute // Metadata is nested inside compute
		if err != nil {
			return metadata, errors.New(fmt.Sprintf("Failed decoding Metadata from metadata endpoint %s", err))
		}
		return metadata, nil
	}
	return metadata, errors.New(fmt.Sprintf("instance meta data service returned non-OK status code: %q", resp.StatusCode))
}
