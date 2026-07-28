package handlers

import (
	"encoding/json"
	"net/http"

	"api-service/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type StockHandler struct {
	stockService *services.StockService
}

func NewStockHandler() *StockHandler {
	return &StockHandler{
		stockService: services.NewStockService(),
	}
}

type AdjustStockRequest struct {
	PartName    string  `json:"partName" binding:"required"`
	PartModel   string  `json:"partModel" binding:"required"`
	Delta       int     `json:"delta" binding:"required"`
	Reason      string  `json:"reason" binding:"required"`
	ReferenceId *string `json:"referenceId"`
}

func (h *StockHandler) AdjustStock(w http.ResponseWriter, r *http.Request) {
	var req AdjustStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	err := h.stockService.AdjustStock(req.PartName, req.PartModel, req.Delta, req.Reason, req.ReferenceId)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to adjust stock"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Stock adjusted successfully"})
}

func (h *StockHandler) GetStockReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.stockService.GetStockReport()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to get stock report"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, report)
}

func (h *StockHandler) GetMovements(w http.ResponseWriter, r *http.Request) {
	partName := r.URL.Query().Get("partName")
	partModel := r.URL.Query().Get("partModel")

	movements, err := h.stockService.GetMovements(partName, partModel)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to get movements"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, movements)
}

type CheckAvailabilityRequest struct {
	PartName    string `json:"partName" binding:"required"`
	PartModel   string `json:"partModel" binding:"required"`
	RequiredQty int    `json:"requiredQty"`
}

func (h *StockHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	var req CheckAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	if req.RequiredQty <= 0 {
		req.RequiredQty = 1
	}

	available, current, err := h.stockService.IsAvailable(req.PartName, req.PartModel, req.RequiredQty)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to check availability"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"partName":  req.PartName,
		"partModel": req.PartModel,
		"available": available,
		"current":   current,
		"required":  req.RequiredQty,
	})
}

type BulkAdjustRequest struct {
	PartName    string  `json:"partName"`
	PartModel   string  `json:"partModel"`
	Delta       int     `json:"delta"`
	Reason      string  `json:"reason"`
	ReferenceId *string `json:"referenceId"`
}

func (h *StockHandler) BulkAdjustFromCSV(w http.ResponseWriter, r *http.Request) {
	var rows []BulkAdjustRequest
	if err := json.NewDecoder(r.Body).Decode(&rows); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	for _, row := range rows {
		if err := h.stockService.AdjustStock(row.PartName, row.PartModel, row.Delta, row.Reason, row.ReferenceId); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to bulk adjust stock"})
			return
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]int{"count": len(rows)})
}

type SeedStockRequest struct {
	PartName     string `json:"partName"`
	PartModel    string `json:"partModel"`
	Quantity     int    `json:"quantity"`
	DeliveryDays int    `json:"deliveryDays"`
}

func (h *StockHandler) SeedInitialStock(w http.ResponseWriter, r *http.Request) {
	var stocks []SeedStockRequest
	if err := json.NewDecoder(r.Body).Decode(&stocks); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	for _, stock := range stocks {
		err := h.stockService.AdjustStock(stock.PartName, stock.PartModel, stock.Quantity, "initial_seed", nil)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to seed stock"})
			return
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Initial stock seeded successfully"})
}

func SetupStockRoutes(r chi.Router) {
	h := NewStockHandler()
	r.Route("/stock", func(r chi.Router) {
		r.Post("/adjust", h.AdjustStock)
		r.Get("/report", h.GetStockReport)
		r.Get("/movements", h.GetMovements)
		r.Post("/check", h.CheckAvailability)
		r.Post("/bulk", h.BulkAdjustFromCSV)
		r.Post("/seed", h.SeedInitialStock)
	})
}
