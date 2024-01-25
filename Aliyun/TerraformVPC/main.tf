provider "alicloud" {}

resource "alicloud_vpc" "node_vpc" {
  vpc_name   = var.vpc_name
  cidr_block = var.vpc_cidr_block
}

resource "alicloud_vswitch" "vswitch0" {
  zone_id      = var.az[0]
  vswitch_name = var.vswitch_name[0]
  cidr_block   = var.vswitch_cidr_blocks[0]
  vpc_id       = alicloud_vpc.node_vpc.id
}

resource "alicloud_vswitch" "vswitch1" {
  zone_id      = var.az[1]
  vswitch_name = var.vswitch_name[1]
  cidr_block   = var.vswitch_cidr_blocks[1]
  vpc_id       = alicloud_vpc.node_vpc.id
}

resource "alicloud_vswitch" "vswitch2" {
  zone_id      = var.az[2]
  vswitch_name = var.vswitch_name[2]
  cidr_block   = var.vswitch_cidr_blocks[2]
  vpc_id       = alicloud_vpc.node_vpc.id
}

resource "alicloud_security_group" "node_sg" {
  name   = "node_sg"
  vpc_id = alicloud_vpc.node_vpc.id
}

resource "alicloud_security_group_rule" "node_rule" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "2522/2522"
  priority          = 1
  security_group_id = alicloud_security_group.node_sg.id
  cidr_ip           = "150.249.195.73/32"
  description       = "For okj office ip"
}

resource "alicloud_security_group_rule" "aptos_node_rule1" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "6180/6180"
  priority          = 1
  security_group_id = alicloud_security_group.node_sg.id
  cidr_ip           = "0.0.0.0/0"
  description       = "For aptos node"
}

resource "alicloud_security_group_rule" "aptos_node_rule2" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "6182/6182"
  priority          = 1
  security_group_id = alicloud_security_group.node_sg.id
  cidr_ip           = "0.0.0.0/0"
  description       = "For aptos node"
}