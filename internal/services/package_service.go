package services

import (
	"Ship_Manager/internal/repositories"
)

type PackSize int

// CalculationResult represents the result of a pack calculation
type CalculationResult struct {
	Packs       map[int]int `json:"packs"`       // Map of pack sizes to the number of packs of that size
	Total       int         `json:"total"`       // Total number of items that will be shipped
	OrderSize   int         `json:"orderSize"`   // Original order size
	ExcessItems int         `json:"excessItems"` // Number of items shipped in excess of the order
	PacksCount  int         `json:"packsCount"`  // Total number of packs used
}

type PackCalculator struct {
	PackSizes []PackSize
}

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

func (ps *packageService) CalculatePacks(orderSize int) map[int]int {
	packSizes := ps.repository.GetSizes()
	largestPack := packSizes[0]
	switch {
	case len(packSizes) == 0:
		{
			return map[int]int{}
		}
		// Handle case where order is smaller than the smallest pack
	case orderSize < packSizes[len(packSizes)-1]:
		{
			return map[int]int{packSizes[len(packSizes)-1]: 1}

		}
	}

	// Initialize dp array
	dp := make([]map[int]int, orderSize+largestPack+1)
	for i := range dp {
		dp[i] = make(map[int]int)
	}

	// Fill dp array
	for i := 1; i <= orderSize+largestPack; i++ {
		for _, size := range packSizes {
			switch {
			case i < size:
				continue
			case i == size:
				dp[i] = map[int]int{size: 1}
			default:
				remainder := i - size
				if len(dp[remainder]) > 0 {
					newSolution, _ := copyMap(dp[remainder])
					newSolution[size]++

					if len(dp[i]) == 0 || len(newSolution) < len(dp[i]) {
						dp[i] = newSolution
					}
				}
			}
		}
	}

	// Find the smallest valid solution
	for i := orderSize; i <= orderSize+largestPack; i++ {
		if len(dp[i]) > 0 {
			return dp[i]
		}
	}

	// If still no solution, return the minimum number of largest packs
	return map[int]int{largestPack: (orderSize + largestPack - 1) / largestPack}
}

func copyMap(originalMap map[int]int) (newSolution map[int]int, totalSpace int) {
	newSolution = make(map[int]int)
	for k, v := range originalMap {
		newSolution[k] = v
		totalSpace += v
	}
	return
}
