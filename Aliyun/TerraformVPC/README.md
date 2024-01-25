# Terrafrom node VPC for alibaba cloud

This terraform repo will create vpc subnets and all necessary components.
This project also create a common bucket to store terraform state.

# Terraform remote state

Use alibaba oss bucket to store state.
When initializing the project, remove the backend config in terraform.tf file, and use terraform init to initialize the project.
That will initialize the bucket to store state.
Then, add the backend config back to terraform.tf file, and use terraform init to initialize the project again.

# Terraform backend config

https://registry.terraform.io/modules/terraform-alicloud-modules/remote-backend
For the first run , You need to delete the terraform.tf file to avoid remote state error, and use terraform apply locally to initialize remote state bucket and dynamodb, then run terraform init `-force-copy` to copy your local state to remote.

# Config credentials

https://registry.terraform.io/providers/aliyun/alicloud/latest/docs
use environment variables to config credentials

```bash
$ export ALICLOUD_ACCESS_KEY="anaccesskey"
$ export ALICLOUD_SECRET_KEY="asecretkey"
$ export ALICLOUD_REGION="cn-beijing"
terraform plan
terraform apply
```

## Problem

> read tcp xxx.xxx.xxx.xxx:xxxx->100.100.1.231:443: read: connection reset by peer

https://github.com/hashicorp/terraform-provider-aws/issues/23614
