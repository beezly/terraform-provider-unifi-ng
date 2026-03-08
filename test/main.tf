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

# Read-only: list all sites
data "unifi_sites" "all" {}

# Read-only: list all devices
data "unifi_devices" "all" {}

output "sites" {
  value = data.unifi_sites.all
}

output "devices" {
  value     = data.unifi_devices.all
  sensitive = true
}
