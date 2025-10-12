package category

// Category represents a category of a product
// Multiple products can be under the same category
type Category struct {
	ID   uint
	Code string
	Name string
}

type Repository interface {
	GetCategoryByCode(code string) (*Category, error)
}
