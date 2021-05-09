package msi

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const instanceMetaDataURL string = "http://169.254.169.254/metadata/instance?api-version=2017-04-02"

type MetaData struct {
	SubscriptionId string `json:"subscriptionId"`
	VMName string `json:"name"`
	VMssName string `json:"vmScaleSetName"`
	ResourceGroupName string `json:"resourceGroupName"`
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

/*GetInstanceMetadata ()
 *Calls the Azure in-VM Instance Metadata service and returns the results to the caller*/
func GetInstanceMetadata() (MetaData, error) {
	var metadata MetaData

	// Build a request to call the instance Azure in-VM metadata service
	req, err := http.NewRequest("GET", instanceMetaDataURL, nil)
	if err != nil {
		log.Error("Failed creating http request --- %s", err)
		return metadata, errors.New("failed creating http request object to retrieve instance metadata")
	}

	// Set the required header for the HTTP request
	req.Header.Add("Metadata", "true")

	// Create the HTTP client and call the instance metadata service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Failed calling instance metadata service --- %s", err)
		return metadata, errors.New("failed calling the in-VM instance metadata service")

	}
	// Complete reading the body
	defer resp.Body.Close()

	// Now return the instance metadata JSON or another error if the status code is not in 2xx range
	if (resp.StatusCode >= 200) && (resp.StatusCode <= 299) {
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&metadata)
		if err != nil {
			log.Error("Failed decoding Metadata from metadata endpoint --- %s", err)
			return metadata, errors.New("failed decoding MSI token from MSI token endpoint")
		}
		return metadata, nil
	}

	log.Error("Failed with Non-200 status code: %q", resp.StatusCode)
	return metadata, errors.New(fmt.Sprintf("instance meta data service returned non-OK status code: %q", resp.StatusCode))
}
