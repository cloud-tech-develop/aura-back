package sync

import (
	"time"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(database *db.DB) *Handler {
	q := database.Wrap(database.DB)
	return &Handler{svc: NewService(q)}
}

// Pull - GET /sync/pull
func (h *Handler) Pull(c *gin.Context) {
	lastSyncStr := c.Query("last_sync")
	lastSync := time.Time{}
	if lastSyncStr != "" {
		if t, err := time.Parse(time.RFC3339, lastSyncStr); err == nil {
			lastSync = t
		}
	}

	batch, err := h.svc.Pull(c.Request.Context(), lastSync)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, batch)
}

// Push - POST /sync/push
func (h *Handler) Push(c *gin.Context) {
	var batch SyncBatch
	if err := c.ShouldBindJSON(&batch); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	stats, err := h.svc.Push(c.Request.Context(), &batch)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, stats)
}
