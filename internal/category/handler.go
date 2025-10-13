package category

import (
	"encoding/json"
	"net/http"
)

type CategoryListResponse struct {
	Categories []CategoryItem `json:"categories"`
}

type CategoryItem struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoryHandler struct {
	srv Service
}

func NewCategoryHandler(s Service) *CategoryHandler {
	return &CategoryHandler{
		srv: s,
	}
}

func (h *CategoryHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	cats, err := h.srv.GetAllCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map response
	categories := make([]CategoryItem, len(cats))
	for i, c := range cats {
		categories[i] = mapCategory(c)
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := CategoryListResponse{
		Categories: categories,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CategoryHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var newCatItem CategoryItem
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newCatItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newCat := mapCategoryItem(newCatItem)

	if err := h.srv.CreateCategory(newCat); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func mapCategory(c Category) CategoryItem {
	return CategoryItem{
		Code: c.Code,
		Name: c.Name,
	}
}

func mapCategoryItem(ci CategoryItem) Category {
	return Category{
		Code: ci.Code,
		Name: ci.Name,
	}
}
