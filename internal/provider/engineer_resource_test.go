package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEngineerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "devops-bootcamp_engineer_resource" "test" {
	name = "test"
    email = "test@test.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first engineer resource has required attributes filled.
					resource.TestCheckResourceAttr("devops-bootcamp_engineer_resource.test", "name", "test"),
					resource.TestCheckResourceAttr("devops-bootcamp_engineer_resource.test", "email", "test@test.com"),
					// Verify first engineer resource has Computed attributes filled.
					resource.TestCheckResourceAttrSet("devops-bootcamp_engineer_resource.test", "id"),
					resource.TestCheckResourceAttrSet("devops-bootcamp_engineer_resource.test", "last_updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "devops-bootcamp_engineer_resource.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "devops-bootcamp_engineer_resource" "test" {
	name = "test.edit"
    email = "test.edit@test.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first engineer resource has required attributes filled.
					resource.TestCheckResourceAttr("devops-bootcamp_engineer_resource.test", "name", "test.edit"),
					resource.TestCheckResourceAttr("devops-bootcamp_engineer_resource.test", "email", "test.edit@test.com"),
					// Verify first engineer resource has Computed attributes filled.
					resource.TestCheckResourceAttrSet("devops-bootcamp_engineer_resource.test", "id"),
					resource.TestCheckResourceAttrSet("devops-bootcamp_engineer_resource.test", "last_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
