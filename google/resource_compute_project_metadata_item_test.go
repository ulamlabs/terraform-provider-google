package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeProjectMetadataItem_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectMetadataItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectMetadataItem_basic("myKey", "myValue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectMetadataItem_hasMetadata("myKey", "myValue"),
				),
			},
		},
	})
}

func TestAccComputeProjectMetadataItem_basicWithEmptyVal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectMetadataItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectMetadataItem_basic("myKey", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectMetadataItem_hasMetadata("myKey", ""),
				),
			},
		},
	})
}

func TestAccComputeProjectMetadataItem_basicUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectMetadataItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectMetadataItem_basic("myKey", "myValue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectMetadataItem_hasMetadata("myKey", "myValue"),
				),
			},
			{
				Config: testAccProjectMetadataItem_basic("myKey", "myUpdatedValue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectMetadataItem_hasMetadata("myKey", "myUpdatedValue"),
				),
			},
		},
	})
}

func testAccCheckProjectMetadataItem_hasMetadata(key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		project, err := config.clientCompute.Projects.Get(config.Project).Do()
		if err != nil {
			return err
		}

		metadata := flattenComputeMetadata(project.CommonInstanceMetadata.Items)

		val, ok := metadata[key]
		if !ok {
			return fmt.Errorf("Unable to find a value for key '%s'", key)
		}
		if val != value {
			return fmt.Errorf("Value for key '%s' does not match. Expected '%s' but found '%s'", key, value, val)
		}
		return nil
	}
}

func testAccCheckProjectMetadataItemDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	project, err := config.clientCompute.Projects.Get(config.Project).Do()
	if err != nil {
		return err
	}

	metadata := flattenComputeMetadata(project.CommonInstanceMetadata.Items)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_project_metadata_item" {
			continue
		}

		_, ok := metadata[rs.Primary.ID]
		if ok {
			return fmt.Errorf("Metadata key/value '%s': '%s' still exist", rs.Primary.Attributes["key"], rs.Primary.Attributes["value"])
		}
	}

	return nil
}

func testAccProjectMetadataItem_basic(key, val string) string {
	return fmt.Sprintf(`
resource "google_compute_project_metadata_item" "foobar" {
  key   = "%s"
  value = "%s"
}
`, key, val)
}
