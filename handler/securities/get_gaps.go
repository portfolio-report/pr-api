package securities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/libs"
)

// GetGaps lists gaps in security prices
func (h *securitiesHandler) GetGaps(c *gin.Context) {
	type Query struct {
		MinDuration int `form:"minDuration"`
		MaxResults  int `form:"maxResults"`
	}

	var q Query
	if err := c.ShouldBindQuery(&q); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	if q.MinDuration == 0 {
		q.MinDuration = 3
	}

	if q.MaxResults == 0 {
		q.MaxResults = 10
	}

	gaps := h.SecurityService.FindGapsInPrices(q.MinDuration, q.MaxResults)

	c.JSON(http.StatusOK, gaps)
}
