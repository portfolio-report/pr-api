package stats

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/libs"
)

func (h *StatsHandler) DeleteClientupdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	err = h.DB.Delete(&db.Clientupdate{}, id).Error
	if err != nil {
		panic(err)
	}

	c.Status(http.StatusNoContent)
}
