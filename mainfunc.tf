resource "alicloud_fcv3_function" "blogShelf" {
  code {
    oss_bucket_name = alicloud_oss_bucket.distStore.id
    oss_object_name = "blogshelf-backend"
  }
  depends_on = [alicloud_oss_bucket_object.dist]

  handler = "blogshelf-backend"
  runtime = "custom.debian10"
  environment_variables = {
    GIN_MODE = "release"
    OTS_NAME = alicloud_ots_instance.blogStore.id
  }
  custom_runtime_config {
    port = 3000
    args = ["-prod"]
    command = ["blogshelf-backend"]
  }
  role = alicloud_ram_role.blogFunc.arn
  internet_access      = true
  cpu                  = "0.5"
  memory_size          = 512
  function_name        = "blog"
  disk_size            = 512
  timeout              = 20
  instance_concurrency = 50
}
resource "alicloud_oss_bucket" "distStore" {
  bucket = "bsdist-${random_string.name-salt.result}"
  resource_group_id = alicloud_resource_manager_resource_group.blog.id
}
resource "alicloud_oss_bucket_object" "dist" {
  bucket = alicloud_oss_bucket.distStore.id
  key    = "blogshelf-backend"
  source = "./dist/blogshelf-backend"
  content_type = "application/octet-stream"
}
resource "alicloud_oss_bucket_acl" "dist" {
  acl    = "private"
  bucket = alicloud_oss_bucket.distStore.bucket
}
resource "alicloud_ram_group" "blog" {
  name = "blogshelf"
}
resource "alicloud_ram_user" "blogFunc" {
  name = "blogShelf-MainFunc"
}
resource "alicloud_ram_group_membership" "allSystemUsers" {
  depends_on = [alicloud_ram_group.blog]
  group_name = alicloud_ram_group.blog.id
  user_names = [
  alicloud_ram_user.blogFunc.name
  ]
}
resource "alicloud_ram_policy" "blogDataFullAccess" {
  description     = "允许访问博客全部数据"
  depends_on = [alicloud_ots_table.comment,alicloud_ots_table.map,alicloud_ots_table.user,alicloud_ots_table.post]
  policy_document = "{\"Version\":\"1\",\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"ots:*Row\",\"ots:GetRange\"],\"Resource\":[\"acs:ots:*:*:instance/blogStore-*/table/user\",\"acs:ots:*:*:instance/blogStore-*/table/post\",\"acs:ots:*:*:instance/blogStore-*/table/map\",\"acs:ots:*:*:instance/blogStore-*/table/comment\"]},{\"Effect\":\"Deny\",\"Action\":[\"ots:DeleteRow\"],\"Resource\":[\"acs:ots:*:*:instance/blogStore-*/table/user\",\"acs:ots:*:*:instance/blogStore-*/table/post\"]}]}"
  policy_name     = "blogFunc"
}
resource "alicloud_ram_role" "blogFunc" {
  name = "blogshelf-mainfunc"
  document = <<POLICY
  {
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": ["fc.aliyuncs.com"]
      }
    }
  ],
  "Version": "1"
}
POLICY
}
resource "alicloud_ram_role_policy_attachment" "blogFunc" {
  policy_name = alicloud_ram_policy.blogDataFullAccess.id
  policy_type = "Custom"
  role_name   = alicloud_ram_role.blogFunc.name
}