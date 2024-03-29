terraform {
  backend "oss" {
    bucket              = "terraform-remote-backend-f1736ed0-fc32-a253-9b4a-f0dadba09fd2"
    prefix              = ""
    key                 = "terraform.tfstate"
    acl                 = "private"
    region              = "ap-northeast-1"
    encrypt             = "false"
    tablestore_endpoint = "https://tf-lock-vpc.ap-northeast-1.ots.aliyuncs.com"
    tablestore_table    = "terraform_remote_backend_lock_table_f1736ed0_fc32_a253_9b4a_f0dadba09fd2"
  }
}

