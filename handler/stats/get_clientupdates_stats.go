package stats

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GetClientupdatesResponse struct {
	Version     string    `json:"version"`
	Count       uint      `json:"count"`
	FirstUpdate time.Time `json:"firstUpdate"`
	LastUpdate  time.Time `json:"lastUpdate"`
}

func (h *StatsHandler) GetClientupdatesStats(c *gin.Context) {
	var results []GetClientupdatesResponse

	err := h.DB.Raw(`
		SELECT
			version,
			count(*) AS count,
			min(timestamp) as "first_update",
			max(timestamp) as "last_update"
		FROM clientupdates
		GROUP BY version
		`).Scan(&results).Error

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, results)
}
