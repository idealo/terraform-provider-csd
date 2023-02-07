package csd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
	UserAgent       string
}

type Zone struct {
	Name        string   `json:"name"`
	NameServers []string `json:"name_servers"`
}

func (c *ApiClient) createZone(zone Zone) (Zone, diag.Diagnostics) {
	var diags diag.Diagnostics

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(zone); err != nil {
		return zone, diag.FromErr(err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/zones", HostURL), buffer)
	if err != nil {
		return zone, diag.FromErr(err)
	}
	authorizationHeaders := signRequest(request, c.AccessKeyId, c.SecretAccessKey, c.SessionToken)
	request.Header.Add("X-Amz-Security-Token", c.SessionToken)
	request.Header.Add("X-Amz-Date", authorizationHeaders.date)
	request.Header.Add("Authorization", authorizationHeaders.authorizationHeaders)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-amz-content-sha256", fmt.Sprintf("%x", authorizationHeaders.payloadHash))
	request.Header.Set("User-Agent", c.UserAgent)

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
	} else if response.StatusCode == 409 {
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return zone, diag.FromErr(err)
		}
		return zone, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't create zone",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode != 201 {
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

func (c *ApiClient) getZone(name string) (Zone, diag.Diagnostics) {
	var diags diag.Diagnostics
	var zone Zone

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/zones/%s", HostURL, name), strings.NewReader(""))
	if err != nil {
		return zone, diag.FromErr(err)
	}
	authorizationHeaders := signRequest(request, c.AccessKeyId, c.SecretAccessKey, c.SessionToken)
	request.Header.Add("X-Amz-Security-Token", c.SessionToken)
	request.Header.Add("X-Amz-Date", authorizationHeaders.date)
	request.Header.Add("Authorization", authorizationHeaders.authorizationHeaders)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-amz-content-sha256", fmt.Sprintf("%x", authorizationHeaders.payloadHash))
	request.Header.Set("User-Agent", c.UserAgent)

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
	} else if response.StatusCode == 404 {
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return zone, diag.FromErr(err)
		}
		return zone, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't find zone with given name",
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

func (c *ApiClient) getZones() ([]Zone, diag.Diagnostics) {
	var diags diag.Diagnostics
	var zones []Zone

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/zones", HostURL), strings.NewReader(""))
	if err != nil {
		return zones, diag.FromErr(err)
	}
	authorizationHeaders := signRequest(request, c.AccessKeyId, c.SecretAccessKey, c.SessionToken)
	request.Header.Add("X-Amz-Security-Token", c.SessionToken)
	request.Header.Add("X-Amz-Date", authorizationHeaders.date)
	request.Header.Add("Authorization", authorizationHeaders.authorizationHeaders)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-amz-content-sha256", fmt.Sprintf("%x", authorizationHeaders.payloadHash))
	request.Header.Set("User-Agent", c.UserAgent)

	response, err := client.Do(request)
	if err != nil {
		return zones, diag.FromErr(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		// Create proper error message if AWS credentials are not valid, probably because they expired
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return zones, diag.FromErr(err)
		}
		return zones, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't authenticate to API, please check AWS credentials",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode != 200 {
		// Create error message for any other unexpected errors
		body, _ := io.ReadAll(response.Body)
		return zones, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected error message from API",
			Detail:   fmt.Sprintf("HTTP %d: %s", response.StatusCode, body),
		})
	}

	if err := json.NewDecoder(response.Body).Decode(&zones); err != nil {
		return zones, diag.FromErr(err)
	}
	return zones, diags
}

func (c *ApiClient) updateZone(zone Zone) (Zone, diag.Diagnostics) {
	var diags diag.Diagnostics

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(zone); err != nil {
		return zone, diag.FromErr(err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/zones/%s", HostURL, zone.Name), buffer)
	if err != nil {
		return zone, diag.FromErr(err)
	}
	authorizationHeaders := signRequest(request, c.AccessKeyId, c.SecretAccessKey, c.SessionToken)
	request.Header.Add("X-Amz-Security-Token", c.SessionToken)
	request.Header.Add("X-Amz-Date", authorizationHeaders.date)
	request.Header.Add("Authorization", authorizationHeaders.authorizationHeaders)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-amz-content-sha256", fmt.Sprintf("%x", authorizationHeaders.payloadHash))
	request.Header.Set("User-Agent", c.UserAgent)

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
	} else if response.StatusCode == 404 {
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return zone, diag.FromErr(err)
		}
		return zone, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't find zone with given name",
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

func (c *ApiClient) deleteZone(name string) diag.Diagnostics {
	var diags diag.Diagnostics

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/zones/%s", HostURL, name), strings.NewReader(""))
	if err != nil {
		return diag.FromErr(err)
	}
	authorizationHeaders := signRequest(request, c.AccessKeyId, c.SecretAccessKey, c.SessionToken)
	request.Header.Add("X-Amz-Security-Token", c.SessionToken)
	request.Header.Add("X-Amz-Date", authorizationHeaders.date)
	request.Header.Add("Authorization", authorizationHeaders.authorizationHeaders)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-amz-content-sha256", fmt.Sprintf("%x", authorizationHeaders.payloadHash))
	request.Header.Set("User-Agent", c.UserAgent)

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
	} else if response.StatusCode == 404 {
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return diag.FromErr(err)
		}
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't find zone with given name",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode != 204 {
		// Create error message for any other unexpected errors
		body, _ := io.ReadAll(response.Body)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected error message from API",
			Detail:   fmt.Sprintf("HTTP %d: %s", response.StatusCode, body),
		})
	}

	return diags
}
