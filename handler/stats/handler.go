package stats

import (
	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"gorm.io/gorm"
)

type statsHandler struct {
	*gorm.DB
	model.SessionService
	model.UserService
	model.GeoipService
}

// NewHandler creates new stats handler and registers routes
func NewHandler(
	R *gin.RouterGroup,
	DB *gorm.DB,
	UserService model.UserService,
	SessionService model.SessionService,
	GeoipService model.GeoipService,
) {
	h := &statsHandler{
		DB:             DB,
		SessionService: SessionService,
		UserService:    UserService,
		GeoipService:   GeoipService,
	}

	g := R.Group("/stats")

	g.HEAD("/update/name.abuchen.portfolio/:version", h.CountClientupdate)
	g.GET("/updates", h.GetClientupdatesStats)
	g.GET("/updates/:version", h.GetClientupdatesStatsVersion)
	g.GET("/",
		middleware.RequireUser(h.SessionService, h.UserService),
		middleware.RequireAdmin(),
		h.GetClientupdates)
	g.DELETE("/:id",
		middleware.RequireUser(h.SessionService, h.UserService),
		middleware.RequireAdmin(),
		h.DeleteClientupdate)
}
