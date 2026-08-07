package id

import "github.com/nrednav/cuid2"

// New generates a new cuid2 ID.
// Used to create IDs for new records, matching the cuid format used by the frontend.
func New() string {
	return cuid2.Generate()
}
