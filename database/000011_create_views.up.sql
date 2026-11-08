BEGIN;

CREATE OR REPLACE VIEW v_measurement_points AS
SELECT
    mp.id AS point_id,
    mp.code AS point_code,
    mp.name AS point_name,
    mp.point_type,
    mp.value_type,
    mp.unit,
    mp.precision,
    mp.status AS point_status,

    e.id AS equipment_id,
    e.code AS equipment_code,
    e.name AS equipment_name,

    a.id AS area_id,
    a.code AS area_code,
    a.name AS area_name,

    p.id AS plant_id,
    p.code AS plant_code,
    p.name AS plant_name

FROM measurement_points mp

LEFT JOIN equipments e
    ON e.id = mp.equipment_id

LEFT JOIN areas a
    ON a.id = e.area_id

LEFT JOIN plants p
    ON p.id = a.plant_id

WHERE mp.deleted_at IS NULL
  AND e.deleted_at IS NULL;

CREATE OR REPLACE VIEW v_active_point_bindings AS
SELECT
    pb.id AS binding_id,
    pb.point_id,
    mp.code AS point_code,
    mp.name AS point_name,
    mp.value_type,
    mp.unit,
    mp.deadband,

    pb.source_id,
    ds.name AS source_name,
    ds.protocol,
    ds.endpoint,

    pb.address_config,
    pb.conversion_config

FROM point_bindings pb

JOIN measurement_points mp
    ON mp.id = pb.point_id

JOIN data_sources ds
    ON ds.id = pb.source_id

WHERE pb.enabled = TRUE
  AND mp.status = 1
  AND mp.deleted_at IS NULL
  AND ds.status = 1;

COMMIT;
