package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func TestAccAIPEDataObjectDeletedManually(t *testing.T) {
	var object = &DataObject{}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testDataObjectDatasource(existingProperty, "bar"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.swp_aipe_data_object.test_object", "properties."+existingProperty, "bar"),
					testAccDataObjectIDFetch("swp_aipe_data_object.test_object", object),
				),
			},
			{
				PreConfig: func() {
					err := aipeClient.DeleteObject(context.Background(), object.ID)
					if err != nil {
						t.Fatalf("Failed to delete object: %s", err)
					}
				},
				Config: testDataObjectDatasource(existingProperty, "bar"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.swp_aipe_data_object.test_object", "properties."+existingProperty, "bar"),
				),
			},
		},
	})
}

type DataObject struct {
	ID string
}

func testAccDataObjectIDFetch(resourceName string, object *DataObject) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			object.ID = ""
			return fmt.Errorf("Not found: %s", resourceName)
		}

		object.ID = resource.Primary.ID
		return nil
	}
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
