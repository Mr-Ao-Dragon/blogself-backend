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
  }
}
variable "region" {
  default = "cn-hongkong"
}
variable "AK" {}
variable "SK" {}
provider "alicloud" {
  access_key = var.AK
  secret_key = var.SK
  region = var.region
}
resource "random_string" "name-salt" {
  length = 7
  special = false
  lower = false
  upper = true
}
resource "alicloud_resource_manager_resource_group" "blog" {
  display_name = "blog shelf"
  resource_group_name = "blog-shelf-${random_string.name-salt.result}"
}
output "resource-group-name" {
  value = alicloud_resource_manager_resource_group.blog.resource_group_name
}