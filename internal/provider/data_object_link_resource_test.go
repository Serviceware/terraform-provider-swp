package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var idDiffTests = []struct {
	stateTargetIDs []string
	planTargetIDs  []string

	expectedAdd    []string
	expectedRemove []string
}{
	{
		stateTargetIDs: []string{"a", "b", "c"},
		planTargetIDs:  []string{"b", "c", "d"},
		expectedAdd:    []string{"d"},
		expectedRemove: []string{"a"},
	},
	{
		stateTargetIDs: []string{"a", "b", "c"},
		planTargetIDs:  []string{"a", "b", "c"},
		expectedAdd:    []string{},
		expectedRemove: []string{},
	},
	{
		stateTargetIDs: []string{"a", "b"},
		planTargetIDs:  []string{"a", "b", "c", "d"},
		expectedAdd:    []string{"c", "d"},
		expectedRemove: []string{},
	},
	{
		stateTargetIDs: []string{"a", "b", "c", "d"},
		planTargetIDs:  []string{"a"},
		expectedAdd:    []string{},
		expectedRemove: []string{"b", "c", "d"},
	},
	{
		stateTargetIDs: []string{},
		planTargetIDs:  []string{"a", "b"},
		expectedAdd:    []string{"a", "b"},
		expectedRemove: []string{},
	},
	{
		stateTargetIDs: []string{"a", "b"},
		planTargetIDs:  []string{},
		expectedAdd:    []string{},
		expectedRemove: []string{"a", "b"},
	},
	{
		stateTargetIDs: []string{"a", "b", "c"},
		planTargetIDs:  []string{"d", "e", "f"},
		expectedAdd:    []string{"d", "e", "f"},
		expectedRemove: []string{"a", "b", "c"},
	},
}

func TestDiffStateAndPlanIDs(t *testing.T) {
	for _, tt := range idDiffTests {
		t.Run(fmt.Sprintf("diff(%v, %v)", tt.stateTargetIDs, tt.planTargetIDs), func(t *testing.T) {
			add, remove := diffStateAndPlanIDs(tt.stateTargetIDs, tt.planTargetIDs)

			if len(add) != len(tt.expectedAdd) {
				t.Errorf("expected %d add ids, got %d", len(tt.expectedAdd), len(add))
			}

			if len(remove) != len(tt.expectedRemove) {
				t.Errorf("expected %d remove ids, got %d", len(tt.expectedRemove), len(remove))
			}

			var addAsMap = make(map[string]bool)
			for _, id := range add {
				addAsMap[id] = true
			}

			var removeAsMap = make(map[string]bool)
			for _, id := range remove {
				removeAsMap[id] = true
			}

			for _, id := range tt.expectedAdd {
				if !addAsMap[id] {
					t.Errorf("expected %s to be in add", id)
				}
			}

			for _, id := range tt.expectedRemove {
				if !removeAsMap[id] {
					t.Errorf("expected %s to be in remove", id)
				}
			}
		})
	}
}

func TestAccAIPEDataObjectLinkHappyCase(t *testing.T) {
	var db01 = &DataObject{}
	var db02 = &DataObject{}
	var cloudInc = &DataObject{}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataObjectLinkSimple(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataObjectIDFetch("swp_aipe_data_object.db01", db01),
					testAccDataObjectIDFetch("swp_aipe_data_object.db02", db02),
					testAccDataObjectIDFetch("swp_aipe_data_object.cloud_inc", cloudInc),
					testAccDataObjectHasLinks(cloudInc, linkNameFromAIPE, relationNameFromAIPE, []*DataObject{db01, db02}),
				),
			},
		},
	})
}

var serverObjectTypeFromAIPE = "server"
var hosterObjectTypeFromAIPE = "hoster"
var linkNameFromAIPE = "server-hosted-by-hoster"
var relationNameFromAIPE = "hosts"

func testAccDataObjectLinkSimple() string {
	return fmt.Sprintf(`
resource "swp_aipe_data_object" "db01" {
	type = "%s"
	properties = {
		"fqdn" = "db01"
	}
}

resource "swp_aipe_data_object" "db02" {
	type = "%s"
	properties = {
		"fqdn" = "db02"
	}
}

resource "swp_aipe_data_object" "cloud_inc" {
	type = "%s"
	properties = {
		"support-portal" = "https://support.cloud-inc.example"
	}
}

resource "swp_aipe_data_object_link" "cloud-inc-hosting-both-dbs" {
	source_id = swp_aipe_data_object.cloud_inc.id

	link_name = "%s"
	relation_name = "%s"

	target_ids = [
		swp_aipe_data_object.db01.id,
		swp_aipe_data_object.db02.id,
	]
}
	`,
		serverObjectTypeFromAIPE,
		serverObjectTypeFromAIPE,
		hosterObjectTypeFromAIPE,
		linkNameFromAIPE,
		relationNameFromAIPE)
}

func testAccDataObjectHasLinks(source *DataObject, linkName string, relationName string, expectedLinks []*DataObject) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ids, err := aipeClient.GetDataObjectLinks(context.Background(), source.ID, linkName, relationName)
		if err != nil {
			return err
		}

		if len(ids) != len(expectedLinks) {
			return fmt.Errorf("expected %d ids, got %d", len(expectedLinks), len(ids))
		}

		var expectedIdsAsMap = make(map[string]bool)
		for _, link := range expectedLinks {
			expectedIdsAsMap[link.ID] = true
		}

		for _, id := range ids {
			if !expectedIdsAsMap[id] {
				return fmt.Errorf("expected %s to be in ids", id)
			}
		}
		return nil
	}
}
