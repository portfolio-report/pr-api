package stats

import (
	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
)

type StatsHandler struct {
	DB             *gorm.DB
	SessionService models.SessionService
	UserService    models.UserService
	GeoipService   models.GeoipService
}

// NewHandler creates new StatsHandler and registers routes
func NewHandler(
	R *gin.RouterGroup,
	DB *gorm.DB,
	UserService models.UserService,
	SessionService models.SessionService,
	GeoipService models.GeoipService,
) {
	h := &StatsHandler{
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
