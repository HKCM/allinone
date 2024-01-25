variable "vpc_name" {
  type        = string
  default     = "prod_vpc"
  description = "vpc_name"
}

variable "region" {
  type        = string
  default     = "ap-northeast-1"
  description = "Alibaba region"
}

variable "vpc_cidr_block" {
  type        = string
  default     = "10.52.0.0/16"
  description = "vpc_cidr_block"
}

variable "vswitch_cidr_blocks" {
  description = "Available cidr blocks for public subnets."
  type        = list(string)
  default = [
    "10.52.1.0/24",
    "10.52.2.0/24",
    "10.52.3.0/24",
  ]
}

variable "az" {
  description = "Available cidr blocks for public subnets."
  type        = list(string)
  default = [
    "ap-northeast-1a",
    "ap-northeast-1b",
    "ap-northeast-1c",
  ]
}

variable "vswitch_name" {
  description = "Available cidr blocks for public subnets."
  type        = list(string)
  default = [
    "node_vswith_1a",
    "node_vswith_1b",
    "node_vswith_1c",
  ]
}