package csd

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"io"
	"net/http"
	"time"
)

// HostURL Set to production endpoint of API
const HostURL string = "https://6zrrgc0ria.execute-api.eu-central-1.amazonaws.com"

// ApiClient that holds authentication details and convenience functions that wrap HTTP communication
type ApiClient struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

func (c *ApiClient) curl(method string, path string, body io.Reader, output interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(method, fmt.Sprintf("%s%s", HostURL, path), body)
	if err != nil {
		return diag.FromErr(err)
	}
	authorizationHeaders := signer(request, c.AccessKeyId, c.SecretAccessKey, c.SessionToken)
	request.Header.Add("X-Amz-Security-Token", c.SessionToken)
	request.Header.Add("X-Amz-Date", authorizationHeaders.date)
	request.Header.Add("Authorization", authorizationHeaders.authorizationHeaders)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-amz-content-sha256", fmt.Sprintf("%x", authorizationHeaders.payloadHash))

	response, err := client.Do(request)
	if err != nil {
		return diag.FromErr(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		// Create proper error message if AWS credentials are not valid, probably because they expired
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return diag.FromErr(err)
		}
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't authenticate to API, please check AWS credentials",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode != 200 {
		// Create error message for any other unexpected errors
		body, _ := io.ReadAll(response.Body)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected error message from API",
			Detail:   fmt.Sprintf("HTTP %d: %s", response.StatusCode, body),
		})
	}

	if output != nil {
		if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
