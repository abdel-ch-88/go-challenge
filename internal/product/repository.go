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

func (r *ProductsRepository) GetAllProductsWith(offset uint, limit uint, catId *uint, priceLimit *float64) ([]Product, uint, error) {
	query := r.db.Model(&ProductModel{})

	if catId != nil {
		query = query.Where("category_id = ?", *catId)
	}
	if priceLimit != nil {
		query = query.Where("price < ?", *priceLimit)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var productModels []ProductModel
	if err := query.Offset(int(offset)).Limit(int(limit)).
		Preload("Variants").Preload("Category").Find(&productModels).Error; err != nil {
		return nil, 0, err
	}

	// Map Db model to Domain one
	products := make([]Product, len(productModels))
	for i, pm := range productModels {
		products[i] = mapProductModel(pm)
	}

	return products, uint(total), nil
}

func mapProductModel(pm ProductModel) Product {
	return Product{
		ID:       pm.ID,
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
			ID:    vm.ID,
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
