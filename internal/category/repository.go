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

func (r *CategoryRepository) GetAllCategories() ([]Category, error) {
	var catModels []CategoryModel
	if err := r.db.Find(&catModels).Error; err != nil {
		return nil, err
	}

	if catModels == nil {
		return nil, nil
	}

	cats := make([]Category, len(catModels))
	for i, cm := range catModels {
		cats[i] = *mapCategoryModel(&cm)
	}

	return cats, nil
}

func (r *CategoryRepository) Save(c Category) error {
	catModel := mapNewCategory(&c)

	res := r.db.Create(catModel)
	if res.Error != nil {
		return res.Error
	}

	return nil
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

func mapNewCategory(c *Category) *CategoryModel {
	return &CategoryModel{
		Code: c.Code,
		Name: c.Name,
	}
}
