package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"swp": providerserver.NewProtocol6WithError(New("test")()),
}

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
}
