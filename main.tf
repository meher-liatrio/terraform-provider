 terraform {
   required_providers {
     devops-bootcamp = {
       source = "liatr.io/terraform/devops-bootcamp"
     }
   }
 }

provider "devops-bootcamp" {
    host = "http://localhost:8080"
}

data "devops-bootcamp_engineer" "test" {}

resource "devops-bootcamp_engineer_resource" "test" {
    name = "alice"
    email = "alice@finches.com"
}

output "test_output" {
  value = data.devops-bootcamp_engineer.test
}
