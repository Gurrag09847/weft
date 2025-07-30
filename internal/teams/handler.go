package teams

import (
	"github.com/Gurrag09847/weft/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)


type Handler interface {
	GetTeamsHandler(c *gin.Context)
}

type handler struct {
	service Service
}

func NewHandler(db *database.Service) Handler {
	return &handler{
		service: NewService(db),
	}
}

func (h *handler) GetTeamsHandler(c *gin.Context) {
	Teams, err := h.service.GetTeamsService()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Teams)
}