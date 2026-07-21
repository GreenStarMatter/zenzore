-- Migration 002: Set up Datastream CDC publication and replication slot
-- Run against the zenzore_registry database on the zenzore-registry instance,
-- AFTER cloudsql.logical_decoding has been enabled and the instance restarted

CREATE PUBLICATION datastream_publication FOR TABLE dim_zyztem, dim_device, dim_sensor;

SELECT PG_CREATE_LOGICAL_REPLICATION_SLOT('datastream_slot', 'pgoutput');
