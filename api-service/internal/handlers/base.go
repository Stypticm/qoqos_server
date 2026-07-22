package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"api-service/internal/repository"
	"github.com/go-chi/chi/v5"
)

type Response struct {
	Data interface{} `json:"data,omitempty"`
}

type ListResponse struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func SendError(w http.ResponseWriter, status int, message string) {
	SendJSON(w, status, ErrorResponse{Error: message})
}

type BaseHandler[T any] struct {
	repo *repository.Repository[T]
}

func NewBaseHandler[T any]() *BaseHandler[T] {
	return &BaseHandler[T]{
		repo: repository.NewRepository[T](),
	}
}

func (h *BaseHandler[T]) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")

	limit := 20
	offset := 0

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	// Handle sorting
	orderBy := query.Get("order_by")
	if orderBy != "" {
		// Quote fields in order_by to handle mixed-case columns in Postgres
		fields := strings.Split(orderBy, ",")
		var quotedFields []string
		for _, f := range fields {
			f = strings.TrimSpace(f)
			if f == "" {
				continue
			}
			parts := strings.Fields(f)
			if len(parts) == 1 {
				// Simple field: "updatedAt" -> "\"updatedAt\""
				quotedFields = append(quotedFields, fmt.Sprintf("\"%s\"", parts[0]))
			} else if len(parts) == 2 {
				// Field with direction: "updatedAt desc" -> "\"updatedAt\" desc"
				quotedFields = append(quotedFields, fmt.Sprintf("\"%s\" %s", parts[0], parts[1]))
			} else {
				// Fallback for complex expressions
				quotedFields = append(quotedFields, f)
			}
		}
		orderBy = strings.Join(quotedFields, ", ")
	} else {
		sort := query.Get("_sort")
		if sort != "" {
			order := query.Get("_order")
			if order == "" {
				order = "asc"
			}
			orderBy = fmt.Sprintf("\"%s\" %s", sort, order)
		}
	}

	// Handle preloads/population
	preloadStr := query.Get("preload")
	if preloadStr == "" {
		preloadStr = query.Get("_populate")
	}
	if preloadStr == "" {
		preloadStr = query.Get("_expand")
	}

	var preloads []string
	if preloadStr != "" {
		preloads = strings.Split(preloadStr, ",")
	}

	// Извлекаем остальные параметры как фильтры
	filters := make(map[string]interface{})
	// Skip known utility params
	skipParams := map[string]bool{
		"limit":     true,
		"offset":    true,
		"order_by":  true,
		"preload":   true,
		"_sort":     true,
		"_order":    true,
		"_start":    true,
		"_end":      true,
		"_populate": true,
		"_expand":   true,
		"_embed":    true,
	}

	for key, values := range query {
		if !skipParams[key] && len(values) > 0 {
			filters[key] = values[0]
		}
	}

	items, total, err := h.repo.GetAll(filters, orderBy, limit, offset, preloads)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendJSON(w, http.StatusOK, ListResponse{Items: items, Total: int64(total)})
}

func (h *BaseHandler[T]) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	query := r.URL.Query()
	preloadStr := query.Get("preload")
	var preloads []string
	if preloadStr != "" {
		preloads = strings.Split(preloadStr, ",")
	}

	item, err := h.repo.GetByID(id, preloads)
	if err != nil {
		SendError(w, http.StatusNotFound, "Resource not found")
		return
	}
	SendJSON(w, http.StatusOK, item)
}

func (h *BaseHandler[T]) Create(w http.ResponseWriter, r *http.Request) {
	var item T
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		SendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Генерируем UUID если ID пустой или передан временный
	// Проверяем через reflection, есть ли поле ID и пустое ли оно
	idField := getIDField(item)
	if idField == "" || idField == "temp" || len(idField) < 10 {
		setIDField(&item, generateUUID())
	}

	if err := h.repo.Create(&item); err != nil {
		SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	SendJSON(w, http.StatusCreated, Response{Data: item})
}

// generateUUID генерирует случайный UUID
func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("id-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (h *BaseHandler[T]) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var item T
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		SendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.repo.Update(id, &item); err != nil {
		SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	SendJSON(w, http.StatusOK, Response{Data: item})
}

func (h *BaseHandler[T]) Patch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		SendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.repo.Patch(id, updates); err != nil {
		SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	// Return updated item
	item, _ := h.repo.GetByID(id, []string{})
	SendJSON(w, http.StatusOK, Response{Data: item})
}

func (h *BaseHandler[T]) Delete(id http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if err := h.repo.Delete(idParam); err != nil {
		SendError(id, http.StatusInternalServerError, err.Error())
		return
	}
	id.WriteHeader(http.StatusNoContent)
}

func (h *BaseHandler[T]) Distinct(w http.ResponseWriter, r *http.Request) {
	field := chi.URLParam(r, "field")
	query := r.URL.Query()

	filters := make(map[string]interface{})
	for key, values := range query {
		if len(values) > 0 {
			filters[key] = values[0]
		}
	}

	results, err := h.repo.GetDistinct(field, filters)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendJSON(w, http.StatusOK, Response{Data: results})
}

// getIDField получает значение поля ID через reflection
func getIDField(item interface{}) string {
	val := reflect.ValueOf(item)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		return ""
	}
	
	if idField.Kind() == reflect.String {
		return idField.String()
	}
	return ""
}

// setIDField устанавливает значение поля ID через reflection
func setIDField(item interface{}, id string) {
	val := reflect.ValueOf(item)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	idField := val.FieldByName("ID")
	if !idField.IsValid() || !idField.CanSet() {
		return
	}
	
	if idField.Kind() == reflect.String {
		idField.SetString(id)
	}
}
