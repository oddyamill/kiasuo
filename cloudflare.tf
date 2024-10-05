terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4"
    }
  }
}

variable "cloudflare_zone_id" {
  type = string
}

locals {
  subdomain = "kiasuo-proxy"
  records = [
    {
      name    = local.subdomain
      content = "188.114.96.3"
      comment = "kiasuo proxy"
    },
    {
      name    = "cloudflare-edge-kja"
      content = "172.70.152.140"
      comment = "data center Krasnoyarsk"
    },
    {
      name    = "cloudflare-edge-kld"
      content = "172.71.17.138"
      comment = "data center Tver"
    },
    {
      name    = "cloudflare-edge-led"
      content = "172.69.8.205"
      comment = "data center Saint Petersburg"
    },
  ]
}

resource "cloudflare_record" "dns" {
  for_each = { for idx, record in local.records : idx => record }

  zone_id         = var.cloudflare_zone_id
  name            = each.value.name
  content         = each.value.content
  comment         = each.value.comment
  type            = "A"
  proxied         = true
  allow_overwrite = true
}

resource "cloudflare_ruleset" "firewall" {
  zone_id = var.cloudflare_zone_id
  name = "Kiasuo ruleset"
  kind    = "zone"
  phase   = "http_request_firewall_custom"

  rules {
    action      = "block"
    expression  = "(http.host eq \"${local.subdomain}.oddya.ru\" and starts_with(http.request.uri.path, \"/internal/\") and http.request.uri.path ne \"/internal/purge-cache\")"
    description = "Block requests to /internal/ except /internal/purge-cache"
    enabled     = true
  }
}
