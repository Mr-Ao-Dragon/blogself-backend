resource "alicloud_ots_instance" "blogStore" {
  name              = "blogStore-${random_string.name-salt.result}"
  description       = "存储你的博客数据"
  accessed_by       = "Vpc"
  resource_group_id = alicloud_resource_manager_resource_group.blog.id
  instance_type     = "HighPerformance"
}
output "ots-name" {
  value = alicloud_ots_instance.blogStore.name
}
resource "alicloud_ots_table" "user" {
  instance_name = alicloud_ots_instance.blogStore.name
  max_version   = 1
  table_name    = "user"
  time_to_live  = -1
  allow_update  = true
  primary_key {
    name = "userid"
    type = "String"
  }
  defined_column {
    name = "username"
    type = "String"
  }
  defined_column {
    name = "email"
    type = "String"
  }
  defined_column {
    name = "role"
    type = "Boolean"
  }
  defined_column {
    name = "banned"
    type = "Boolean"
  }
  defined_column {
    name = "allowpost"
    type = "Boolean"
  }
  defined_column {
    name = "githubtoken"
    type = "String"
  }
  defined_column {
    name = "createdat"
    type = "Integer"
  }
}
resource "alicloud_ots_table" "post" {
  instance_name = alicloud_ots_instance.blogStore.name
  max_version   = 1
  table_name    = "post"
  time_to_live  = -1
  allow_update  = true
  primary_key {
    name = "postid"
    type = "String"
  }
  defined_column {
    name = "authorid"
    type = "String"
  }
  defined_column {
    name = "title"
    type = "String"
  }
  defined_column {
    name = "content"
    type = "String"
  }
  defined_column {
    name = "status"
    type = "String" // soft delete
  }
  defined_column {
    name = "createdat"
    type = "Integer"
  }
  defined_column {
    name = "updateat"
    type = "Integer"
  }
}
resource "alicloud_ots_table" "map" {
  instance_name = alicloud_ots_instance.blogStore.name
  max_version   = 1
  table_name    = "map"
  time_to_live  = -1
  allow_update = false
  primary_key {
    name = "postid"
    type = "String"
  }
  defined_column {
    name = "commentid"
    type = "String"
  }
}
resource "alicloud_ots_table" "comment" {
  instance_name = alicloud_ots_instance.blogStore.name
  max_version   = 1
  table_name    = "comment"
  time_to_live  = -1
  allow_update = true
  primary_key {
    name = "postid"
    type = "String"
  }
  primary_key {
    name = "commentid"
    type = "String"
  }
  defined_column {
    name = "authorid"
    type = "String"
  }
  defined_column {
    name = "status" // soft delete
    type = "Boolean"
  }
  defined_column {
    name = "createdat"
    type = "Integer"
  }
}