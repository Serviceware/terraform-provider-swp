package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var existingProperty = "foo"
var typoProperty = "oof"

func TestAccAIPEDataObjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testDataObjectDatasource(existingProperty, "bar"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.swp_aipe_data_object.test_object", "properties."+existingProperty, "bar"),
				),
			},
		},
	})
}

func TestAccAIPEDataObjectErrorInFirstCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// IF a typo is made in the property name, the first creation will fail
			{
				Config:      testDataObjectDatasource(typoProperty, "bar"),
				ExpectError: regexp.MustCompile("unexpected status code: 400"),
			},
			// but then, a corrected version should work
			{
				Config: testDataObjectDatasource(existingProperty, "bar"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.swp_aipe_data_object.test_object", "properties."+existingProperty, "bar"),
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
