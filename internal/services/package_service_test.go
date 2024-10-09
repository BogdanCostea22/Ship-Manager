package services_test

import (
	"Ship_Manager/internal/repositories"
	"Ship_Manager/internal/services"
	"reflect"
	"testing"
)

func TestPackageService(t *testing.T) {
	repo := repositories.NewPackageRepository()
	service := services.NewPackageService(repo)

	t.Run("AddPack", func(t *testing.T) {
		err := service.AddPack(250)
		if err != nil {
			t.Errorf("Failed to add pack: %v", err)
		}

		sizes := service.GetPackSizes()
		if len(sizes) != 1 || sizes[0] != 250 {
			t.Errorf("Expected pack sizes [250], got %v", sizes)
		}
	})

	t.Run("AddPackAndCheckIfThereAreInTheRightOrder", func(t *testing.T) {
		service.ClearPacks()
		service.AddPack(250)
		service.AddPack(500)

		sizes := service.GetPackSizes()
		if len(sizes) != 2 || sizes[0] != 500 || sizes[1] != 250 {
			t.Errorf("Expected pack sizes [250], got %v", sizes)
		}
	})

	t.Run("ClearPacks", func(t *testing.T) {
		service.AddPack(500)
		service.ClearPacks()

		sizes := service.GetPackSizes()
		if len(sizes) != 0 {
			t.Errorf("Expected empty pack sizes, got %v", sizes)
		}
	})

	t.Run("GetPackSizes", func(t *testing.T) {
		service.ClearPacks()
		service.AddPack(250)
		service.AddPack(500)
		service.AddPack(1000)

		expected := []int{1000, 500, 250}
		sizes := service.GetPackSizes()
		if !reflect.DeepEqual(sizes, expected) {
			t.Errorf("Expected pack sizes %v, got %v", expected, sizes)
		}
	})

	t.Run("CalculatePacks", func(t *testing.T) {
		service.ClearPacks()
		service.AddPack(250)
		service.AddPack(500)
		service.AddPack(1000)
		service.AddPack(2000)
		service.AddPack(5000)

		testCases := []struct {
			order    int
			expected map[int]int
		}{
			{1, map[int]int{250: 1}},
			{250, map[int]int{250: 1}},
			{251, map[int]int{500: 1}},
			{501, map[int]int{500: 1, 250: 1}},
			{12001, map[int]int{5000: 2, 2000: 1, 250: 1}},
		}

		for _, tc := range testCases {
			result := service.CalculatePacks(tc.order)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("For order %d, expected %v, got %v", tc.order, tc.expected, result)
			}
		}
	})

	t.Run("Calculate optimal order", func(t *testing.T) {
		service.ClearPacks()
		service.AddPack(5)
		service.AddPack(12)

		testCases := []struct {
			order    int
			expected map[int]int
		}{
			{15, map[int]int{5: 3}},
			{18, map[int]int{5: 4}},
		}
		for _, tc := range testCases {
			result := service.CalculatePacks(tc.order)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("For order %d, expected %v, got %v", tc.order, tc.expected, result)
			}
		}
	})
}
