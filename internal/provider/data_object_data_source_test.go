package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAIPEDataObjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testDataObjectDatasource("foo", "bar"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.swp_aipe_data_object.test_object", "properties.foo", "bar"),
				),
			},
		},
	})
}

func testDataObjectDatasource(propertyName string, propertyValue string) string {
	return fmt.Sprintf(`
resource "swp_aipe_data_object" "test_object" {
	type = "test-object"
	properties = {
		%s = "%s"
	}
}

data "swp_aipe_data_object" "test_object" {
  id = swp_aipe_data_object.test_object.id
}
`, propertyName, propertyValue)
}
