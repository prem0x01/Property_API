package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	TotalItems int64       `json:"total_items"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status: "success",
		Data:   data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Status: "error",
		Error:  message,
	})
}

func PaginatedSuccessResponse(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	totalPages := (int(total) + pageSize - 1) / pageSize

	c.JSON(http.StatusOK, Response{
		Status: "success",
		Data: PaginatedResponse{
			Items:      items,
			TotalItems: total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	})
}
