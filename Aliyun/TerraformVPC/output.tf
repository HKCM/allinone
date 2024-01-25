output "node_vpc_id" {
  value = alicloud_vpc.node_vpc.id
}

output "vswitch0" {
  value = alicloud_vswitch.vswitch0.id
}

output "vswitch1" {
  value = alicloud_vswitch.vswitch1.id
}

output "vswitch2" {
  value = alicloud_vswitch.vswitch2.id
}

output "node_rule" {
  value = alicloud_security_group.node_sg.id
}