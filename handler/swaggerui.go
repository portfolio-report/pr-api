package handler

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/swaggerui"
)

//go:embed openapi.json
var openApiSpec []byte

func (h *Handler) RegisterSwaggerUi(g *gin.RouterGroup, prefix string) {
	g.GET(prefix+"/*w", gin.WrapH(http.StripPrefix(prefix, swaggerui.Handler(openApiSpec))))
}
