package portfolios

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
)

type putAccountRequest struct {
	Type                 string    `json:"type" binding:"oneof=deposit securities"`
	Name                 string    `json:"name" binding:"required"`
	CurrencyCode         *string   `json:"currencyCode" binding:"required_if=Type deposit,max=3"`
	ReferenceAccountUuid *string   `json:"referenceAccountUuid" binding:"omitempty,uuid"`
	Active               bool      `json:"active"`
	Note                 string    `json:"note"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

// PutAccount creates or updates account in portfolio
func (h *portfoliosHandler) PutAccount(c *gin.Context) {
	portfolioId := middleware.PortfolioFromContext(c).ID
	uuid := c.Param("uuid")
	if err := h.Validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	var req putAccountRequest
	if err := c.BindJSON(&req); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	var account db.PortfolioAccount

	err := h.DB.FirstOrInit(&account, db.PortfolioAccount{PortfolioID: uint(portfolioId), UUID: uuid}).Error
	if err != nil {
		panic(err)
	}

	account.Type = req.Type
	account.Name = req.Name
	account.Active = req.Active
	account.Note = req.Note
	account.UpdatedAt = req.UpdatedAt

	switch req.Type {
	case "deposit":
		account.CurrencyCode = req.CurrencyCode
		account.ReferenceAccountUUID = nil
	case "securities":
		account.CurrencyCode = nil
		if req.ReferenceAccountUuid == nil {
			libs.HandleBadRequestError(c, "referenceAccountUuid is missing")
			return
		}
		account.ReferenceAccountUUID = req.ReferenceAccountUuid
	default:
		panic(fmt.Errorf("invalid type: %s", req.Type))
	}

	if err := h.DB.Save(&account).Error; err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23503" {
			libs.HandleBadRequestError(c, "data violates constraint "+pgErr.Constraint)
			return
		}

		panic(err)
	}

	c.JSON(http.StatusOK, models.PortfolioAccountResponseFromDB(&account))
}
