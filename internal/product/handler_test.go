package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/category"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MOCK the product Service
type MockService struct {
	mock.Mock
}

func (m *MockService) BuildFindCriteria(offset *uint, limit *uint, catCode *string, priceLessThan *float64) (*FindCriteria, error) {
	args := m.Called(offset, limit, catCode, priceLessThan)

	fc := args.Get(0)
	if fc != nil {
		return fc.(*FindCriteria), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockService) FindProducts(fc FindCriteria) ([]Product, uint, error) {
	args := m.Called(fc)

	prods := args.Get(0)
	if prods != nil {
		return prods.([]Product), args.Get(1).(uint), args.Error(2)
	}

	return nil, args.Get(1).(uint), args.Error(2)
}

func (m *MockService) GetProduct(code string) (*Product, error) {
	args := m.Called(code)

	prod := args.Get(0)
	if prod != nil {
		return prod.(*Product), args.Error(1)
	}

	return nil, args.Error(1)
}

func TestHandleList(t *testing.T) {

	t.Run("Test Happy path", func(t *testing.T) {
		mockService := new(MockService)
		testMux := http.NewServeMux()

		// Given GET params
		offset := uint(20)
		limit := uint(10)
		cat := "CAT"
		catId := uint(1)
		priceLimit := float64(9.99)

		// Stage Service Calls
		findCrt := FindCriteria{Offset: offset, Limit: limit, CategoryID: &catId, PriceLessThan: &priceLimit}

		mockService.On("BuildFindCriteria", &offset, &limit, &cat, &priceLimit).
			Return(&findCrt, nil)

		total := uint(13)
		givenCat := category.Category{Code: "CAT", Name: "other"}
		givenProds := []Product{
			{Code: "001", Price: 1.1, Category: &givenCat},
			{Code: "002", Price: 2.2, Category: &givenCat},
		}

		mockService.On("FindProducts", findCrt).
			Return(givenProds, total, nil)

		// Set routing
		handler := NewCatalogHandler(mockService)
		testMux.HandleFunc("GET /catalog", handler.HandleList)

		// When
		req := httptest.NewRequest("GET", fmt.Sprintf("/catalog?offset=%d&limit=%d&category=%s&priceLessThan=%f",
			offset, limit, cat, priceLimit), nil)
		record := httptest.NewRecorder()
		testMux.ServeHTTP(record, req)

		// Then
		var actualRes CatalogResponse
		err := json.Unmarshal(record.Body.Bytes(), &actualRes)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, record.Code)

		expectedRes := CatalogResponse{
			Total: total,
			Products: []ProductItem{
				{Code: "001", Price: 1.1, Category: "other"},
				{Code: "002", Price: 2.2, Category: "other"},
			},
		}
		assert.Equal(t, expectedRes, actualRes)
	})

	t.Run("Test Invalid GET Params", func(t *testing.T) {
		type testCase struct {
			Name    string
			GetPath string
		}

		invalidParams := []testCase{
			{"Invalid offset", "/catalog?offset=one&limit=10&category=CAT0&priceLessThan=9.99"},
			{"Invalid limit", "/catalog?offset=0&limit=bar&category=CAT0&priceLessThan=9.99"},
			{"Invalid priceLessThan", "/catalog?offset=0&limit=10&category=CAT0&priceLessThan=nil"},
		}

		for _, tCase := range invalidParams {
			t.Run(tCase.Name, func(t *testing.T) {
				mockService := new(MockService)
				testMux := http.NewServeMux()

				// Set routing
				handler := NewCatalogHandler(mockService)
				testMux.HandleFunc("GET /catalog", handler.HandleList)

				// When
				req := httptest.NewRequest("GET", tCase.GetPath, nil)
				record := httptest.NewRecorder()
				testMux.ServeHTTP(record, req)

				// Then
				assert.Equal(t, http.StatusBadRequest, record.Code)
			})
		}
	})

	t.Run("Test Invalid criteria params", func(t *testing.T) {
		mockService := new(MockService)
		testMux := http.NewServeMux()

		// Stage Service Calls
		mockService.On("BuildFindCriteria", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("invalid criteria"))

		// Set routing
		handler := NewCatalogHandler(mockService)
		testMux.HandleFunc("GET /catalog", handler.HandleList)

		// When
		req := httptest.NewRequest("GET", "/catalog", nil)
		record := httptest.NewRecorder()
		testMux.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusBadRequest, record.Code)
	})

	t.Run("Test Internal failure", func(t *testing.T) {
		mockService := new(MockService)
		testMux := http.NewServeMux()

		// Stage Service Calls
		fc := &FindCriteria{}
		mockService.On("BuildFindCriteria", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(fc, nil)

		mockService.On("FindProducts", *fc).Return(nil, uint(0), errors.New("something went wrong"))

		// Set routing
		handler := NewCatalogHandler(mockService)
		testMux.HandleFunc("GET /catalog", handler.HandleList)

		// When
		req := httptest.NewRequest("GET", "/catalog", nil)
		record := httptest.NewRecorder()
		testMux.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusInternalServerError, record.Code)
	})

}

func TestHandleGet(t *testing.T) {

	t.Run("Test Happy path", func(t *testing.T) {
		mockService := new(MockService)
		testMux := http.NewServeMux()

		// Given
		prodCode := "PROD010"

		// Stage Service Calls
		price1 := float64(1.1)
		price2 := float64(2.2)
		givenCat := category.Category{Code: "CAT", Name: "other"}
		givenVars := []Variant{
			{Name: "one", SKU: "558464", Price: &price1},
			{Name: "two", SKU: "558464", Price: &price2},
		}
		givenProd := Product{Code: prodCode, Price: price1, Category: &givenCat, Variants: givenVars}

		mockService.On("GetProduct", prodCode).Return(&givenProd, nil)

		// Set routing
		handler := NewCatalogHandler(mockService)
		testMux.HandleFunc("GET /catalog/{code}", handler.HandleGet)

		// When
		req := httptest.NewRequest("GET", "/catalog/"+prodCode, nil)
		record := httptest.NewRecorder()
		testMux.ServeHTTP(record, req)

		// Then
		var actualRes ProductResponse
		err := json.Unmarshal(record.Body.Bytes(), &actualRes)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, record.Code)

		expectedRes := ProductResponse{
			Code:     prodCode,
			Price:    price1,
			Category: givenCat.Name,
			Variants: []VariantItem{
				{Name: "one", SKU: "558464", Price: price1},
				{Name: "two", SKU: "558464", Price: price2},
			},
		}
		assert.Equal(t, expectedRes, actualRes)
	})

	t.Run("Test Internal failure", func(t *testing.T) {
		mockService := new(MockService)
		testMux := http.NewServeMux()

		mockService.On("GetProduct", mock.Anything).Return(nil, errors.New("unhandled error"))

		// Set routing
		handler := NewCatalogHandler(mockService)
		testMux.HandleFunc("GET /catalog/{code}", handler.HandleGet)

		// When
		req := httptest.NewRequest("GET", "/catalog/bla", nil)
		record := httptest.NewRecorder()
		testMux.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusInternalServerError, record.Code)
	})

	t.Run("Test not found product", func(t *testing.T) {
		mockService := new(MockService)
		testMux := http.NewServeMux()

		mockService.On("GetProduct", mock.Anything).Return(nil, nil)

		// Set routing
		handler := NewCatalogHandler(mockService)
		testMux.HandleFunc("GET /catalog/{code}", handler.HandleGet)

		// When
		req := httptest.NewRequest("GET", "/catalog/foo", nil)
		record := httptest.NewRecorder()
		testMux.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusNotFound, record.Code)
	})
}
