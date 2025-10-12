package product

import (
	"github.com/mytheresa/go-hiring-challenge/internal/category"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductModel struct {
	ID         uint            `gorm:"primaryKey"`
	Code       string          `gorm:"uniqueIndex;not null"`
	Price      decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	Variants   []VariantModel  `gorm:"foreignKey:ProductID"`
	CategoryID *uint
	Category   *category.CategoryModel `gorm:"foreignKey:CategoryID"`
}

func (p *ProductModel) TableName() string {
	return "products"
}

type VariantModel struct {
	ID        uint            `gorm:"primaryKey"`
	ProductID uint            `gorm:"not null"`
	Name      string          `gorm:"not null"`
	SKU       string          `gorm:"uniqueIndex;not null"`
	Price     decimal.Decimal `gorm:"type:decimal(10,2);null"`
}

func (v *VariantModel) TableName() string {
	return "product_variants"
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) Repository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts() ([]Product, error) {
	var productModels []ProductModel
	if err := r.db.Preload("Variants").Preload("Category").Find(&productModels).Error; err != nil {
		return nil, err
	}

	// Map Db model to Domain one
	products := make([]Product, len(productModels))
	for i, pm := range productModels {
		products[i] = mapProductModel(pm)
	}

	return products, nil
}

func mapProductModel(pm ProductModel) Product {
	return Product{
		Code:     pm.Code,
		Price:    pm.Price.InexactFloat64(),
		Variants: mapVariantModels(pm.Variants),
		Category: mapCategoryModel(pm.Category),
	}
}

func mapVariantModels(vms []VariantModel) []Variant {
	if vms == nil {
		return nil
	}

	variants := make([]Variant, len(vms))
	for i, vm := range vms {
		variants[i] = Variant{
			Name:  vm.Name,
			SKU:   vm.SKU,
			Price: vm.Price.InexactFloat64(),
		}
	}

	return variants
}

func mapCategoryModel(c *category.CategoryModel) *category.Category {
	if c == nil {
		return nil
	}

	return &category.Category{
		Code: c.Code,
		Name: c.Name,
	}
}
