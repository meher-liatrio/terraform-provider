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

# data "devops-bootcamp_engineer" "test" {}

data "devops-bootcamp_devs" "finches" {}

resource "devops-bootcamp_engineer_resource" "alice" {
    name = "alice"
    email = "alice@finches.com"
}

import {
  to = devops-bootcamp_engineer_resource.bob
  id = "AJK5Y"
}

resource "devops-bootcamp_engineer_resource" "bob" {
  name = "bob ross"
  email = "bob@ross.com"
}

resource "devops-bootcamp_engineer_resource" "eve" {
    name = "eve"
    email = "eve@finches.com"
}


# resource "devops-bootcamp_engineer_resource" "sloane" {
#     name = "sloane"
#     email = "sloane@finches.com"
# }

resource "devops-bootcamp_engineer_resource" "grant" {
    name = "grant"
    email = "grant@finches.com"
}

resource "devops-bootcamp_engineer_resource" "myles" {
    name = "myles"
    email = "myles@finches.com"
}

# resource "devops-bootcamp_engineer_resource" "meher" {
#     name = "meher"
#     email = "meher@finches.com"
# }

resource "devops-bootcamp_dev_resource" "dev_finches" {
        name = "dev_finches"
        engineers = [
          # {id=devops-bootcamp_engineer_resource.meher.id},
          {id=devops-bootcamp_engineer_resource.myles.id},
          {id=devops-bootcamp_engineer_resource.grant.id},
          {id=devops-bootcamp_engineer_resource.bob.id}
          ]     
}

resource "devops-bootcamp_dev_resource" "dev_cryptos" {
        name = "dev_cryptos"
        engineers = [
          {id=devops-bootcamp_engineer_resource.eve.id},
          {id=devops-bootcamp_engineer_resource.bob.id},
          {id=devops-bootcamp_engineer_resource.alice.id}
          ]     
}
#  terraform plan -generate-config-out=generated.tf
# import {
#   to = devops-bootcamp_engineer_resource.ryan
#   id = "UWJVB"
# }

# import {
#   to = devops-bootcamp_engineer_resource.zach
#   id = "D6LWM"
# }


# output "test_output" {
#   value = devops-bootcamp_engineer_resource.test
# }

# output "test_output2" {
#   value = devops-bootcamp_engineer_resource.bob
# }

# output "test_output3" {
#   value = devops-bootcamp_engineer_resource.eve
# }
