# 数据库迁移

表结构以本目录下的 SQL 脚本为准，应用不使用 GORM `AutoMigrate`。

## 执行迁移

按编号顺序执行 `*.up.sql`（每个脚本自带 `BEGIN/COMMIT`）：

```bash
for f in database/0000*.up.sql; do
  psql "postgres://postgres@localhost:5432/gogogo?sslmode=disable" \
    -v ON_ERROR_STOP=1 -f "$f" || exit 1
done
```

回滚时按编号倒序执行 `*.down.sql`。

## 分区

`point_samples` 是按月 RANGE 分区表，迁移只创建父表。应用启动时会自动创建
上个月、当前月以及未来 `database.partition_months_ahead` 个月的分区，
并每天检查一次。因此首次启动前必须先执行迁移。

## 应用配置

非敏感连接参数在 `config/config.yaml` 的 `database` 段配置，密码通过环境变量注入：

| 环境变量 | 说明 |
| --- | --- |
| `DB_HOST` / `DB_PORT` / `DB_NAME` / `DB_USER` | 连接地址、端口、库名、用户 |
| `DB_PASSWORD` | 数据库密码（不要写入配置文件） |
| `DB_SSLMODE` | `disable` / `require` 等 |
| `DB_DSN` | 完整连接串，设置后覆盖以上单项配置 |
| `DB_MAX_OPEN_CONNS` / `DB_MAX_IDLE_CONNS` | 连接池大小 |
| `DB_CONN_MAX_LIFETIME` / `DB_CONN_MAX_IDLE_TIME` | 连接生命周期，如 `1h`、`10m` |
| `DB_PARTITION_MONTHS_AHEAD` | 预建未来分区数量，默认 2 |
