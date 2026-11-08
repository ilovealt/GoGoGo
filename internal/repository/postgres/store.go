package postgres

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"gogogo/internal/pkg/apperrors"
)

// Store 是 Store 接口的 GORM/PostgreSQL 实现。
// 表结构由 database/ 下的 SQL 迁移维护，不使用 AutoMigrate。
type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

// DB 暴露底层的 GORM 句柄，供启动任务（分区维护、健康检查等）使用。
func (s *Store) DB() *gorm.DB {
	return s.db
}

func now() time.Time {
	return time.Now().UTC()
}

func invalid(message string) error {
	return apperrors.New(apperrors.ErrInvalid, message)
}

func notFound(message string) error {
	return apperrors.New(apperrors.ErrNotFound, message)
}

func conflict(message string) error {
	return apperrors.New(apperrors.ErrConflict, message)
}

// recordNotFound 把 gorm.ErrRecordNotFound 转成带业务文案的 NotFound。
func recordNotFound(err error, message string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound(message)
	}
	return err
}

// duplicateAs 把唯一约束冲突转成 Conflict。
func duplicateAs(err error, message string) error {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return conflict(message)
	}
	return err
}

// foreignKeyAs 把外键冲突转成 Conflict。
func foreignKeyAs(err error, message string) error {
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return conflict(message)
	}
	return err
}

func trim(value string) string {
	return strings.TrimSpace(value)
}

// likePattern 生成大小写不敏感的前缀/包含匹配表达式所需的正则文本，
// 同时转义 LIKE 的通配符，保证关键字按字面量匹配。
func likePattern(keyword string) string {
	escaped := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(strings.ToLower(trim(keyword)))
	return "%" + escaped + "%"
}
