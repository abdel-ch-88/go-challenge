package product

import (
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/category"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MOCK the category repository
type MockCatRepository struct {
	mock.Mock
}

func (m *MockCatRepository) GetAllCategories() ([]category.Category, error) {
	args := m.Called()

	return args.Get(0).([]category.Category), args.Error(1)
}

func (m *MockCatRepository) GetCategoryByCode(newCat string) (*category.Category, error) {
	args := m.Called(newCat)

	cat := args.Get(0)
	if cat != nil {
		return cat.(*category.Category), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockCatRepository) Save(c category.Category) error {
	args := m.Called(c)

	return args.Error(0)
}

// MOCK the repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetAllProductsWith(offset uint, limit uint, categoryID *uint, priceLimit *float64) ([]Product, uint, error) {
	args := m.Called(offset, limit, categoryID, priceLimit)

	return args.Get(0).([]Product), args.Get(1).(uint), args.Error(2)
}

func (m *MockRepository) GetProductByCode(code string) (*Product, error) {
	args := m.Called(code)

	prod := args.Get(0)
	if prod != nil {
		return prod.(*Product), args.Error(1)
	}

	return nil, args.Error(1)
}

func TestBuildFindCriteria(t *testing.T) {

	t.Run("Happy Path all params", func(t *testing.T) {
		mockCatRepo := new(MockCatRepository)

		// Given
		offset := uint(5)
		limit := uint(20)
		catCode := "C07"
		catId := uint(7)
		priceLessThan := float64(4.99)

		// Stage catalog Repository call
		mockCatRepo.On("GetCategoryByCode", catCode).Return(&category.Category{ID: catId}, nil)

		serv := NewService(nil, mockCatRepo)

		// When
		fc, err := serv.BuildFindCriteria(&offset, &limit, &catCode, &priceLessThan)

		// Then
		expected := &FindCriteria{
			Offset:        offset,
			Limit:         limit,
			CategoryID:    &catId,
			PriceLessThan: &priceLessThan,
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, fc)
	})

	t.Run("Default with all params nil", func(t *testing.T) {
		serv := NewService(nil, nil)

		// When
		fc, err := serv.BuildFindCriteria(nil, nil, nil, nil)

		// Then
		expected := &FindCriteria{
			Offset:        DefaultOffset,
			Limit:         DefaultLimit,
			CategoryID:    nil,
			PriceLessThan: nil,
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, fc)
	})

	t.Run("Invalid limits", func(t *testing.T) {
		serv := NewService(nil, nil)

		// Given invalid limits
		limit0 := uint(0)
		limit1 := uint(101)

		// When
		fc0, err0 := serv.BuildFindCriteria(nil, &limit0, nil, nil)
		fc1, err1 := serv.BuildFindCriteria(nil, &limit1, nil, nil)

		// Then
		assert.Nil(t, fc0)
		assert.Nil(t, fc1)
		assert.ErrorContains(t, err0, "limit should be between 1 and 100")
		assert.ErrorContains(t, err1, "limit should be between 1 and 100")
	})

	t.Run("Invlaid category", func(t *testing.T) {
		mockCatRepo := new(MockCatRepository)

		serv := NewService(nil, mockCatRepo)

		// Given
		catCode := "inj"

		// Stage catalog repository call
		mockCatRepo.On("GetCategoryByCode", catCode).Return(nil, errors.New("migrated"))

		// When
		fc, err := serv.BuildFindCriteria(nil, nil, &catCode, nil)

		// Then
		assert.Nil(t, fc)
		assert.ErrorContains(t, err, "migrated")
	})

	t.Run("Invalid priceLessThan", func(t *testing.T) {
		serv := NewService(nil, nil)

		// Given invalid price
		priceLessThan := float64(0.001)

		// When
		fc, err := serv.BuildFindCriteria(nil, nil, nil, &priceLessThan)

		// Then
		assert.Nil(t, fc)
		assert.ErrorContains(t, err, "price limit cannot be lower than 0.01")
	})
}

func TestFindProducts(t *testing.T) {
	mockRepo := new(MockRepository)

	serv := NewService(mockRepo, nil)

	// Given
	findCrit := FindCriteria{Offset: 0, Limit: 10, CategoryID: nil, PriceLessThan: nil}

	total := uint(33)
	price := float64(9.99)
	vars := []Variant{
		{Name: "new one", SKU: "MP5543", Price: &price},
		{Name: "very old", SKU: "MP5542", Price: &price},
	}
	cat := category.Category{Code: "C00X", Name: "exclusive"}
	prods := []Product{
		{Code: "P00A", Price: 4.5, Variants: vars, Category: &cat},
		{Code: "P00B", Price: 10.0, Variants: vars, Category: &cat},
	}

	// Stage Call
	mockRepo.On("GetAllProductsWith", findCrit.Offset, findCrit.Limit, findCrit.CategoryID, findCrit.PriceLessThan).
		Return(prods, total, nil)

	// When
	res, tot, err := serv.FindProducts(findCrit)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, total, tot)
	assert.Equal(t, prods, res)
}

func TestGetProduct(t *testing.T) {

	t.Run("Happy Path", func(t *testing.T) {
		mockRepo := new(MockRepository)

		serv := NewService(mockRepo, nil)

		// Given
		prodCode := "PR002"

		price := float64(5.01)
		vars := []Variant{
			{Name: "left", SKU: "MP5543", Price: &price},
			{Name: "center", SKU: "MP5542", Price: nil},
			{Name: "right", SKU: "MP5542", Price: nil},
		}
		cat := category.Category{Code: "C00X", Name: "exclusive"}
		prod := Product{Code: prodCode, Price: 6.51, Variants: vars, Category: &cat}

		mockRepo.On("GetProductByCode", prodCode).Return(&prod, nil)

		// When
		res, err := serv.GetProduct(prodCode)

		// Then
		expectedVars := []Variant{
			{Name: "left", SKU: "MP5543", Price: &price},
			{Name: "center", SKU: "MP5542", Price: &prod.Price},
			{Name: "right", SKU: "MP5542", Price: &prod.Price},
		}
		expectedProd := &Product{Code: prodCode, Price: prod.Price, Variants: expectedVars, Category: &cat}
		assert.Nil(t, err)
		assert.Equal(t, expectedProd, res)
	})

	t.Run("Failed to get Product", func(t *testing.T) {
		mockRepo := new(MockRepository)

		serv := NewService(mockRepo, nil)

		// Given
		prodCode := "PR404"

		mockRepo.On("GetProductByCode", prodCode).Return(nil, errors.New("missing link"))

		// When
		res, err := serv.GetProduct(prodCode)

		// Then
		assert.Nil(t, res)
		assert.ErrorContains(t, err, "missing link")
	})

	t.Run("No Prodcuct found nil", func(t *testing.T) {
		mockRepo := new(MockRepository)

		serv := NewService(mockRepo, nil)

		// Given
		prodCode := "PR--8"

		mockRepo.On("GetProductByCode", prodCode).Return(nil, nil)

		// When
		res, err := serv.GetProduct(prodCode)

		// Then
		assert.Nil(t, res)
		assert.Nil(t, err)
	})

	t.Run("No Prodcuct found empty", func(t *testing.T) {
		mockRepo := new(MockRepository)

		serv := NewService(mockRepo, nil)

		// Given
		prodCode := "PR--9"

		mockRepo.On("GetProductByCode", prodCode).Return(&Product{}, nil)

		// When
		res, err := serv.GetProduct(prodCode)

		// Then
		assert.Nil(t, res)
		assert.Nil(t, err)
	})
}
