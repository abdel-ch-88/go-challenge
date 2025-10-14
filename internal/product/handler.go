package product

import (
	"encoding/json"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

type CatalogResponse struct {
	Total    uint          `jons:"total"`
	Products []ProductItem `json:"products"`
}

type ProductItem struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type ProductResponse struct {
	Code     string        `json:"code"`
	Price    float64       `json:"price"`
	Category string        `json:"category"`
	Variants []VariantItem `json:"variants"`
}

type VariantItem struct {
	Name  string  `json:"name"`
	SKU   string  `json:"SKU"`
	Price float64 `json:"price"`
}

type CatalogHandler struct {
	srv Service
}

func NewCatalogHandler(s Service) *CatalogHandler {
	return &CatalogHandler{
		srv: s,
	}
}

func (h *CatalogHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	findCrt, err := h.extrctParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, total, err := h.srv.FindProducts(*findCrt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map response
	products := make([]ProductItem, len(res))
	for i, p := range res {
		products[i] = mapProductToItem(p)
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := CatalogResponse{
		Total:    total,
		Products: products,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	prodCode := r.PathValue("code")
	if prodCode == "" {
		http.Error(w, "missing mandatory path param: code", http.StatusBadRequest)
		return
	}

	res, err := h.srv.GetProduct(prodCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res == nil {
		http.Error(w, "no product was found with the provided code", http.StatusNotFound)
		return
	}

	// Map response
	response := mapProductToResponse(*res)

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CatalogHandler) extrctParams(params url.Values) (*FindCriteria, error) {
	var err error
	var offset *uint
	var limit *uint
	var catCode *string
	var priceLessThan *float64

	// Check and validate the format of each GET param
	offset, err = parseUint(params.Get("offset"))
	if err != nil {
		return nil, err
	}

	limit, err = parseUint(params.Get("limit"))
	if err != nil {
		return nil, err
	}

	if params.Has("category") {
		catVal := params.Get("category")
		catCode = &catVal
	}

	priceLessThan, err = parseFloat(params.Get("priceLessThan"))
	if err != nil {
		return nil, err
	}

	// Check the busines logic validations of each param
	return h.srv.BuildFindCriteria(offset, limit, catCode, priceLessThan)
}

func parseUint(val string) (*uint, error) {
	if val == "" {
		return nil, nil
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return nil, err
	}

	res := uint(intVal)
	return &res, nil
}

func parseFloat(val string) (*float64, error) {
	if val == "" {
		return nil, nil
	}

	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return nil, err
	}

	return &floatVal, nil
}

func mapProductToItem(p Product) ProductItem {
	item := ProductItem{
		Code:  p.Code,
		Price: math.Round(p.Price*100) / 100,
	}

	if p.Category != nil {
		item.Category = p.Category.Name
	}

	return item
}

func mapProductToResponse(p Product) ProductResponse {
	res := ProductResponse{
		Code:  p.Code,
		Price: p.Price,
	}

	if p.Category != nil {
		res.Category = p.Category.Name
	}
	res.Variants = mapVariantToItem(p.Variants)

	return res
}

func mapVariantToItem(vs []Variant) []VariantItem {
	if vs == nil {
		return nil
	}

	items := make([]VariantItem, len(vs))

	for i, v := range vs {
		items[i] = VariantItem{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: *v.Price,
		}
	}

	return items
}
