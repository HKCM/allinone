
# 由于module不能正确创建instance， 这里手动创建
# The instance type HighPerformance is not available in the region ap-northeast-1.
resource "alicloud_ots_instance" "state_lock_instance" {
  name          = "tf-lock-vpc"
  description   = "Terraform remote backend state lock."
  accessed_by   = "Any"
  instance_type = "Capacity"
  tags = {
    Purpose = "Terraform state lock for state "
  }
}

module "remote_state" {
  source                    = "terraform-alicloud-modules/remote-backend/alicloud"
  create_backend_bucket     = true
  create_ots_lock_instance  = false
  create_ots_lock_table     = true
  backend_ots_lock_instance = alicloud_ots_instance.state_lock_instance.id
  region                    = var.region
  state_name                = "terraform.tfstate"
  encrypt_state             = false

}