package product

import (
	"encoding/json"
	"math"
	"net/http"
)

type Response struct {
	Products []ProductItem `json:"products"`
}

type ProductItem struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type CatalogHandler struct {
	srv *Service
}

func NewCatalogHandler(s *Service) *CatalogHandler {
	return &CatalogHandler{
		srv: s,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	res, err := h.srv.FindProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map response
	products := make([]ProductItem, len(res))
	for i, p := range res {
		products[i] = mapProduct(p)
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Products: products,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func mapProduct(p Product) ProductItem {
	item := ProductItem{
		Code:  p.Code,
		Price: math.Round(p.Price*100) / 100,
	}

	if p.Category != nil {
		item.Category = p.Category.Name
	}

	return item
}
