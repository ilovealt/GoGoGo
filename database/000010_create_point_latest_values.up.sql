BEGIN;

CREATE TABLE point_latest_values (
    point_id BIGINT PRIMARY KEY,

    ts TIMESTAMPTZ NOT NULL,
    value_type SMALLINT NOT NULL DEFAULT 1,
    value_num DOUBLE PRECISION,
    value_text TEXT,
    value_bool BOOLEAN,
    quality SMALLINT NOT NULL DEFAULT 0,
    source_id BIGINT,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_point_latest_values_point
        FOREIGN KEY (point_id)
        REFERENCES measurement_points(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_point_latest_values_source
        FOREIGN KEY (source_id)
        REFERENCES data_sources(id)
        ON DELETE SET NULL,

    CONSTRAINT ck_point_latest_values_value_type
        CHECK (value_type IN (1, 2, 3)),

    CONSTRAINT ck_point_latest_values_value
        CHECK (
            (value_type = 1 AND value_num IS NOT NULL AND value_text IS NULL AND value_bool IS NULL)
            OR (value_type = 2 AND value_num IS NULL AND value_text IS NOT NULL AND value_bool IS NULL)
            OR (value_type = 3 AND value_num IS NULL AND value_text IS NULL AND value_bool IS NOT NULL)
        ),

    CONSTRAINT ck_point_latest_values_quality
        CHECK (quality BETWEEN 0 AND 5)
);

CREATE INDEX idx_point_latest_values_source_id
    ON point_latest_values (source_id);

COMMENT ON TABLE point_latest_values IS '测点最新值表，供当前值查询和仪表盘使用';

COMMIT;
