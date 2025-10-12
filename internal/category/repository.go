package category

type CategoryModel struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"uniqueIndex;not null"`
	Name string `gorm:"not null"`
}

func (p *CategoryModel) TableName() string {
	return "categories"
}
