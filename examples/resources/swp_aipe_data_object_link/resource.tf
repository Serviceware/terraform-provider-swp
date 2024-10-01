resource "swp_aipe_data_object" "example_server" {
  type = "cloud-server"
  properties = {
    "name" = "db01.example",
    "ip"   = "10.1.2.3"
  }
}


resource "swp_aipe_data_object" "example_hoster" {
  type = "hoster"
  properties = {
    "name"            = "Cloud Provider Inc.",
    "support_contact" = "https://support.cloudprovider.example"
  }
}


resource "swp_aipe_data_object_link" "server_is_hosted" {
  source_id = swp_aipe_data_object.example_server.id

  link          = "servers-to-hoster"
  relation_name = "hosted-by"

  target_ids = [
    swp_aipe_data_object.example_hoster.id
  ]
}
