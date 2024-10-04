package repositories

import (
	"reflect"
	"testing"
)

func TestPackageRepository(t *testing.T) {
	t.Run("Add and GetSizes", func(t *testing.T) {
		repo := NewPackageRepository()

		// Test adding pack sizes
		sizes := []int{500, 250, 1000, 2000, 5000}
		for _, size := range sizes {
			err := repo.Add(size)
			if err != nil {
				t.Errorf("Failed to add size %d: %v", size, err)
			}
		}

		// Test getting sizes (should be sorted in descending order)
		expected := []int{5000, 2000, 1000, 500, 250}
		actual := repo.GetSizes()
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("GetSizes() = %v, want %v", actual, expected)
		}
	})

	t.Run("Add duplicate size", func(t *testing.T) {
		repo := NewPackageRepository()

		// Add a size
		err := repo.Add(1000)
		if err != nil {
			t.Errorf("Failed to add size 1000: %v", err)
		}

		// Try to add the same size again
		err = repo.Add(1000)
		if err != ErrSizeAlreadyExists {
			t.Errorf("Expected ErrSizeAlreadyExists, got %v", err)
		}
	})

	t.Run("DeleteAll", func(t *testing.T) {
		repo := NewPackageRepository()

		// Add some sizes
		sizes := []int{500, 250, 1000}
		for _, size := range sizes {
			err := repo.Add(size)
			if err != nil {
				t.Errorf("Failed to add size %d: %v", size, err)
			}
		}

		// Delete all sizes
		repo.DeleteAll()

		// Check if the repository is empty
		actual := repo.GetSizes()
		if len(actual) != 0 {
			t.Errorf("Expected empty repository after DeleteAll, got %v", actual)
		}
	})
}
