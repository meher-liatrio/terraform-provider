package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	devops_resource "github.com/liatrio/devops-bootcamp/examples/ch7/devops-resources"
)

// GetDev - Returns a single dev
func (c *Client) GetDev(devID string) (*devops_resource.Dev, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dev/id/%s", c.HostURL, devID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	dev := devops_resource.Dev{}
	err = json.Unmarshal(body, &dev)
	if err != nil {
		return nil, err
	}

	return &dev, nil
}

// GetDevs - Returns list of devs
func (c *Client) GetDevs() ([]devops_resource.Dev, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dev", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	devs := []devops_resource.Dev{}
	err = json.Unmarshal(body, &devs)
	if err != nil {
		return nil, err
	}

	return devs, nil
}

// CreateDev - Create a new Dev
func (c *Client) CreateDev(dev devops_resource.Dev) (*devops_resource.Dev, error) {
	// Marshal the single Dev into JSON
	rb, err := json.Marshal(dev)
	if err != nil {
		return nil, err
	}

	// Create a new POST request with the JSON body
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/dev", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	// Perform the HTTP request
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response into an Dev struct
	devObj := devops_resource.Dev{}
	err = json.Unmarshal(body, &devObj)
	if err != nil {
		return nil, err
	}

	return &devObj, nil
}

// UpdateDev - Update an existing dev
func (c *Client) UpdateDev(dev devops_resource.Dev) (*devops_resource.Dev, error) {
	log.Printf("\nUpdating dev: %+v\n", dev) // Add debug log

	// Marshal the single Dev into JSON
	rb, err := json.Marshal(dev)
	if err != nil {
		log.Printf("\nError marshalling dev: %s\n", err) // Add debug log
		return nil, err
	}

	// Create a new PUT request with the JSON body
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/dev/%s", c.HostURL, strings.Trim(dev.Id, "\"")), strings.NewReader(string(rb)))
	if err != nil {
		log.Printf("\nError creating request: %v\n", req) // Add debug log
		log.Printf("\nError creating request: %s\n", err) // Add debug log
		return nil, err
	}

	// Perform the HTTP request
	body, err := c.doRequest(req)
	log.Printf("\nResponse body: %s\n", body) // Add debug log
	if err != nil {
		log.Printf("\nError performing request: %s\n", err) // Add debug log
		return nil, err
	}

	// Unmarshal the response into an Dev struct
	// an_dev := devops_resource.Dev{}
	err = json.Unmarshal(body, &dev)
	if err != nil {
		log.Printf("\nError unmarshalling response: %s\n", err) // Add debug log
		return nil, err
	}

	return &dev, nil
}

// DeleteDev - Delete an existing dev
func (c *Client) DeleteDev(id string) error {
	log.Printf("\nDeleting dev: %+s\n", id) // Add debug log

	// Create a new Delete request with the JSON body
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/dev/%s", c.HostURL, strings.Trim(id, "\"")), nil)
	if err != nil {
		log.Printf("\nError creating request: %v\n", req) // Add debug log
		log.Printf("\nError creating request: %s\n", err) // Add debug log
		return err
	}

	// Perform the HTTP request
	body, err := c.doRequest(req)
	log.Printf("\nResponse body: %s\n", body) // Add debug log
	if err != nil {
		log.Printf("\nError performing request: %s\n", err) // Add debug log
		return err
	}

	// Unmarshal the response into an Dev struct
	// an_dev := devops_resource.Dev{}
	if err != nil {
		log.Printf("\nError unmarshalling response: %s\n", err) // Add debug log
		return err
	}

	return nil
}

type EngineerPayload struct {
	EngineerId string `json:"id"`
}

// AddEngToDev - adds engineer to dev engineers list
func (c *Client) AddEngToDev(DevId string, EngId string) error {
	// Create the payload
	payload := EngineerPayload{
		EngineerId: EngId,
	}
	// Marshal the engineer id into JSON
	rb, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create a new POST request with the JSON body
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/dev/%s", c.HostURL, DevId), bytes.NewBuffer(rb))
	if err != nil {
		return err
	}

	// Perform the HTTP request
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	// Unmarshal the response into an Dev struct
	devObj := devops_resource.Dev{}
	err = json.Unmarshal(body, &devObj)
	if err != nil {
		return err
	}

	return nil
}
