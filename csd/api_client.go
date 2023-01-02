package csd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const HostURL string = "https://6zrrgc0ria.execute-api.eu-central-1.amazonaws.com"

type AuthInfo struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

type ApiClient struct {
	AuthInfo AuthInfo
}

func (c *ApiClient) CreateZone(zone Zone) error {
	client := &http.Client{Timeout: 10 * time.Second}

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(zone); err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/zones", HostURL), buffer)
	if err != nil {
		return err
	}

	authorizationHeaders := signer(&c.AuthInfo, request)
	request.Header.Add("X-Amz-Security-Token", c.AuthInfo.SessionToken)
	request.Header.Add("X-Amz-Date", authorizationHeaders.date)
	request.Header.Add("Authorization", authorizationHeaders.authorizationHeaders)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-amz-content-sha256", fmt.Sprintf("%x", authorizationHeaders.payloadHash))

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return err
}

func (c *ApiClient) ReadZone(name string) (Zone, error) {
	zone := Zone{}

	client := &http.Client{Timeout: 10 * time.Second}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/zones/%s", HostURL, name), nil)
	if err != nil {
		return zone, err
	}

	authorizationHeaders := signer(&c.AuthInfo, request)
	request.Header.Add("X-Amz-Security-Token", c.AuthInfo.SessionToken)
	request.Header.Add("X-Amz-Date", authorizationHeaders.date)
	request.Header.Add("Authorization", authorizationHeaders.authorizationHeaders)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-amz-content-sha256", fmt.Sprintf("%x", authorizationHeaders.payloadHash))

	response, err := client.Do(request)
	if err != nil {
		return zone, err
	}
	defer response.Body.Close()

	// decode the response
	err = json.NewDecoder(response.Body).Decode(&zone)
	if err != nil {
		return zone, err
	}

	return zone, nil
}

func (c *ApiClient) ReadZones() ([]Zone, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	var zones []Zone

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/zones", HostURL), nil)
	if err != nil {
		return zones, err
	}

	response, err := client.Do(request)
	if err != nil {
		return zones, err
	}
	defer response.Body.Close()

	// decode the response

	err = json.NewDecoder(response.Body).Decode(&zones)
	if err != nil {
		return zones, err
	}

	return zones, nil
}

func (c *ApiClient) UpdateZone(zone Zone) error {
	client := &http.Client{Timeout: 10 * time.Second}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(zone); err != nil {
		return err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/zones/%s", HostURL, zone.Name), buf)
	if err != nil {
		return err
	}

	authorizationHeaders := signer(&c.AuthInfo, request)
	request.Header.Add("X-Amz-Security-Token", c.AuthInfo.SessionToken)
	request.Header.Add("X-Amz-Date", authorizationHeaders.date)
	request.Header.Add("Authorization", authorizationHeaders.authorizationHeaders)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-amz-content-sha256", fmt.Sprintf("%x", authorizationHeaders.payloadHash))

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

func (c *ApiClient) DeleteZone(name string) error {
	client := &http.Client{Timeout: 10 * time.Second}

	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/zones/%s", HostURL, name), nil)
	if err != nil {
		return err
	}

	authorizationHeaders := signer(&c.AuthInfo, request)
	request.Header.Add("X-Amz-Security-Token", c.AuthInfo.SessionToken)
	request.Header.Add("X-Amz-Date", authorizationHeaders.date)
	request.Header.Add("Authorization", authorizationHeaders.authorizationHeaders)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-amz-content-sha256", fmt.Sprintf("%x", authorizationHeaders.payloadHash))

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}
