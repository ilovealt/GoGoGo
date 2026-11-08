package bootstrap

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// EnsureSamplePartitions 为 point_samples 预建上个月、当前月及未来 monthsAhead 个月的分区。
// 上个月用于容纳月初启动时的回填/演示数据和跨月迟到的样本。
// 可重复执行：已存在的分区会被跳过。
func EnsureSamplePartitions(ctx context.Context, db *gorm.DB, monthsAhead int) error {
	if monthsAhead < 0 {
		monthsAhead = 0
	}

	month := time.Now().UTC()
	month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)

	startMonth := month.AddDate(0, -1, 0)
	for offset := 0; offset <= monthsAhead+1; offset++ {
		from := startMonth.AddDate(0, offset, 0)
		to := from.AddDate(0, 1, 0)
		table := "point_samples_" + from.Format("2006_01")
		statement := fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s PARTITION OF point_samples FOR VALUES FROM ('%s') TO ('%s')",
			table,
			from.Format(time.RFC3339),
			to.Format(time.RFC3339),
		)
		if err := db.WithContext(ctx).Exec(statement).Error; err != nil {
			return fmt.Errorf("create partition %s: %w", table, err)
		}
	}
	return nil
}
