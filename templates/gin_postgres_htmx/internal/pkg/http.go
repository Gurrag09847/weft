package pkg

import "github.com/gin-gonic/gin"

type Response struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrResponse struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func HttpResponse(c *gin.Context, r Response) {
	c.JSON(r.Status, r)
}

func ErrorResponse(c *gin.Context, r ErrResponse) {
	c.JSON(r.Status, r)
}
