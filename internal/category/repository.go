package category

import (
	"gorm.io/gorm"
)

type CategoryModel struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"uniqueIndex;not null"`
	Name string `gorm:"not null"`
}

func (p *CategoryModel) TableName() string {
	return "categories"
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) Repository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) GetCategoryByCode(code string) (*Category, error) {
	var cat *CategoryModel
	if err := r.db.Where("code = ?", code).Find(&cat).Error; err != nil {
		return nil, err
	}

	if cat != nil {
		return mapCategoryModel(cat), nil
	}

	return nil, nil
}

func mapCategoryModel(c *CategoryModel) *Category {
	if c == nil {
		return nil
	}

	return &Category{
		ID:   c.ID,
		Code: c.Code,
		Name: c.Name,
	}
}
