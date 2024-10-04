package handlers

import (
	"Ship_Manager/cmd/web"
	"Ship_Manager/internal/repositories"
	"Ship_Manager/internal/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

// PackageHandler is responsible for handling HTTP requests related to package management.
// It provides methods for adding pack sizes, calculating packs for orders, clearing all packs,
// and retrieving pack sizes.
type PackageHandler struct {
	service services.PackageService
}

// NewPackageHandler creates a new instance of PackageHandler with the given PackageService.
func NewPackageHandler(service services.PackageService) *PackageHandler {
	return &PackageHandler{
		service: service,
	}
}

// AddPack handles POST requests to add a new pack size.
// It expects a form value "size" with the pack size to add.
// Returns HTTP 409 if the pack size already exists.
// Triggers "packSizesChanged" event on success.
func (ph *PackageHandler) AddPack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sizeStr := r.FormValue("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		http.Error(w, "Invalid pack size", http.StatusBadRequest)
		return
	}

	err = ph.service.AddPack(size)

	if err != nil {
		var statusCode int
		var errorMessage string

		switch err {
		case repositories.ErrSizeAlreadyExists:
			statusCode = http.StatusConflict
			errorMessage = "Pack size already exists"
		default:
			statusCode = http.StatusInternalServerError
			errorMessage = "An error occurred while adding the pack size"
		}

		w.WriteHeader(statusCode)
		w.Header().Set("HX-Trigger", "errorMessage")
		templ.Handler(web.ErrorMessage(errorMessage)).ServeHTTP(w, r)
		return
	}

	w.Header().Set("HX-Trigger", "packSizesChanged")
	templ.Handler(web.PackSizesList(ph.service.GetPackSizes())).ServeHTTP(w, r)
}

// Calculate handles POST requests to calculate packs for an order.
// It expects a form value "order" with the order size.
// Returns a JSON response with the calculated packs.
func (ph *PackageHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderStr := r.FormValue("order")
	order, err := strconv.Atoi(orderStr)
	if err != nil {
		http.Error(w, "Invalid order size", http.StatusBadRequest)
		return
	}

	result := ph.service.CalculatePacks(order)

	jsonResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Error encoding result", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}

// ClearPacks handles POST requests to clear all pack sizes.
// Triggers "packSizesChanged" event on success.
func (ph *PackageHandler) ClearPacks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ph.service.ClearPacks()

	w.Header().Set("HX-Trigger", "packSizesChanged")
	templ.Handler(web.PackSizesList(ph.service.GetPackSizes())).ServeHTTP(w, r)
}

// PackSizes handles requests to retrieve all pack sizes.
// Returns an HTML component with the list of pack sizes.
func (ph *PackageHandler) PackSizes(w http.ResponseWriter, r *http.Request) {
	templ.Handler(web.PackSizesList(ph.service.GetPackSizes())).ServeHTTP(w, r)
}

// CalculatorIndex handles requests for the main calculator page.
// Returns the HTML for the calculator index page.
func (ph *PackageHandler) CalculatorIndex(w http.ResponseWriter, r *http.Request) {
	templ.Handler(web.IndexPage(ph.service.GetPackSizes())).ServeHTTP(w, r)
}
