resource "archive_file" "zip" {
  source_file      = "dist/blogshelf-bakcend"
  output_path      = "dist/code.zip"
  type             = "zip"
  output_file_mode = "0444"
}
resource "alicloud_fcv3_function" "blogShelf" {
  code {
    zip_file = filebase64(tostring(archive_file.zip.output_path))
  }
  handler = "blogshelf-backend"
  runtime = "custom.debian10"
  environment_variables = {
    GIN_MODE = "release"
    OTS_NAME = alicloud_ots_instance.blogStore.name
  }
  custom_runtime_config {
    port = 3000
    args = []
    command = ["blogshelf-backend"]
  }
  role = alicloud_ram_role.blogFunc.name
  internet_access      = true
  cpu                  = "0.5"
  memory_size          = 512
  function_name        = "blog"
  disk_size            = 512
  timeout              = 20
  instance_concurrency = 50
}
resource "alicloud_ram_group" "blog" {
  name = "blogshelf"
}
resource "alicloud_ram_user" "blogFunc" {

  name = "blogShelf-MainFunc"
}
resource "alicloud_ram_group_membership" "allSystemUsers" {
  group_name = "加入所有成员到用户组"
  user_names = [
  alicloud_ram_user.blogFunc.name
  ]
}
resource "alicloud_ram_policy" "blogDataFullAccess" {
  description     = "允许访问博客全部数据"
  policy_document = "{\"Version\":\"1\",\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"ots:*Row\",\"ots:GetRange\",],\"Resource\":[\"acs:ots:*:*:instance/blogStore-*/table/user\",\"acs:ots:*:*:instance/blogStore-*/table/post\",\"acs:ots:*:*:instance/blogStore-*/table/map\",\"acs:ots:*:*:instance/blogStore-*/table/comment\"],\"Condition\":{}},{\"Effect\":\"Deny\",\"Action\":[\"ots:DeleteRow\"],\"Resource\":[\"acs:ots:*:*:instance/blogStore-*/table/user\",\"acs:ots:*:*:instance/blogStore-*/table/post\"],\"Condition\":{}}]}"
  policy_name     = "blogFunc"
}
resource "alicloud_ram_role" "blogFunc" {
  name = "博客主体角色"
  document = <<EOF
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "RAM": [
          "acs:ram::*:blogShelf-MainFunc"
        ]
      }
    }
  ],
  "Version": "1"
}
EOF
}
resource "alicloud_ram_role_policy_attachment" "blogFunc" {
  policy_name = alicloud_ram_policy.blogDataFullAccess.name
  policy_type = "Custom"
  role_name   = alicloud_ram_role.blogFunc.name
}