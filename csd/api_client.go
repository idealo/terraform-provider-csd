package csd

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"strings"
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

func (c *ApiClient) getZone(name string) (Zone, diag.Diagnostics) {
	var diags diag.Diagnostics
	var zone Zone

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/zones/%s", HostURL, name), strings.NewReader(""))
	if err != nil {
		return zone, diag.FromErr(err)
	}
	signRequest(request, c.AccessKeyId, c.SecretAccessKey, c.SessionToken)
	if err != nil {
		return zone, diag.FromErr(err)
	}
	response, err := client.Do(request)
	if err != nil {
		return zone, diag.FromErr(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		// Create proper error message if AWS credentials are not valid, probably because they expired
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return zone, diag.FromErr(err)
		}
		return zone, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't authenticate to API, please check AWS credentials",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode != 200 {
		// Create error message for any other unexpected errors
		body, _ := io.ReadAll(response.Body)
		return zone, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected error message from API",
			Detail:   fmt.Sprintf("HTTP %d: %s", response.StatusCode, body),
		})
	}

	if err := json.NewDecoder(response.Body).Decode(&zone); err != nil {
		return zone, diag.FromErr(err)
	}
	return zone, diags
}

func (c *ApiClient) curl(method string, path string, body io.Reader) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(method, fmt.Sprintf("%s%s", HostURL, path), body)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	signRequest(request, c.AccessKeyId, c.SecretAccessKey, c.SessionToken)

	response, err := client.Do(request)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	defer response.Body.Close()

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
