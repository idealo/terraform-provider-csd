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

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(zone); err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/zones", HostURL), buf)
	if err != nil {
		return err
	}

	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return err
}
