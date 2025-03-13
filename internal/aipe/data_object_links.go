package aipe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type GetDataObjectLinksResponse struct {
	TotalElements int `json:"totalElements"`
	Objects       []struct {
		System struct {
			ID string `json:"id"`
		}
	} `json:"objects"`
}

func (c *AIPEClient) GetDataObjectLinks(ctx context.Context, id string, linkName string, relationName string) ([]string, error) {
	var objectIDs []string = nil
	totalElements := 1
	page := 0

	for len(objectIDs) < totalElements {

		// Make a request to the AIPE API to get the object with the specified ID.
		objectURL := fmt.Sprintf("%s/data/api/v1/objects/%s/links?linkDefinitionName=%s&relationName=%s&page=%d", c.URL, id, linkName, relationName, page)
		page++

		req, err := http.NewRequestWithContext(ctx, "GET", objectURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.OIDCToken))

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, &ApiError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode)}
		}

		tflog.Info(ctx, "Reading object links", map[string]interface{}{"url": objectURL})
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var objectLinks GetDataObjectLinksResponse
		json.Unmarshal(bodyBytes, &objectLinks)

		if len(objectIDs) == 0 { // first iteration
			totalElements = objectLinks.TotalElements
		}

		for _, object := range objectLinks.Objects {
			objectIDs = append(objectIDs, object.System.ID)
		}
	}

	slices.Sort(objectIDs)
	return objectIDs, nil
}

type LinkDefinition struct {
	LinkName     string   `json:"linkDefinitionName"`
	RelationName string   `json:"relationName"`
	Add          []string `json:"add,omitempty"`
	Remove       []string `json:"remove,omitempty"`
}

type UpdateDataObjectLinksRequest struct {
	Links []LinkDefinition `json:"links"`
}

func (c *AIPEClient) UpdateDataObjectLinks(ctx context.Context, id string, linkName string, relationName string, add []string, remove []string) error {
	tflog.Info(ctx, "Updating data object links", map[string]interface{}{"url": c.URL, "id": id, "linkName": linkName, "relationName": relationName, "add": add, "remove": remove})
	objectURL := fmt.Sprintf("%s/data/api/v1/objects/%s", c.URL, id)

	payload := UpdateDataObjectLinksRequest{
		Links: []LinkDefinition{
			{
				LinkName:     linkName,
				RelationName: relationName,
				Add:          add,
				Remove:       remove,
			},
		},
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", objectURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.OIDCToken))
	req.Header.Set("Content-Type", "application/json")

	tflog.Info(ctx, "Sending request to AIPE API", map[string]interface{}{"url": objectURL, "payload": string(payloadJSON)})
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		responseBody, _ := io.ReadAll(resp.Body)
		tflog.Info(ctx, "update object links failed", map[string]interface{}{"status": resp.StatusCode, "objectURL": objectURL, "response": string(responseBody)})
		return &ApiError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode)}
	}

	return nil
}
