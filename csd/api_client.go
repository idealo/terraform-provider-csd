package csd

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"time"
)

// HostURL Set to production endpoint of API
const HostURL string = "https://api.common-short-domain.aws.idealo.cloud"

// ApiClient that holds authentication details and convenience functions that wrap HTTP communication
type ApiClient struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

func (c *ApiClient) curl(method string, path string, body io.Reader) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(method, fmt.Sprintf("%s%s", HostURL, path), body)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	authorizationHeaders := signer(request, c.AccessKeyId, c.SecretAccessKey, c.SessionToken)
	request.Header.Add("X-Amz-Security-Token", c.SessionToken)
	request.Header.Add("X-Amz-Date", authorizationHeaders.date)
	request.Header.Add("Authorization", authorizationHeaders.authorizationHeaders)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-amz-content-sha256", fmt.Sprintf("%x", authorizationHeaders.payloadHash))

	response, err := client.Do(request)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	defer response.Body.Close()

	// TODO: react on 409 Conflict
	if response.StatusCode == 403 {
		// Create proper error message if AWS credentials are not valid, probably because they expired
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return nil, diag.FromErr(err)
		}
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't authenticate to API, please check AWS credentials",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode == 409 {
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return nil, diag.FromErr(err)
		}
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't create zone",
			Detail:   responseBody["message"],
		})
	} else if !slices.Contains([]int{200, 201, 204}, response.StatusCode) {
		// Create error message for any other unexpected errors
		body, _ := io.ReadAll(response.Body)
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected error message from API",
			Detail:   fmt.Sprintf("HTTP %d: %s", response.StatusCode, body),
		})
	}

	var output interface{}
	if output != nil {
		if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
			return nil, diag.FromErr(err)
		}
	}

	return output, diags
}
