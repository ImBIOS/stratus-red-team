terraform {
    required_providers {
        google = {
            source  = "hashicorp/google"
            version = "~> 4.28.0"
        }
        google-beta = {
            source  = "hashicorp/google-beta"
            version = "~> 4.28.0"
        }
    }
}
locals {
    num_secrets     = 20
    resource_prefix = "stratus-red-team-retrieve-secret"
}
resource "random_string" "secrets" {
    count       = local.num_secrets
    length      = 16
    min_lower   = 16
}

resource "google_secret_manager_secret" "secrets" {
    provider    = google-beta
    count       = local.num_secrets
    secret_id   = "${local.resource_prefix}-${count.index}"

    replication {
        automatic = true
    }
    depends_on  = [google_project_service.secretmanager]
}

resource "google_secret_manager_secret_version" "secret-values" {
    count       = local.num_secrets
    secret      = google_secret_manager_secret.secrets[count.index].id
    secret_data = random_string.secrets[count.index].result
}

output "display" {
    value = format("%s Secrets Manager secrets ready", local.num_secrets)
}