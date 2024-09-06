locals {
  username = "application-account-hkraemer-api-test"
  password = "IJ4Idc0bNaYUuNMw382MoJpNtvroV0wB"
}
terraform {
  required_providers {
    hashicups = {
      source = "serviceware/aipe"
    }
  }
}

provider "aipe" {
  aipe_url = "https://saas-ops-1.ai-process-engine.testing.swops.cloud"  
  authenticator_realm_url = "https://auth.sso.testing.swops.cloud/auth/realms/saas-ops-1"

  application_username = local.username
  application_password = local.password
}

data "aipe_data_object" "test_object" {
  id = "436099"
}

output "test_object" {
  value = data.aipe_data_object.test_object
}
