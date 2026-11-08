package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeOK           = 0
	CodeInvalid      = 1001
	CodeUnauthorized = 1002
	CodeForbidden    = 1003
	CodeNotFound     = 1004
	CodeConflict     = 1005
	CodeInternal     = 1006
)

type Envelope struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func requestID(c *gin.Context) string {
	if value, ok := c.Get("request_id"); ok {
		if id, ok := value.(string); ok {
			return id
		}
	}
	return ""
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{
		Code:      CodeOK,
		Message:   "ok",
		Data:      data,
		RequestID: requestID(c),
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Envelope{
		Code:      CodeOK,
		Message:   "created",
		Data:      data,
		RequestID: requestID(c),
	})
}

func Error(c *gin.Context, status, code int, message string) {
	c.AbortWithStatusJSON(status, Envelope{
		Code:      code,
		Message:   message,
		RequestID: requestID(c),
	})
}

type Page struct {
	Items    any `json:"items"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

func NewPage(items any, total, page, pageSize int) Page {
	return Page{Items: items, Page: page, PageSize: pageSize, Total: total}
}

// Paginate 对内存切片做分页，返回当前页数据和总数。
func Paginate[T any](items []T, page, pageSize int) ([]T, int) {
	total := len(items)
	start := (page - 1) * pageSize
	if start >= total {
		return []T{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return items[start:end], total
}
