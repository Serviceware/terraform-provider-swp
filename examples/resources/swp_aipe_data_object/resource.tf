resource "swp_aipe_data_object" "example_server" {
  type = "cloud-server"
  properties = {
    "name" = "db01.example",
    "ip"   = "10.1.2.3"
  }
}
