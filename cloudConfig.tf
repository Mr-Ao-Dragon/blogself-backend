terraform {
  required_providers {
    alicloud = {
      source  = "aliyun/alicloud"
      version = "1.242.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.6.2"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "2.7.0"
    }
  }
  cloud {
    organization = "blogshelf"
    workspaces {
      name = "blogshelf"
    }
  }
}
variable "region" {
  default = "cn-hongkong"
}
variable "AK" {
  nullable = false
  sensitive = false
}
variable "SK" {
  nullable = false
  sensitive = false
}
variable "GH_BASIC_CLIENT_ID" {
  nullable = false
  sensitive = false

}
variable "GH_BASIC_SECRET_SECRET" {
  nullable = false
  sensitive = false
}
provider "alicloud" {
  access_key = var.AK
  secret_key = var.SK
  region = var.region
}
resource "random_string" "name-salt" {
  length = 4
  special = false
  lower = true
  upper = false
}
resource "alicloud_resource_manager_resource_group" "blog" {
  display_name = "blog shelf"
  resource_group_name = "blog-shelf-${random_string.name-salt.result}"
}
output "resource-group-name" {
  value = alicloud_resource_manager_resource_group.blog.resource_group_name
}