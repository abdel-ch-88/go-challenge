package product

import "github.com/mytheresa/go-hiring-challenge/internal/category"

// Product represents a product in the catalog.
// It includes a unique code and a price.
// It can have multiple variants and belongs only in one category
type Product struct {
	ID       uint
	Code     string
	Price    float64
	Variants []Variant
	Category *category.Category
}

// Variant represents a product variant in the catalog.
// It includes a unique name, SKU, and an optional price.
// Variants can be used to represent different configurations or options for a product.
type Variant struct {
	ID    uint
	Name  string
	SKU   string
	Price *float64
}

// A Service that will be implemented with the business logic of the doamin product
type FindCriteria struct {
	Offset        uint
	Limit         uint
	CategoryID    *uint
	PriceLessThan *float64
}

type Service interface {
	BuildFindCriteria(offset *uint, limit *uint, catCode *string, priceLessThan *float64) (*FindCriteria, error)
	FindProducts(fc FindCriteria) ([]Product, uint, error)
	GetProduct(code string) (*Product, error)
}

type Repository interface {
	GetAllProductsWith(offset uint, limit uint, categoryID *uint, priceLimit *float64) ([]Product, uint, error)
	GetProductByCode(code string) (*Product, error)
}
