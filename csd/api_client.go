package csd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const HostURL string = "https://common-short-domain.aws.idealo.cloud"

type ApiClient struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *ApiClient) CreateZone(zone Zone) error {
	client := &http.Client{Timeout: 10 * time.Second}

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(zone); err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/zones", HostURL), buffer)
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return err
}

func (c *ApiClient) ReadZone(zone Zone) error {

}

// ReadZones TODO: maybe remove zones completely and only keep zone
func (c *ApiClient) ReadZones(zone Zone) error {

}

func (c *ApiClient) UpdateZone(zone Zone) error {

}

func (c *ApiClient) DeleteZone(zone Zone) error {

}
