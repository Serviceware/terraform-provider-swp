package authenticator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type AuthenticatorClient struct {
	Client *http.Client

	ApplicationUsername string
	ApplicationPassword string

	URL string
}

func (c AuthenticatorClient) Authenticate(ctx context.Context) (string, error) {
	tokenUrl := fmt.Sprintf("%s/protocol/openid-connect/token", c.URL)
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	req, err := http.NewRequestWithContext(ctx, "POST", tokenUrl, strings.NewReader(form.Encode()))
	tflog.Info(ctx, "creating request", map[string]interface{}{"url": tokenUrl, "err": err})

	if err != nil {
		return "", err
	}

	req.SetBasicAuth(c.ApplicationUsername, c.ApplicationPassword)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		tflog.Info(ctx, "login failed", map[string]interface{}{
			"status": resp.StatusCode,
			"body":   data,
		})
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}
