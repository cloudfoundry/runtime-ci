variable "env_name" {}

variable "dns_suffix" {}

variable "project" {}

variable "service_account_key" {}

provider "google" {
  project     = "${var.project}"
  region      = "${var.region}"
  credentials = "${var.service_account_key}"
}

variable "projectid" {
    type = "string"
    default = "cf-relint-bosh-lite"
}

variable "region" {
    type = "string"
    default = "us-central1"
}

variable "zone" {
    type = "string"
    default = "us-central1-a"
}

resource "google_compute_network" "bosh-lite" {
  name       = "bosh-lite"
}

// Static IP for the BOSH director
resource "google_compute_address" "bosh-lite" {
  name = "bosh-lite"
  project = "cf-relint-bosh-lite"
  region = "us-central1"
}

// Subnet for the BOSH director
resource "google_compute_subnetwork" "bosh-lite" {
  name          = "bosh-lite"
  ip_cidr_range = "10.0.1.0/16"
  network       = "${google_compute_network.bosh-lite.self_link}"
}

resource "google_dns_managed_zone" "env_dns_zone" {
  name        = "${var.env_name}-zone"
  dns_name    = "${var.env_name}.${var.dns_suffix}."
  description = "DNS zone for the ${var.env_name} environment"
}

resource "google_dns_record_set" "wildcard-dns" {
  name       = "*.${google_dns_managed_zone.env_dns_zone.dns_name}"
  depends_on = ["google_compute_address.bosh-lite"]
  type       = "A"
  ttl        = 300

  managed_zone = "${google_dns_managed_zone.env_dns_zone.name}"

  rrdatas = ["${google_compute_address.bosh-lite.address}"]
}

// Allow ssh & mbus access to director
resource "google_compute_firewall" "bosh-lite" {
  name    = "bosh-lite"
  network = "${google_compute_network.bosh-lite.name}"

  allow {
    protocol = "tcp"
    ports = ["22", "6868", "25555", "80", "443"]
  }

  allow {
    protocol = "icmp"
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags = ["bosh-lite"]
}

// Outputs go below here:
output "external_ip" {
  value = "${google_compute_address.bosh-lite.address}"
}

output "projectid" {
  value = "${var.projectid}"
}

output "credentials" {
  value = "${var.service_account_key}"
}

output "region" {
  value = "${var.region}"
}

output "zone" {
  value = "${var.zone}"
}

output "network_name" {
  value = "${google_compute_network.bosh-lite.name}"
}

output "subnetwork_name" {
  value = "${google_compute_subnetwork.bosh-lite.name}"
}

output "internal_cidr" {
  value = "${google_compute_subnetwork.bosh-lite.ip_cidr_range}"
}

output "internal_gw" {
  value = "${google_compute_subnetwork.bosh-lite.gateway_address}"
}

output "internal_ip" {
  value = "10.0.1.6"
}
