package aipe

import (
	"context"
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
	DataObject map[string]string `json:"dataObject"`
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
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tflog.Info(ctx, "Successfully retrieved object", map[string]interface{}{"bytes": bodyBytes})

	var object ObjectAPIResponse
	json.Unmarshal(bodyBytes, &object)
	return object.DataObject, nil
}
