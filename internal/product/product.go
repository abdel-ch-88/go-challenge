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
	Price float64
}

type Repository interface {
	GetAllProductsWith(offset uint, limit uint, categoryID *uint, priceLimit *float64) ([]Product, uint, error)
}
