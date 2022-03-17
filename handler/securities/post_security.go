package securities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
)

type CreateUpdateSecurityRequest struct {
	Name         *string `json:"name"`
	Isin         *string `json:"isin"`
	Wkn          *string `json:"wkn"`
	SecurityType *string `json:"securityType"`
	SymbolXfra   *string `json:"symbolXfra"`
	SymbolXnas   *string `json:"symbolXnas"`
	SymbolXnys   *string `json:"symbolXnys"`
}

// PostSecurity creates new security
func (h *SecuritiesHandler) PostSecurity(c *gin.Context) {
	var request CreateUpdateSecurityRequest
	if err := c.BindJSON(&request); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	security := db.Security{
		UUID:         uuid.New().String(),
		Name:         request.Name,
		Isin:         request.Isin,
		Wkn:          request.Wkn,
		SecurityType: request.SecurityType,
		SymbolXfra:   request.SymbolXfra,
		SymbolXnas:   request.SymbolXnas,
		SymbolXnys:   request.SymbolXnys,
	}

	if err := h.DB.Create(&security).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, models.SecurityResponseFromDB(&security))
}
