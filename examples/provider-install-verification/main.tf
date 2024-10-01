locals {
  username = "application-account-hkraemer-api-test"
  password = "IJ4Idc0bNaYUuNMw382MoJpNtvroV0wB"
}
/*
terraform {
  required_providers {
    serviceware-platform = {
      source = "serviceware/swp"
    }
  }
}
*/
provider "swp" {
  aipe_url = "https://saas-ops-1.ai-process-engine.testing.swops.cloud" 
  authenticator_realm_url = "https://auth.sso.testing.swops.cloud/auth/realms/saas-ops-1"

  application_username = local.username
  application_password = local.password
}

/*
data "aipe_data_object" "test_object" {
  id = "436099"
}

output "test_object" {
  value = data.aipe_data_object.test_object
}
*/

locals {
  clusters = {
    "lizzy-labs" = {
      name = "Delizzifyed"
      region = "eu"
    }
    "greta-internal" = {
      name = "Greta"
      region = "eu"
    }
    "patty-production" = {
      name = "Patty"
      region = "us"
    }
  }
}

resource "swp_aipe_data_object" "nomad_cluster" {
  for_each = local.clusters
  type = "nomad-cluster-2"
  properties = {
    "name" = each.value.name,
    "region-2" = each.value.region,
  }
}

resource "swp_aipe_data_object" "lizzy_fsn" {
  type = "nomad-datacenter"
  properties = {
    name = "Lizzy FSN"
  }
}
resource "swp_aipe_data_object" "lizzy_hel1" {
  type = "nomad-datacenter"
  properties = {
    name = "Lizzy HEL1"
  }
}

resource "swp_aipe_data_object_link" "lizzy_has_lizzy_fsn" {
  source_id = swp_aipe_data_object.nomad_cluster["lizzy-labs"].id
  link_name = "nomad-datacenter"
  relation_name = "contains"
  target_ids = [
    swp_aipe_data_object.lizzy_fsn.id,
    swp_aipe_data_object.lizzy_hel1.id,
  ]
}
resource "swp_aipe_data_object" "greta_fsn" {
  type = "nomad-datacenter"
  properties = {
    name = "Greta FSN1"
  }
}
resource "swp_aipe_data_object_link" "greta_has_dcs" {
  source_id = swp_aipe_data_object.nomad_cluster["greta-internal"].id
  link_name = "nomad-datacenter"
  relation_name = "contains"
  target_ids = [
    swp_aipe_data_object.greta_fsn.id,
  ]
}
