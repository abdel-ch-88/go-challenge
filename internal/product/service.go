package product

import (
	"errors"

	"github.com/mytheresa/go-hiring-challenge/internal/category"
)

const DefaultOffset = 0
const DefaultLimit = 10

type Service struct {
	prodRepo Repository
	catRepo  category.Repository
}

type FindCriteria struct {
	Offset        uint
	Limit         uint
	CategoryID    *uint
	PriceLessThan *float64
}

func NewService(prodRepo Repository, catRepo category.Repository) *Service {
	return &Service{
		prodRepo: prodRepo,
		catRepo:  catRepo,
	}
}

func (srv *Service) BuildFindCriteria(offset *uint, limit *uint, catCode *string, priceLessThan *float64) (*FindCriteria, error) {
	criteria := FindCriteria{}

	if offset == nil {
		criteria.Offset = DefaultOffset
	} else {
		criteria.Offset = *offset
	}

	if limit == nil {
		criteria.Limit = DefaultLimit
	} else {
		if *limit < 1 || *limit > 100 {
			return nil, errors.New("limit should be between 1 and 100")
		}
		criteria.Limit = *limit
	}

	if catCode != nil {
		cat, err := srv.catRepo.GetCategoryByCode(*catCode)
		if err != nil {
			return nil, err
		}
		if cat != nil {
			criteria.CategoryID = &cat.ID
		}
	}

	if priceLessThan != nil && *priceLessThan < float64(0.01) {
		return nil, errors.New("price limit cannot be lower than 0.01")
	}
	criteria.PriceLessThan = priceLessThan

	return &criteria, nil
}

func (srv *Service) FindProducts(fc FindCriteria) ([]Product, uint, error) {
	return srv.prodRepo.GetAllProductsWith(fc.Offset, fc.Limit, fc.CategoryID, fc.PriceLessThan)
}
