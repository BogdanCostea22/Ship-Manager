// Package repositories provides functionality for managing package sizes.
package repositories

import (
	"fmt"
	"sort"
	"sync"
)

// ErrSizeAlreadyExists is returned when attempting to add a package size that already exists.
var ErrSizeAlreadyExists = fmt.Errorf("pack size already exists")

// packCache represents the in-memory storage for pack sizes.
type packCache struct {
	packSizes []int
	mu        sync.Mutex
}

// PackageRepository defines the interface for managing package sizes.
type PackageRepository interface {
	// Add inserts a new pack size into the repository.
	// It returns an error if the size already exists.
	Add(size int) error

	// DeleteAll removes all pack sizes from the repository.
	DeleteAll()

	// GetSizes returns a slice of all pack sizes in descending order.
	GetSizes() []int
}

// packageRepository implements the PackageRepository interface.
type packageRepository struct {
	cache *packCache
}

// NewPackageRepository creates and returns a new instance of PackageRepository.
func NewPackageRepository() PackageRepository {
	pc := packCache{
		packSizes: []int{},
	}
	return &packageRepository{
		cache: &pc,
	}
}

// Add inserts a new pack size into the repository.
// The sizes are maintained in descending order.
// It returns ErrSizeAlreadyExists if the size is already in the repository.
func (pr *packageRepository) Add(size int) error {
	pr.cache.mu.Lock()
	defer pr.cache.mu.Unlock()

	// Find the correct position for the new size
	index := sort.Search(len(pr.cache.packSizes), func(i int) bool {
		return pr.cache.packSizes[i] <= size
	})

	// Check if size already exists
	if index < len(pr.cache.packSizes) && pr.cache.packSizes[index] == size {
		return ErrSizeAlreadyExists
	}

	// Insert the new size at the correct position
	pr.cache.packSizes = append(pr.cache.packSizes, 0)
	copy(pr.cache.packSizes[index+1:], pr.cache.packSizes[index:])
	pr.cache.packSizes[index] = size

	return nil
}

// DeleteAll removes all pack sizes from the repository.
func (pr *packageRepository) DeleteAll() {
	pr.cache.mu.Lock()
	defer pr.cache.mu.Unlock()
	pr.cache.packSizes = []int{}
}

// GetSizes returns a copy of all pack sizes in descending order.
func (pr *packageRepository) GetSizes() []int {
	pr.cache.mu.Lock()
	defer pr.cache.mu.Unlock()
	return append([]int{}, pr.cache.packSizes...)
}
