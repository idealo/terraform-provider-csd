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
const HostURL string = "https://csd.idealo.tools"

// ApiClient that holds authentication details and convenience functions that wrap HTTP communication
type ApiClient struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	UserAgent       string
}

// Zone Delegation

type ZoneDelegation struct {
	Name        string   `json:"name"`
	NameServers []string `json:"name_servers"`
}

func (c *ApiClient) createZoneDelegation(zoneDelegation ZoneDelegation) (ZoneDelegation, diag.Diagnostics) {
	var diags diag.Diagnostics

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(zoneDelegation); err != nil {
		return zoneDelegation, diag.FromErr(err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/zone_delegations", HostURL), buffer)
	if err != nil {
		return zoneDelegation, diag.FromErr(err)
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
		return zoneDelegation, diag.FromErr(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		// Create proper error message if AWS credentials are not valid, probably because they expired
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return zoneDelegation, diag.FromErr(err)
		}
		return zoneDelegation, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't authenticate to API, please check AWS credentials",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode == 409 {
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return zoneDelegation, diag.FromErr(err)
		}
		return zoneDelegation, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't create zone delegation",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode != 201 {
		// Create error message for any other unexpected errors
		body, _ := io.ReadAll(response.Body)
		return zoneDelegation, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected error message from API",
			Detail:   fmt.Sprintf("HTTP %d: %s", response.StatusCode, body),
		})
	}

	if err := json.NewDecoder(response.Body).Decode(&zoneDelegation); err != nil {
		return zoneDelegation, diag.FromErr(err)
	}
	return zoneDelegation, diags
}

func (c *ApiClient) getZoneDelegation(name string) (ZoneDelegation, diag.Diagnostics) {
	var diags diag.Diagnostics
	var zoneDelegation ZoneDelegation

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/zone_delegations/%s", HostURL, name), strings.NewReader(""))
	if err != nil {
		return zoneDelegation, diag.FromErr(err)
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
		return zoneDelegation, diag.FromErr(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		// Create proper error message if AWS credentials are not valid, probably because they expired
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return zoneDelegation, diag.FromErr(err)
		}
		return zoneDelegation, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't authenticate to API, please check AWS credentials",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode == 404 {
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return zoneDelegation, diag.FromErr(err)
		}
		return zoneDelegation, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't find zone delegation with given name",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode != 200 {
		// Create error message for any other unexpected errors
		body, _ := io.ReadAll(response.Body)
		return zoneDelegation, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected error message from API",
			Detail:   fmt.Sprintf("HTTP %d: %s", response.StatusCode, body),
		})
	}

	if err := json.NewDecoder(response.Body).Decode(&zoneDelegation); err != nil {
		return zoneDelegation, diag.FromErr(err)
	}
	return zoneDelegation, diags
}

func (c *ApiClient) getZoneDelegations() ([]ZoneDelegation, diag.Diagnostics) {
	var diags diag.Diagnostics
	var zoneDelegations []ZoneDelegation

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/zone_delegations", HostURL), strings.NewReader(""))
	if err != nil {
		return zoneDelegations, diag.FromErr(err)
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
		return zoneDelegations, diag.FromErr(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		// Create proper error message if AWS credentials are not valid, probably because they expired
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return zoneDelegations, diag.FromErr(err)
		}
		return zoneDelegations, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't authenticate to API, please check AWS credentials",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode != 200 {
		// Create error message for any other unexpected errors
		body, _ := io.ReadAll(response.Body)
		return zoneDelegations, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected error message from API",
			Detail:   fmt.Sprintf("HTTP %d: %s", response.StatusCode, body),
		})
	}

	if err := json.NewDecoder(response.Body).Decode(&zoneDelegations); err != nil {
		return zoneDelegations, diag.FromErr(err)
	}
	return zoneDelegations, diags
}

func (c *ApiClient) updateZoneDelegation(zoneDelegation ZoneDelegation) (ZoneDelegation, diag.Diagnostics) {
	var diags diag.Diagnostics

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(zoneDelegation); err != nil {
		return zoneDelegation, diag.FromErr(err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v2/zone_delegations/%s", HostURL, zoneDelegation.Name), buffer)
	if err != nil {
		return zoneDelegation, diag.FromErr(err)
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
		return zoneDelegation, diag.FromErr(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		// Create proper error message if AWS credentials are not valid, probably because they expired
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return zoneDelegation, diag.FromErr(err)
		}
		return zoneDelegation, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't authenticate to API, please check AWS credentials",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode == 404 {
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return zoneDelegation, diag.FromErr(err)
		}
		return zoneDelegation, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't find zone delegation with given name",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode != 200 {
		// Create error message for any other unexpected errors
		body, _ := io.ReadAll(response.Body)
		return zoneDelegation, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected error message from API",
			Detail:   fmt.Sprintf("HTTP %d: %s", response.StatusCode, body),
		})
	}

	if err := json.NewDecoder(response.Body).Decode(&zoneDelegation); err != nil {
		return zoneDelegation, diag.FromErr(err)
	}
	return zoneDelegation, diags
}

func (c *ApiClient) deleteZoneDelegation(name string) diag.Diagnostics {
	var diags diag.Diagnostics

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v2/zone_delegations/%s", HostURL, name), strings.NewReader(""))
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
			Summary:  "Couldn't find zone delegation with given name",
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

// Record

type Record struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	TTL    int    `json:"ttl"`
	RRType string `json:"rrtype"`
}

func (c *ApiClient) createRecord(record Record) (Record, diag.Diagnostics) {
	var diags diag.Diagnostics

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(record); err != nil {
		return record, diag.FromErr(err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/records", HostURL), buffer)
	if err != nil {
		return record, diag.FromErr(err)
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
		return record, diag.FromErr(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		// Create proper error message if AWS credentials are not valid, probably because they expired
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return record, diag.FromErr(err)
		}
		return record, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't authenticate to API, please check AWS credentials",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode == 409 {
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return record, diag.FromErr(err)
		}
		return record, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't create record",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode != 201 {
		// Create error message for any other unexpected errors
		body, _ := io.ReadAll(response.Body)
		return record, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected error message from API",
			Detail:   fmt.Sprintf("HTTP %d: %s", response.StatusCode, body),
		})
	}

	if err := json.NewDecoder(response.Body).Decode(&record); err != nil {
		return record, diag.FromErr(err)
	}
	return record, diags
}

func (c *ApiClient) getRecord(name string) (Record, diag.Diagnostics) {
	var diags diag.Diagnostics
	var record Record

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/records/%s", HostURL, name), strings.NewReader(""))
	if err != nil {
		return record, diag.FromErr(err)
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
		return record, diag.FromErr(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		// Create proper error message if AWS credentials are not valid, probably because they expired
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return record, diag.FromErr(err)
		}
		return record, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't authenticate to API, please check AWS credentials",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode == 404 {
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return record, diag.FromErr(err)
		}
		return record, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't find record with given name",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode != 200 {
		// Create error message for any other unexpected errors
		body, _ := io.ReadAll(response.Body)
		return record, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected error message from API",
			Detail:   fmt.Sprintf("HTTP %d: %s", response.StatusCode, body),
		})
	}

	if err := json.NewDecoder(response.Body).Decode(&record); err != nil {
		return record, diag.FromErr(err)
	}
	return record, diags
}

func (c *ApiClient) getRecords() ([]Record, diag.Diagnostics) {
	var diags diag.Diagnostics
	var record []Record

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/records", HostURL), strings.NewReader(""))
	if err != nil {
		return record, diag.FromErr(err)
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
		return record, diag.FromErr(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		// Create proper error message if AWS credentials are not valid, probably because they expired
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return record, diag.FromErr(err)
		}
		return record, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't authenticate to API, please check AWS credentials",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode != 200 {
		// Create error message for any other unexpected errors
		body, _ := io.ReadAll(response.Body)
		return record, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected error message from API",
			Detail:   fmt.Sprintf("HTTP %d: %s", response.StatusCode, body),
		})
	}

	if err := json.NewDecoder(response.Body).Decode(&record); err != nil {
		return record, diag.FromErr(err)
	}
	return record, diags
}

func (c *ApiClient) updateRecord(record Record) (Record, diag.Diagnostics) {
	var diags diag.Diagnostics

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(record); err != nil {
		return record, diag.FromErr(err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v2/records/%s", HostURL, record.Name), buffer)
	if err != nil {
		return record, diag.FromErr(err)
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
		return record, diag.FromErr(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		// Create proper error message if AWS credentials are not valid, probably because they expired
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return record, diag.FromErr(err)
		}
		return record, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't authenticate to API, please check AWS credentials",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode == 404 {
		var responseBody map[string]string
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return record, diag.FromErr(err)
		}
		return record, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Couldn't find record with given id",
			Detail:   responseBody["message"],
		})
	} else if response.StatusCode != 200 {
		// Create error message for any other unexpected errors
		body, _ := io.ReadAll(response.Body)
		return record, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected error message from API",
			Detail:   fmt.Sprintf("HTTP %d: %s", response.StatusCode, body),
		})
	}

	if err := json.NewDecoder(response.Body).Decode(&record); err != nil {
		return record, diag.FromErr(err)
	}
	return record, diags
}

func (c *ApiClient) deleteRecord(name string) diag.Diagnostics {
	var diags diag.Diagnostics

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v2/records/%s", HostURL, name), strings.NewReader(""))
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
			Summary:  "Couldn't find record with given name",
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
