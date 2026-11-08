BEGIN;

CREATE TABLE point_samples (
    point_id BIGINT NOT NULL,

    ts TIMESTAMPTZ NOT NULL,

    value_type SMALLINT NOT NULL DEFAULT 1,

    value_num DOUBLE PRECISION,
    value_text TEXT,
    value_bool BOOLEAN,

    quality SMALLINT NOT NULL DEFAULT 0,

    source_id BIGINT,

    collected_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),

    CONSTRAINT pk_point_samples
        PRIMARY KEY (point_id, ts),

    CONSTRAINT fk_point_samples_point
        FOREIGN KEY (point_id)
        REFERENCES measurement_points(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_point_samples_source
        FOREIGN KEY (source_id)
        REFERENCES data_sources(id)
        ON DELETE SET NULL,

    CONSTRAINT ck_point_samples_value_type
        CHECK (value_type IN (1, 2, 3)),

    CONSTRAINT ck_point_samples_value
        CHECK (
            (
                value_type = 1
                AND value_num IS NOT NULL
                AND value_text IS NULL
                AND value_bool IS NULL
            )
            OR
            (
                value_type = 2
                AND value_num IS NULL
                AND value_text IS NOT NULL
                AND value_bool IS NULL
            )
            OR
            (
                value_type = 3
                AND value_num IS NULL
                AND value_text IS NULL
                AND value_bool IS NOT NULL
            )
        ),

    CONSTRAINT ck_point_samples_quality
        CHECK (quality BETWEEN 0 AND 5)
) PARTITION BY RANGE (ts);

-- 具体月份分区由应用启动任务（EnsureSamplePartitions）和第 000010 号迁移之外的
-- 日常维护任务创建，避免在迁移脚本中硬编码未来日期。

-- 单测点时间范围查询使用主键 (point_id, ts)
-- 跨测点按时间范围查询使用 BRIN，索引体积远小于 B-tree
CREATE INDEX idx_point_samples_ts_brin
    ON point_samples USING BRIN (ts);

-- 数据源删除时把 source_id 置空，外键检查需要该索引
CREATE INDEX idx_point_samples_source_id
    ON point_samples (source_id);

COMMENT ON TABLE point_samples IS '测点历史时序数据表（按月分区）';
COMMENT ON COLUMN point_samples.point_id IS '测点ID';
COMMENT ON COLUMN point_samples.ts IS '数据发生时间';
COMMENT ON COLUMN point_samples.value_type IS '数据类型：1数值，2字符串，3布尔';
COMMENT ON COLUMN point_samples.value_num IS '数值型数据';
COMMENT ON COLUMN point_samples.value_text IS '字符串型数据';
COMMENT ON COLUMN point_samples.value_bool IS '布尔型数据';
COMMENT ON COLUMN point_samples.quality IS '质量：0正常，1不确定，2坏质量，3通讯中断，4未连接，5超量程';
COMMENT ON COLUMN point_samples.source_id IS '数据来源ID，数据源删除后置空';
COMMENT ON COLUMN point_samples.collected_at IS '系统接收时间，使用clock_timestamp()避免批量插入取到同一事务时间';

COMMIT;
