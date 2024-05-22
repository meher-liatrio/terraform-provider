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


resource "devops-bootcamp_engineer_resource" "meher" {
    name = "meherL"
    email = "meher@finches.com"
}

resource "devops-bootcamp_engineer_resource" "sloane" {
    name = "sloane"
    email = "sloane@finches.com"
}

resource "devops-bootcamp_dev_resource" "dev_finches" {
        name = "dev_finches"
        engineers = [
          {id=devops-bootcamp_engineer_resource.meher.id},
          {id=devops-bootcamp_engineer_resource.sloane.id}
          ]     
}

data "devops-bootcamp_devs" "devs" {}

# output "finches_output" {
#   value = devops-bootcamp_dev_resource.dev_finches
# }

output "meher_output" {
  value = devops-bootcamp_engineer_resource.meher
}