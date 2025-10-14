package category

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MOCK the repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetAllCategories() ([]Category, error) {
	args := m.Called()

	return args.Get(0).([]Category), args.Error(1)
}

func (m *MockRepository) GetCategoryByCode(newCat string) (*Category, error) {
	args := m.Called(newCat)

	cat := args.Get(0)
	if cat != nil {
		return cat.(*Category), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockRepository) Save(c Category) error {
	args := m.Called(c)

	return args.Error(0)
}

func TestGetAllCategories(t *testing.T) {
	mockRepo := new(MockRepository)

	serv := NewService(mockRepo)

	// Given
	cats := []Category{
		{Code: "P+", Name: "plus"},
		{Code: "*", Name: "premium"},
	}

	// Stage Call
	mockRepo.On("GetAllCategories").Return(cats, nil)

	// When
	res, err := serv.GetAllCategories()

	// Then
	assert.Nil(t, err)
	assert.Equal(t, cats, res)
}

func TestCreateCategory(t *testing.T) {

	t.Run("Happy Path", func(t *testing.T) {
		mockRepo := new(MockRepository)

		serv := NewService(mockRepo)

		// Given
		newCat := Category{Code: "C123", Name: "metal"}

		// Stage Call
		mockRepo.On("GetCategoryByCode", newCat.Code).Return(nil, nil)
		mockRepo.On("Save", newCat).Return(nil)

		// When
		err := serv.CreateCategory(newCat)

		// Then
		assert.Nil(t, err)
	})

	t.Run("Invalid category", func(t *testing.T) {
		type testCase struct {
			Name     string
			Category Category
		}

		invalidCats := []testCase{
			{"Invalid category", Category{Code: "", Name: ""}},
			{"Invalid Name", Category{Code: "C012", Name: ""}},
			{"Invalid Code", Category{Code: "", Name: "High"}},
		}

		mockRepo := new(MockRepository)
		serv := NewService(mockRepo)

		for _, tCase := range invalidCats {
			t.Run(tCase.Name, func(t *testing.T) {
				// When
				err := serv.CreateCategory(tCase.Category)

				// Then
				assert.ErrorContains(t, err, "invalid Category data")
			})
		}
	})

	t.Run("Failed to cehck existing Category", func(t *testing.T) {
		mockRepo := new(MockRepository)

		serv := NewService(mockRepo)

		// Given
		newCat := Category{Code: "C666", Name: "down"}

		// Stage Call
		mockRepo.On("GetCategoryByCode", newCat.Code).Return(nil, errors.New("disconnected"))

		// When
		err := serv.CreateCategory(newCat)

		// Then
		assert.ErrorContains(t, err, "disconnected")
	})

	t.Run("Category already exists", func(t *testing.T) {
		mockRepo := new(MockRepository)

		serv := NewService(mockRepo)

		// Given
		newCat := Category{Code: "C000", Name: "zero"}

		// Stage Call
		mockRepo.On("GetCategoryByCode", newCat.Code).Return(&Category{ID: 1, Code: "C000"}, nil)

		// When
		err := serv.CreateCategory(newCat)

		// Then
		assert.ErrorContains(t, err, "ategody code already exists")
	})
}
