package stats

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/libs"
)

// GetClientupdates lists updates
func (h *statsHandler) GetClientupdates(c *gin.Context) {
	type Query struct {
		Limit      int    `form:"limit"`
		Skip       int    `form:"skip"`
		Sort       string `form:"sort" binding:"omitempty,oneof=id timestamp version country useragent"`
		Descending bool   `form:"descending"`
		Version    string `form:"version"`
	}
	var q Query
	err := c.BindQuery(&q)
	if err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	if q.Limit == 0 {
		q.Limit = 10
	}

	if q.Sort == "" {
		q.Sort = "timestamp"
	}

	query := h.DB

	if q.Version != "" {
		query = query.Where("version = ?", q.Version)
	}

	var totalCount int64
	if err := query.Table("clientupdates").Count(&totalCount).Error; err != nil {
		panic(err)
	}

	order := q.Sort
	if q.Descending {
		order += " desc"
	}
	query = query.Order(order)

	query = query.Limit(q.Limit).Offset(q.Skip)

	var clientupdates []db.Clientupdate
	if err := query.Find(&clientupdates).Error; err != nil {
		panic(err)
	}

	entries := []gin.H{}
	for _, e := range clientupdates {
		entries = append(entries, gin.H{
			"id":        e.ID,
			"timestamp": e.Timestamp,
			"version":   e.Version,
			"country":   e.Country,
			"useragent": e.Useragent,
		})
	}

	c.JSON(http.StatusOK, gin.H{"entries": entries, "params": gin.H{"totalCount": totalCount}})
}
