// pagination/state.go

package pagination

// PaginationState holds all the calculated data needed to render pagination controls.
type PaginationState struct {
	// Core Data
	CurrentPage  int
	TotalRecords int
	PageLimit    int
	TotalPages   int

	// Navigation States
	FirstPageNum int // The page number for the 'First Page' button (usually 1)
	PrevPageNum  int // The page number for the 'Previous Page' button
	NextPageNum  int // The page number for the 'Next Page' button
	LastPageNum  int // The page number for the 'Last Page' button

	// Disabled Flags (used to disable HTMX buttons)
	IsFirstDisabled bool
	IsPrevDisabled  bool
	IsNextDisabled  bool
	IsLastDisabled  bool

	// URL values
	FirstPageURL string
	PrevPageURL  string
	NextPageURL  string
	LastPageURL  string
}
