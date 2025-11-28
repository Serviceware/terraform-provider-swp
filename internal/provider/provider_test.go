package provider

import (
	"net/http"
	"os"
	"testing"

	"github.com/Serviceware/terraform-provider-swp/internal/aipe"
	"github.com/Serviceware/terraform-provider-swp/internal/authenticator"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"swp": providerserver.NewProtocol6WithError(New("test")()),
}

var aipeClient aipe.AIPEClient

var requiredEnvironmentVariables = []string{
	"SWP_APPLICATION_USER_USERNAME",
	"SWP_APPLICATION_USER_PASSWORD",
	"SWP_AUTHENTICATOR_URL",
	"SWP_AIPE_URL",
}

func testAccPreCheck(t *testing.T) {
	for _, envVar := range requiredEnvironmentVariables {
		if os.Getenv(envVar) == "" {
			t.Fatalf("Environment variable %s must be set for acceptance tests", envVar)
		}
	}

	client := http.DefaultClient
	authenticatorClient := authenticator.AuthenticatorClient{
		Client:              client,
		ApplicationUsername: os.Getenv("SWP_APPLICATION_USER_USERNAME"),
		ApplicationPassword: os.Getenv("SWP_APPLICATION_USER_PASSWORD"),
		URL:                 os.Getenv("SWP_AUTHENTICATOR_URL"),
	}

	aipeClient = aipe.AIPEClient{
		HTTPClient:    client,
		URL:           os.Getenv("SWP_AIPE_URL"),
		Authenticator: &authenticatorClient,
	}
}
