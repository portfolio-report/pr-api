package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/graph/dataloaders"
	"github.com/portfolio-report/pr-api/graph/generated"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/shopspring/decimal"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

func (r *exchangerateResolver) Prices(ctx context.Context, obj *model.Exchangerate, from *string) ([]*model.ExchangeratePrice, error) {
	return r.CurrenciesService.GetExchangeratePrices(obj.ID, from)
}

func (r *mutationResolver) Register(ctx context.Context, username string, password string) (*model.Session, error) {
	user, err := r.UserService.Create(username)
	if err != nil {
		return nil, err
	}
	err = r.UserService.UpdatePassword(ctx, user, password)
	if err != nil {
		panic(err)
	}

	useragent := middleware.UseragentFromContext(ctx)
	session, err := r.SessionService.CreateSession(user, useragent)
	if err != nil {
		panic(err)
	}

	return session, nil
}

func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*model.Session, error) {
	user, err := r.UserService.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gqlerror.Errorf("Unauthorized")
		}

		panic(err)
	}

	valid, err := r.UserService.VerifyPassword(ctx, user, password)
	if err != nil {
		panic(err)
	}
	if !valid {
		return nil, gqlerror.Errorf("Unauthorized")
	}

	useragent := middleware.UseragentFromContext(ctx)

	session, err := r.SessionService.CreateSession(user, useragent)
	if err != nil {
		panic(err)
	}

	return session, nil
}

func (r *mutationResolver) CreateSession(ctx context.Context, note string) (*model.Session, error) {
	user := middleware.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("Access denied")
	}

	session, err := r.SessionService.CreateSession(user, note)
	if err != nil {
		panic(err)
	}

	return session, nil
}

func (r *mutationResolver) DeleteSession(ctx context.Context, token string) (*model.Session, error) {
	user := middleware.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("Access denied")
	}

	session, err := r.SessionService.DeleteSession(token)
	if err != nil {
		panic(err)
	}
	return session, nil
}

func (r *mutationResolver) CreatePortfolio(ctx context.Context, portfolio model.PortfolioInput) (*model.Portfolio, error) {
	user := middleware.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("Access denied")
	}
	return r.PortfolioService.CreatePortfolio(user, &portfolio)
}

func (r *mutationResolver) UpdatePortfolio(ctx context.Context, id int, portfolio model.PortfolioInput) (*model.Portfolio, error) {
	user := middleware.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("Access denied")
	}

	_, err := r.PortfolioService.GetPortfolioOfUserByID(user, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Not found")
		}
		panic(err)
	}

	return r.PortfolioService.UpdatePortfolio(uint(id), &portfolio)
}

func (r *mutationResolver) DeletePortfolio(ctx context.Context, id int) (*model.Portfolio, error) {
	user := middleware.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("Access denied")
	}

	_, err := r.PortfolioService.GetPortfolioOfUserByID(user, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Not found")
		}
		panic(err)
	}

	return r.PortfolioService.DeletePortfolio(uint(id))
}

func (r *portfolioAccountResolver) Value(ctx context.Context, obj *model.PortfolioAccount, currencyCode *string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *portfolioSecurityResolver) Shares(ctx context.Context, obj *model.PortfolioSecurity) (*decimal.Decimal, error) {
	key := model.PortfolioSecurityKey{PortfolioID: obj.PortfolioID, UUID: obj.UUID}
	return dataloaders.For(ctx).PortfolioSecuritySharesByUUID.Load(key)
}

func (r *queryResolver) Currencies(ctx context.Context) ([]*model.Currency, error) {
	return r.CurrenciesService.GetCurrencies()
}

func (r *queryResolver) Exchangerate(ctx context.Context, baseCurrencyCode string, quoteCurrencyCode string) (*model.Exchangerate, error) {
	return r.CurrenciesService.GetExchangerate(baseCurrencyCode, quoteCurrencyCode)
}

func (r *queryResolver) Portfolios(ctx context.Context) ([]*model.Portfolio, error) {
	user := middleware.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("Access denied")
	}

	return r.PortfolioService.GetAllOfUser(user)
}

func (r *queryResolver) Portfolio(ctx context.Context, id int) (*model.Portfolio, error) {
	user := middleware.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("Access denied")
	}

	return r.PortfolioService.GetPortfolioOfUserByID(user, uint(id))
}

func (r *queryResolver) PortfolioAccounts(ctx context.Context, portfolioID int) ([]*model.PortfolioAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PortfolioSecurities(ctx context.Context, portfolioID int) ([]*model.PortfolioSecurity, error) {
	user := middleware.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("Access denied")
	}

	return r.PortfolioService.GetPortfolioSecuritiesOfPortfolio(portfolioID)
}

func (r *queryResolver) PortfolioSecurity(ctx context.Context, portfolioID int, uuid uuid.UUID) (*model.PortfolioSecurity, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Security(ctx context.Context, uuid uuid.UUID) (*model.Security, error) {
	security, err := r.SecurityService.GetSecurityByUUID(uuid)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Not found")
	}
	if err != nil {
		panic(err)
	}
	return security, nil
}

func (r *queryResolver) Sessions(ctx context.Context) ([]*model.Session, error) {
	user := middleware.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("Access denied")
	}

	return r.SessionService.GetAllOfUser(user)
}

func (r *securityResolver) SecurityTaxonomies(ctx context.Context, obj *model.Security) ([]*model.SecurityTaxonomy, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *securityResolver) Events(ctx context.Context, obj *model.Security) ([]*model.Event, error) {
	return r.SecurityService.GetEventsOfSecurity(obj)
}

func (r *securityTaxonomyResolver) Taxonomy(ctx context.Context, obj *model.SecurityTaxonomy) (*model.Taxonomy, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *sessionResolver) User(ctx context.Context, obj *model.Session) (*model.User, error) {
	return dataloaders.For(ctx).UserByID.Load(int(obj.UserID))
}

// Exchangerate returns generated.ExchangerateResolver implementation.
func (r *Resolver) Exchangerate() generated.ExchangerateResolver { return &exchangerateResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// PortfolioAccount returns generated.PortfolioAccountResolver implementation.
func (r *Resolver) PortfolioAccount() generated.PortfolioAccountResolver {
	return &portfolioAccountResolver{r}
}

// PortfolioSecurity returns generated.PortfolioSecurityResolver implementation.
func (r *Resolver) PortfolioSecurity() generated.PortfolioSecurityResolver {
	return &portfolioSecurityResolver{r}
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Security returns generated.SecurityResolver implementation.
func (r *Resolver) Security() generated.SecurityResolver { return &securityResolver{r} }

// SecurityTaxonomy returns generated.SecurityTaxonomyResolver implementation.
func (r *Resolver) SecurityTaxonomy() generated.SecurityTaxonomyResolver {
	return &securityTaxonomyResolver{r}
}

// Session returns generated.SessionResolver implementation.
func (r *Resolver) Session() generated.SessionResolver { return &sessionResolver{r} }

type exchangerateResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type portfolioAccountResolver struct{ *Resolver }
type portfolioSecurityResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type securityResolver struct{ *Resolver }
type securityTaxonomyResolver struct{ *Resolver }
type sessionResolver struct{ *Resolver }
