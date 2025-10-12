package product

// Product represents a product in the catalog.
// It includes a unique code and a price.
type Product struct {
	Code     string
	Price    float64
	Variants []Variant
}

// Variant represents a product variant in the catalog.
// It includes a unique name, SKU, and an optional price.
// Variants can be used to represent different configurations or options for a product.
type Variant struct {
	Name  string
	SKU   string
	Price float64
}

type Repository interface {
	GetAllProducts() ([]Product, error)
}
