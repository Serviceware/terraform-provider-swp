package aipe

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type AIPEClient struct {
	// HTTPClient is the client used to make HTTP requests.
	HTTPClient *http.Client

	// URL is the base URL for the AIPE API.
	URL string

	// OIDCToken is the OIDC token used to authenticate with the AIPE API.
	OIDCToken string
}

type ObjectAPIResponse struct {
	DataObject map[string]interface{} `json:"dataObject"`
}

func (c *AIPEClient) GetObject(ctx context.Context, id string) (map[string]string, error) {
	// Make a request to the AIPE API to get the object with the specified ID.

	objectURL := fmt.Sprintf("%s/data/api/v1/objects/%s", c.URL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", objectURL, nil)
	tflog.Info(ctx, "creating request", map[string]interface{}{"url": objectURL, "err": err})
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.OIDCToken))

	tflog.Info(ctx, "reading object")
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tflog.Info(ctx, "get object failed", map[string]interface{}{"status": resp.StatusCode, "objectURL": objectURL})
		return nil, &ApiError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode)}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tflog.Info(ctx, "Successfully retrieved object", map[string]interface{}{"bytes": string(bodyBytes)})

	var object ObjectAPIResponse
	json.Unmarshal(bodyBytes, &object)
	delete(object.DataObject, "system")
	return convertPropertiesToString(object.DataObject), nil
}

type ObjectCreateRequest struct {
	Type       string                 `json:"typeName"`
	DataObject map[string]interface{} `json:"dataObject"`
}

type ObjectCreateResponse struct {
	ID string `json:"dataObjectId"`
}

func (c *AIPEClient) CreateObject(ctx context.Context, objectType string, data map[string]string) (string, error) {
	// Make a request to the AIPE API to create an object with the specified data.

	objectURL := fmt.Sprintf("%s/data/api/v1/objects", c.URL)

	requestObject := ObjectCreateRequest{
		Type:       objectType,
		DataObject: convertPropertiesFromString(data),
	}
	body, err := json.Marshal(requestObject)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", objectURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.OIDCToken))
	req.Header.Set("Content-Type", "application/json")

	tflog.Info(ctx, "creating object", map[string]interface{}{"objectType": objectType})
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respData, _ := io.ReadAll(resp.Body)
		blub := base64.StdEncoding.EncodeToString(respData)
		tflog.Info(ctx, "create object failed", map[string]interface{}{"status": resp.StatusCode, "objectURL": objectURL, "respData": blub})
		return "", &ApiError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode)}
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	tflog.Info(ctx, "Successfully created object", map[string]interface{}{"objectType": objectType, "data": data, "response": string(respData)})

	var createResponse ObjectCreateResponse
	err = json.Unmarshal(respData, &createResponse)
	if err != nil {
		return "", err
	}

	return createResponse.ID, nil
}

type ObjectUpdateRequest struct {
	DataObject map[string]interface{} `json:"dataObject"`
}

func (c *AIPEClient) UpdateObject(ctx context.Context, id string, data map[string]string) error {
	// Make a request to the AIPE API to update the object with the specified ID.
	objectURL := fmt.Sprintf("%s/data/api/v1/objects/%s", c.URL, id)

	requestObject := ObjectUpdateRequest{
		DataObject: convertPropertiesFromString(data),
	}
	body, err := json.Marshal(requestObject)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", objectURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.OIDCToken))
	req.Header.Set("Content-Type", "application/json")

	tflog.Info(ctx, "updating object", map[string]interface{}{"id": id})
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		tflog.Info(ctx, "update object failed", map[string]interface{}{"status": resp.StatusCode, "objectURL": objectURL, "respData": string(respData)})
		return &ApiError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode)}
	}

	tflog.Info(ctx, "Successfully updated object", map[string]interface{}{"id": id, "data": data})

	return nil
}

func (c *AIPEClient) DeleteObject(ctx context.Context, id string) error {
	// Make a request to the AIPE API to delete the object with the specified ID.
	objectURL := fmt.Sprintf("%s/data/api/v1/objects/%s", c.URL, id)

	req, err := http.NewRequestWithContext(ctx, "DELETE", objectURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.OIDCToken))

	tflog.Info(ctx, "deleting object", map[string]interface{}{"id": id})
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		respData, _ := io.ReadAll(resp.Body)
		tflog.Info(ctx, "delete object failed", map[string]interface{}{"status": resp.StatusCode, "objectURL": objectURL, "respData": string(respData)})
		return &ApiError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode)}
	}

	tflog.Info(ctx, "Successfully deleted object", map[string]interface{}{"id": id})

	return nil
}
