package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPackageService is a mock of PackageService
type MockPackageService struct {
	mock.Mock
}

func (m *MockPackageService) AddPack(size int) error {
	args := m.Called(size)
	return args.Error(0)
}

func (m *MockPackageService) ClearPacks() {
	m.Called()
}

func (m *MockPackageService) GetPackSizes() []int {
	args := m.Called()
	return args.Get(0).([]int)
}

func (m *MockPackageService) CalculatePacks(order int) map[int]int {
	args := m.Called(order)
	return args.Get(0).(map[int]int)
}

func TestAddPack(t *testing.T) {
	mockService := new(MockPackageService)
	handler := NewPackageHandler(mockService)

	t.Run("Successful add", func(t *testing.T) {
		mockService.On("AddPack", 100).Return(nil).Once()
		mockService.On("GetPackSizes").Return([]int{100}).Once()

		form := url.Values{}
		form.Add("size", "100")
		req, _ := http.NewRequest("POST", "/add-pack", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler.AddPack(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Header().Get("HX-Trigger"), "packSizesChanged")
	})

	t.Run("Invalid size", func(t *testing.T) {
		form := url.Values{}
		form.Add("size", "invalid")
		req, _ := http.NewRequest("POST", "/add-pack", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler.AddPack(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestCalculate(t *testing.T) {
	mockService := new(MockPackageService)
	handler := NewPackageHandler(mockService)

	t.Run("Successful calculation", func(t *testing.T) {
		mockService.On("CalculatePacks", 250).Return(map[int]int{250: 1}).Once()

		form := url.Values{}
		form.Add("order", "250")
		req, _ := http.NewRequest("POST", "/calculate", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler.Calculate(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var result map[int]int
		json.NewDecoder(rr.Body).Decode(&result)
		assert.Equal(t, map[int]int{250: 1}, result)
	})

	t.Run("Invalid order size", func(t *testing.T) {
		form := url.Values{}
		form.Add("order", "invalid")
		req, _ := http.NewRequest("POST", "/calculate", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler.Calculate(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestClearPacks(t *testing.T) {
	mockService := new(MockPackageService)
	handler := NewPackageHandler(mockService)

	mockService.On("ClearPacks").Once()
	mockService.On("GetPackSizes").Return([]int{}).Once()

	req, _ := http.NewRequest("POST", "/clear-packs", nil)
	rr := httptest.NewRecorder()

	handler.ClearPacks(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("HX-Trigger"), "packSizesChanged")
}

func TestPackSizes(t *testing.T) {
	mockService := new(MockPackageService)
	handler := NewPackageHandler(mockService)

	mockService.On("GetPackSizes").Return([]int{100, 250, 500}).Once()

	req, _ := http.NewRequest("GET", "/pack-sizes", nil)
	rr := httptest.NewRecorder()

	handler.PackSizes(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "100")
	assert.Contains(t, rr.Body.String(), "250")
	assert.Contains(t, rr.Body.String(), "500")
}

func TestCalculatorIndex(t *testing.T) {
	mockService := new(MockPackageService)
	handler := NewPackageHandler(mockService)

	mockService.On("GetPackSizes").Return([]int{100, 250, 500}).Once()

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.CalculatorIndex(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Calculator")
}
