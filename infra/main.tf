data "google_project" "project" {}

# Enable required APIs

resource "google_project_service" "sqladmin" {
  service            = "sqladmin.googleapis.com"
  disable_on_destroy = false
}

resource "google_sql_database_instance" "zenzore_registry" {
  name                = "zenzore-registry"
  database_version    = "POSTGRES_15"
  region              = var.region
  deletion_protection = false

  settings {
    tier              = "db-f1-micro"
    activation_policy = "ALWAYS"
  }

  lifecycle {
    ignore_changes = [settings[0].activation_policy]
  }

  depends_on = [google_project_service.sqladmin]
}

resource "google_sql_database" "zenzore_registry" {
  name     = "zenzore_registry"
  instance = google_sql_database_instance.zenzore_registry.name
}

resource "google_sql_user" "registry_admin" {
  name     = "registry_admin"
  instance = google_sql_database_instance.zenzore_registry.name
  password = var.cloudsql_password
}

resource "google_pubsub_subscription_iam_member" "dead_letter_subscriber" {
  subscription = google_pubsub_subscription.zenzore_ingest_bq_sub.name
  role         = "roles/pubsub.subscriber"
  member       = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

resource "google_storage_bucket_iam_member" "pubsub_gcs_writer" {
  bucket = google_storage_bucket.storage_bucket_name.name
  role   = "roles/storage.objectCreator"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

resource "google_storage_bucket_iam_member" "pubsub_gcs_reader" {
  bucket = google_storage_bucket.storage_bucket_name.name
  role   = "roles/storage.legacyBucketReader"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

resource "google_pubsub_topic_iam_member" "dead_letter_publisher" {
  topic  = google_pubsub_topic.zenzore_deadletter.name
  role   = "roles/pubsub.publisher"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

resource "google_bigquery_dataset_iam_member" "pubsub_bq_editor" {
  dataset_id = google_bigquery_dataset.zenzore_raw.dataset_id
  role       = "roles/bigquery.dataEditor"
  member     = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"

  depends_on = [google_bigquery_dataset.zenzore_raw]
}

resource "google_project_service" "pubsub" {
  service            = "pubsub.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "bigquery" {
  service            = "bigquery.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "storage" {
  service            = "storage.googleapis.com"
  disable_on_destroy = false
}

resource "google_bigquery_dataset" "zenzore_raw" {
  dataset_id = var.bigquery_dataset_id
  location   = var.region

  depends_on = [google_project_service.bigquery]
}

resource "google_bigquery_table" "events" {
  dataset_id          = google_bigquery_dataset.zenzore_raw.dataset_id
  table_id            = var.bigquery_table_id
  deletion_protection = false
  depends_on          = [google_bigquery_dataset.zenzore_raw]

  schema = jsonencode([
    {

      name        = "ID"
      type        = "STRING"
      mode        = "NULLABLE"
      description = "Zyztem ID"
    },
    {
      name        = "Devices"
      type        = "RECORD"
      mode        = "REPEATED"
      description = "List of devices in this message"
      fields = [
        {
          name        = "DeviceSN"
          type        = "STRING"
          mode        = "NULLABLE"
          description = "Device serial number"
        },
        {
          name        = "DevicePN"
          type        = "STRING"
          mode        = "NULLABLE"
          description = "Device part number"
        },
        {
          name        = "DeviceExportTime"
          type        = "STRING"
          mode        = "NULLABLE"
          description = "Time the device exported this record"
        },
        {
          name        = "Sensors"
          type        = "RECORD"
          mode        = "REPEATED"
          description = "List of sensors on this device"
          fields = [
            {
              name        = "SensorSN"
              type        = "STRING"
              mode        = "NULLABLE"
              description = "Sensor serial number"
            },
            {
              name        = "SensorPN"
              type        = "STRING"
              mode        = "NULLABLE"
              description = "Sensor part number"
            },
            {
              name        = "SampleValue"
              type        = "FLOAT"
              mode        = "NULLABLE"
              description = "Sensor sample reading"
            },
            {
              name        = "SampleValueTime"
              type        = "STRING"
              mode        = "NULLABLE"
              description = "Time the sample was taken"
            }
          ]
        }
      ]
    }
  ])

}

# Cloud Storage bucket for dead letter messages
resource "google_storage_bucket" "storage_bucket_name" {
  name                     = var.storage_bucket_name
  location                 = var.region
  force_destroy            = true
  public_access_prevention = "enforced"

  depends_on = [google_project_service.storage]
}

# Pub/Sub topic
resource "google_pubsub_topic" "zenzore_ingest" {
  name = var.pubsub_topic_name

  depends_on = [google_project_service.pubsub]
}


resource "google_pubsub_topic" "zenzore_deadletter" {
  name = var.pubsub_dead_topic_name

  depends_on = [google_project_service.pubsub]
}

# Pub/Sub subscription → BigQuery
resource "google_pubsub_subscription" "zenzore_ingest_bq_sub" {
  name  = var.pubsub_subscription_name
  topic = google_pubsub_topic.zenzore_ingest.name

  bigquery_config {
    table            = "${var.project_id}.${var.bigquery_dataset_id}.${var.bigquery_table_id}"
    use_table_schema = true
  }

  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.zenzore_deadletter.id
    max_delivery_attempts = 5
  }

  depends_on = [
    google_project_service.pubsub,
    google_bigquery_table.events,
    google_bigquery_dataset_iam_member.pubsub_bq_editor
  ]
}


# Pub/Sub dead-letter subscription → Cloud Storage
resource "google_pubsub_subscription" "zenzore_deadletter_gcs_sub" {
  name  = "${var.pubsub_topic_name}-deadletter-gcs-sub"
  topic = google_pubsub_topic.zenzore_deadletter.name

  cloud_storage_config {
    bucket = google_storage_bucket.storage_bucket_name.name

    filename_prefix = "deadletter-"
    filename_suffix = ".json"

    max_bytes    = 1000000
    max_duration = "300s"
  }

  depends_on = [
    google_pubsub_topic.zenzore_deadletter,
    google_storage_bucket.storage_bucket_name,
    google_pubsub_topic_iam_member.dead_letter_publisher,
    google_storage_bucket_iam_member.pubsub_gcs_writer,
    google_storage_bucket_iam_member.pubsub_gcs_reader
  ]
}
