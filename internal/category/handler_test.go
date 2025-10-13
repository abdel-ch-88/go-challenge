package category

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MOCK the categpry Service
type MockService struct {
	mock.Mock
}

func (m *MockService) GetAllCategories() ([]Category, error) {
	args := m.Called()

	cats := args.Get(0)
	if cats != nil {
		return cats.([]Category), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockService) CreateCategory(newCat Category) error {
	args := m.Called(newCat)

	return args.Error(0)
}

func TestHandleList(t *testing.T) {

	t.Run("Test Happy Path", func(t *testing.T) {
		mockService := new(MockService)
		testMux := http.NewServeMux()

		// Stage Service Call
		gienCats := []Category{
			{Code: "C001", Name: "fist"},
			{Code: "C002", Name: "second"},
			{Code: "C003", Name: "third"},
		}

		mockService.On("GetAllCategories").Return(gienCats, nil)

		// Set routing
		handler := NewCategoryHandler(mockService)
		testMux.HandleFunc("GET /categories", handler.HandleList)

		// When
		req := httptest.NewRequest("GET", "/categories", nil)
		record := httptest.NewRecorder()
		testMux.ServeHTTP(record, req)

		// Then
		var actualRes CategoryListResponse
		err := json.Unmarshal(record.Body.Bytes(), &actualRes)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, record.Code)

		expectedRes := CategoryListResponse{
			Categories: []CategoryItem{
				{Code: "C001", Name: "fist"},
				{Code: "C002", Name: "second"},
				{Code: "C003", Name: "third"},
			},
		}
		assert.Equal(t, expectedRes, actualRes)
	})

	t.Run("Test Internal Failure", func(t *testing.T) {
		mockService := new(MockService)
		testMux := http.NewServeMux()

		// Stage Service Call
		mockService.On("GetAllCategories").Return(nil, errors.New("categorical error"))

		// Set routing
		handler := NewCategoryHandler(mockService)
		testMux.HandleFunc("GET /categories", handler.HandleList)

		// When
		req := httptest.NewRequest("GET", "/categories", nil)
		record := httptest.NewRecorder()
		testMux.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusInternalServerError, record.Code)
	})
}

func TestHandleCreate(t *testing.T) {

	t.Run("Test Happy Path", func(t *testing.T) {
		mockService := new(MockService)
		testMux := http.NewServeMux()

		// Given
		newCat := CategoryItem{Code: "000", Name: "VIP"}

		newCatBody, err := json.Marshal(newCat)
		assert.NoError(t, err)

		// Stage Service Call
		cat := Category{Code: newCat.Code, Name: newCat.Name}
		mockService.On("CreateCategory", cat).Return(nil)

		// Set routing
		handler := NewCategoryHandler(mockService)
		testMux.HandleFunc("POST /categories", handler.HandleCreate)

		// When
		reader := bytes.NewReader(newCatBody)
		req := httptest.NewRequest("POST", "/categories", reader)
		record := httptest.NewRecorder()
		testMux.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusNoContent, record.Code)
	})

	t.Run("Test Invalid catalog", func(t *testing.T) {
		mockService := new(MockService)
		testMux := http.NewServeMux()

		// Given
		badBody := `[{}]`

		// Set routing
		handler := NewCategoryHandler(mockService)
		testMux.HandleFunc("POST /categories", handler.HandleCreate)

		// When
		reader := strings.NewReader(badBody)
		req := httptest.NewRequest("POST", "/categories", reader)
		record := httptest.NewRecorder()
		testMux.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusBadRequest, record.Code)
	})

	t.Run("Test failed creation", func(t *testing.T) {
		mockService := new(MockService)
		testMux := http.NewServeMux()

		// Given
		catBody := `{"Code": "A", "Name": "fiber"}`

		// Stage Service Call
		mockService.On("CreateCategory", mock.Anything).Return(errors.New("something flew"))

		// Set routing
		handler := NewCategoryHandler(mockService)
		testMux.HandleFunc("POST /categories", handler.HandleCreate)

		// When
		reader := strings.NewReader(catBody)
		req := httptest.NewRequest("POST", "/categories", reader)
		record := httptest.NewRecorder()
		testMux.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusInternalServerError, record.Code)
	})
}
