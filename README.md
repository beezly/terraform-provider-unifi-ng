# terraform-provider-unifi

A Terraform provider for UniFi Network, built on the **official UniFi Integration API** (OpenAPI-based) rather than the legacy controller API.

## Why a new provider?

Existing providers ([paultyng/unifi](https://github.com/paultyng/terraform-provider-unifi), [ubiquiti-community/unifi](https://github.com/ubiquiti-community/terraform-provider-unifi)) use the legacy `/api/s/{site}/rest/` API. This API has become increasingly unreliable on newer UniFi OS firmware (5.x+), with type mismatches causing import failures.

This provider uses the new Integration API (`/v1/sites/{siteId}/...`), which:
- Supports **API key authentication** (no username/password)
- Is **officially documented** via OpenAPI specs
- Is **more stable** across firmware versions

OpenAPI specs sourced from: **<https://beez.ly/unifi-apis/>**

## Resources

| Resource | API Endpoint |
|----------|-------------|
| `unifi_network` | `/v1/sites/{siteId}/networks` |
| `unifi_firewall_policy` | `/v1/sites/{siteId}/firewall/policies` |
| `unifi_firewall_zone` | `/v1/sites/{siteId}/firewall/zones` |
| `unifi_dns_policy` | `/v1/sites/{siteId}/dns/policies` |
| `unifi_traffic_matching_list` | `/v1/sites/{siteId}/traffic-matching-lists` |

## Data Sources

| Data Source | API Endpoint |
|------------|-------------|
| `unifi_sites` | `/v1/sites` |
| `unifi_devices` | `/v1/sites/{siteId}/devices` |

## Usage

```hcl
provider "unifi" {
  api_key        = var.unifi_api_key   # or UNIFI_API_KEY env var
  api_url        = "https://192.168.0.1"
  site_id        = "88f7af54-98f8-306a-a1c7-c9349722b1f6"
  allow_insecure = true  # for self-signed certs
}

resource "unifi_network" "iot" {
  name    = "IoT"
  vlan_id = 3
}
```

## Build

```bash
# Install codegen tools (first time only)
go install github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi@latest
go install github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework@latest

# Full pipeline: fetch latest spec → normalize → generate → build
make update

# Or step by step:
make fetch      # download latest spec from beez.ly/unifi-apis/
make normalize  # fix schema names for codegen compatibility
make generate   # run codegen tools
make build      # compile binary
```

## Updating for new UniFi API versions

When Ubiquiti releases a new UniFi Network version with an updated OpenAPI spec on <https://beez.ly/unifi-apis/>:

```bash
make update
```

That's it. The pipeline will detect the new version, regenerate all schema types, and rebuild.

## Status

🚧 **Work in progress** — schema generation is working, HTTP CRUD implementations in progress.
