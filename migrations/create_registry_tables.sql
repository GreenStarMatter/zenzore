-- Migration 001: Create zenzore_registry tables
-- Run against the zenzore_registry database on the zenzore-registry Cloud SQL instance

CREATE DATABASE zenzore_registry;

CREATE TABLE IF NOT EXISTS dim_zyztem (
    zyztem_key SERIAL PRIMARY KEY,
    zyztem_id VARCHAR NOT NULL UNIQUE,
    zyztem_pn VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    zyztem_sn VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    manufacturer VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    model VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    revision VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS dim_device (
    device_key SERIAL PRIMARY KEY,
    device_pn VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    device_sn VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    manufacturer VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    oem_device_pn VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    oem_device_sn VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    firmware_version VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    hardware_revision VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS dim_sensor (
    sensor_key SERIAL PRIMARY KEY,
    sensor_id VARCHAR NOT NULL UNIQUE,
    sensor_pn VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    sensor_sn VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    manufacturer VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    oem_sensor_pn VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    oem_sensor_sn VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    sensor_type VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    engineering_units VARCHAR NOT NULL DEFAULT 'UNKNOWN',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
