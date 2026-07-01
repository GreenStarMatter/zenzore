variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "Default region"
  type        = string
  default     = "us-central1"
}

variable "pubsub_topic_name" {
  description = "Pub/Sub topic name"
  type        = string
}


variable "pubsub_dead_topic_name" {
  description = "Pub/Sub deadletter topic name"
  type        = string
}

variable "pubsub_subscription_name" {
  description = "Pub/Sub subscription name"
  type        = string
}

variable "bigquery_dataset_id" {
  description = "BigQuery dataset name"
  type        = string
}


variable "bigquery_table_id" {
  description = "BigQuery table name"
  type        = string
}

variable "storage_bucket_name" {
  description = "Cloud Storage bucket name"
  type        = string
}
