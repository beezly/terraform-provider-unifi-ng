terraform {
  required_providers {
    unifi = {
      source  = "beezly/unifi"
      version = "0.0.1"
    }
  }
}

provider "unifi" {
  api_key        = var.unifi_api_key
  api_url        = "https://192.168.0.1"
  site_id        = "88f7af54-98f8-306a-a1c7-c9349722b1f6"
  allow_insecure = true
}

variable "unifi_api_key" {
  sensitive = true
}

# --- Data sources ---

data "unifi_sites" "all" {}
data "unifi_devices" "all" {}

output "sites" {
  value = data.unifi_sites.all
}

output "devices" {
  value     = data.unifi_devices.all
  sensitive = true
}

# --- Networks ---

resource "unifi_network" "admin" {
  name         = "Admin"
  enabled      = true
  vlan_id      = 1
  ip_subnet    = "192.168.0.0/24"
  dhcp_enabled = true
}

resource "unifi_network" "clients" {
  name         = "Clients"
  enabled      = true
  vlan_id      = 2
  ip_subnet    = "192.168.2.0/24"
  dhcp_enabled = true
}

resource "unifi_network" "iot" {
  name         = "iot"
  enabled      = true
  vlan_id      = 3
  ip_subnet    = "192.168.3.0/24"
  dhcp_enabled = true
}

resource "unifi_network" "cameras" {
  name         = "Cameras"
  enabled      = true
  vlan_id      = 4
  ip_subnet    = "192.168.4.0/24"
  dhcp_enabled = true
}

resource "unifi_network" "infra" {
  name         = "Infra"
  enabled      = true
  vlan_id      = 6
  ip_subnet    = "192.168.6.0/24"
  dhcp_enabled = true
}

resource "unifi_network" "leo" {
  name         = "Leo"
  enabled      = true
  vlan_id      = 101
  ip_subnet    = "192.168.253.0/24"
  dhcp_enabled = true
}

resource "unifi_network" "inter_vlan" {
  name         = "Inter-VLAN routing"
  enabled      = true
  vlan_id      = 4040
  ip_subnet    = "10.255.253.0/24"
  dhcp_enabled = false
}

resource "unifi_network" "isp_stix" {
  name    = "ISP-Stix"
  enabled = true
  vlan_id = 1002
}
