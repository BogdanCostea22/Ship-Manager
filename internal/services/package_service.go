package services

import "Ship_Manager/internal/repositories"

// PackageService provides methods to manage and calculate packing for orders.
type PackageService interface {
	// AddPack adds a new pack size to the available pack sizes.
	// It returns an error if the pack size already exists.
	AddPack(size int) error

	// ClearPacks removes all pack sizes from the service.
	ClearPacks()

	// GetPackSizes returns a slice of all available pack sizes, sorted in descending order.
	GetPackSizes() []int

	// CalculatePacks determines the optimal combination of packs for a given order size.
	// It returns a map where the keys are pack sizes and the values are the number of packs needed.
	CalculatePacks(order int) map[int]int
}

type packageService struct {
	repository repositories.PackageRepository
}

// NewPackageService creates a new instance of PackageService with the given repository.
func NewPackageService(repository repositories.PackageRepository) PackageService {
	return &packageService{
		repository: repository,
	}
}

func (ps *packageService) AddPack(size int) error {
	return ps.repository.Add(size)
}

func (ps *packageService) ClearPacks() {
	ps.repository.DeleteAll()
}

func (ps *packageService) GetPackSizes() []int {
	return ps.repository.GetSizes()
}

func (ps *packageService) CalculatePacks(order int) map[int]int {
	result := make(map[int]int)
	remaining := order
	packSizes := ps.repository.GetSizes()

	for _, size := range packSizes {
		count := remaining / size
		if count > 0 {
			result[size] = count
			remaining -= count * size
		}
	}

	if remaining > 0 && len(packSizes) > 0 {
		result[packSizes[len(packSizes)-1]]++
	}

	return result
}
