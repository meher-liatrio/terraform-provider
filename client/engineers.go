package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	devops_resource "github.com/liatrio/devops-bootcamp/examples/ch7/devops-resources"
)

func (c *Client) GetEngineer(engineerID string) (*devops_resource.Engineer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers/id/%s", c.HostURL, engineerID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	engineer := devops_resource.Engineer{}
	err = json.Unmarshal(body, &engineer)
	if err != nil {
		return nil, err
	}

	return &engineer, nil
}

// GetEngineers - Returns list of engineers (no auth required)
func (c *Client) GetEngineers() ([]devops_resource.Engineer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	engineers := []devops_resource.Engineer{}
	err = json.Unmarshal(body, &engineers)
	if err != nil {
		return nil, err
	}

	return engineers, nil
}

// CreateEngineer - Create a new order with a single order item
func (c *Client) CreateEngineer(engineer devops_resource.Engineer) (*devops_resource.Engineer, error) {
	// Marshal the single Engineer into JSON
	rb, err := json.Marshal(engineer)
	if err != nil {
		return nil, err
	}

	// Create a new POST request with the JSON body
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/engineers", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	// Perform the HTTP request
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response into an Engineer struct
	order := devops_resource.Engineer{}
	err = json.Unmarshal(body, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// UpdateEngineer - Update an existing engineer
func (c *Client) UpdateEngineer(engineer devops_resource.Engineer) (*devops_resource.Engineer, error) {
	log.Printf("\nUpdating engineer: %+v\n", engineer) // Add debug log

	// Marshal the single Engineer into JSON
	rb, err := json.Marshal(engineer)
	if err != nil {
		log.Printf("\nError marshalling engineer: %s\n", err) // Add debug log
		return nil, err
	}

	// Create a new PUT request with the JSON body
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/engineers/%s", c.HostURL, strings.Trim(engineer.Id, "\"")), strings.NewReader(string(rb)))
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

	// Unmarshal the response into an Engineer struct
	// an_engineer := devops_resource.Engineer{}
	err = json.Unmarshal(body, &engineer)
	if err != nil {
		log.Printf("\nError unmarshalling response: %s\n", err) // Add debug log
		return nil, err
	}

	return &engineer, nil
}
