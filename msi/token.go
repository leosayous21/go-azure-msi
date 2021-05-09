package msi

import (
	"errors"
	"net/url"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"fmt"
	"net/http"
	"encoding/json"
)

const msiTokenURL string = "http://169.254.169.254/metadata/identity/oauth2/token"
const resourceURL string = "https://management.azure.com/"

/*MsiToken ()
 *Encapsulates a token retrieved by the MSI extension on the local machine*/
type MsiToken struct {
	AccessToken 	string `json:"access_token"`
	RefreshToken 	string `json:"refresh_token"`
	ExpiresIn 		string `json:"expires_in"`
	ExpiresOn 		string `json:"expires_on"`
	NotBefore 		string `json:"not_before"`
	Resource 		string `json:"resource"`
	TokenType 		string `json:"token_type"`
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

/*GetMsiToken ()
 *Uses the Managed Service Identity Extension to retrieve a token that allows the VM to call into
 *the Azure Resource Manager APIs*/
func GetMsiToken() (token MsiToken, err error) {
	var myToken MsiToken

	// Build a request to call the MSI Extension OAuth2 Service
	// The request must contain the resource for which we request the token
	finalRequestURL := fmt.Sprintf("%s?api-version=2018-02-01&resource=%s", fmt.Sprintf(msiTokenURL), url.QueryEscape(resourceURL))
	req, err := http.NewRequest("GET", finalRequestURL, nil)
	if err != nil {
		return myToken, errors.New(fmt.Sprintf("failed creating http request object to request MSI token %s", err))
	}
	// Set the required header for the HTTP request
	req.Header.Add("Metadata", "true")

	// Create the HTTP client and call the instance metadata service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return myToken, errors.New(fmt.Sprintf("failed calling MSI token service %s", err))
	}
	// Complete reading the body
	defer resp.Body.Close()

	// Now return the instance metadata JSON or another error if the status code is not in 2xx range
	if (resp.StatusCode >= 200) && (resp.StatusCode <= 299) {
		dec := json.NewDecoder(resp.Body)
		err := dec.Decode(&myToken)
		if err != nil {
			return myToken, errors.New(fmt.Sprintf("failed decoding MSI token from MSI token endpoint  %s", err))
		}
		return myToken, nil
	}

	// Try to read the body and log the error details, nevertheless
	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(fmt.Sprintf("Failed reading response body from http response for more error details %s", err))
	}

	return myToken, errors.New(fmt.Sprintf("instance meta data service returned non-OK status code: %d - %s", resp.StatusCode, bodyContent))
}