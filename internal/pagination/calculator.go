// pagination/calculator.go

package pagination

import (
	"fmt"
	"math"
)

// CalculateState computes all the pagination data points.
func CalculateState(currentPage int, limit int, totalRecords int) PaginationState {
	// Ensure limit and current page are valid non-zero values
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if currentPage <= 0 {
		currentPage = 1 // Default to first page
	}
	if totalRecords < 0 {
		totalRecords = 0
	}

	// 1. Calculate Total Pages
	// Use math.Ceil to round up (e.g., 11 records with a limit of 10 is 2 pages)
	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	// If there are no records, there is 1 total page (the empty page)
	if totalRecords == 0 {
		totalPages = 1
	}

	// Adjust CurrentPage if it somehow exceeds TotalPages
	if currentPage > totalPages {
		currentPage = totalPages
	}

	// 2. Determine Navigation Page Numbers
	firstPage := 1
	lastPage := totalPages

	prevPage := currentPage - 1
	if prevPage < firstPage {
		prevPage = firstPage // Should be disabled, but set to 1
	}

	nextPage := currentPage + 1
	if nextPage > lastPage {
		nextPage = lastPage // Should be disabled, but set to TotalPages
	}

	// 3. Determine Disabled States
	isFirstDisabled := currentPage == firstPage
	isPrevDisabled := currentPage == firstPage
	isNextDisabled := currentPage == lastPage
	isLastDisabled := currentPage == lastPage

	// URL Values
	firstPageURL := fmt.Sprintf("/api/update?page=%d", firstPage)
	prevPageURL := fmt.Sprintf("/api/update?page=%d", prevPage)
	nextPageURL := fmt.Sprintf("/api/update?page=%d", nextPage)
	lastPageURL := fmt.Sprintf("/api/update?page=%d", lastPage)

	return PaginationState{
		CurrentPage:  currentPage,
		TotalRecords: totalRecords,
		PageLimit:    limit,
		TotalPages:   totalPages,

		FirstPageNum: firstPage,
		PrevPageNum:  prevPage,
		NextPageNum:  nextPage,
		LastPageNum:  lastPage,

		IsFirstDisabled: isFirstDisabled,
		IsPrevDisabled:  isPrevDisabled,
		IsNextDisabled:  isNextDisabled,
		IsLastDisabled:  isLastDisabled,

		FirstPageURL: firstPageURL,
		PrevPageURL:  prevPageURL,
		NextPageURL:  nextPageURL,
		LastPageURL:  lastPageURL,
	}
}
