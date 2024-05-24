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
    name = "meher"
    email = "meher@finches.com"
}


resource "devops-bootcamp_dev_resource" "dev_finches" {
        name = "dev_finches"
        engineers = [
          {id=devops-bootcamp_engineer_resource.meher.id}
          ]     
}

#  terraform plan -generate-config-out=generated.tf
# import {
#   to = devops-bootcamp_engineer_resource.ryan
#   id = "UWJVB"
# }

