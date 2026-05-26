package foodtosave

import "time"

// Gondola corresponds to an item returned by /merchants/{id}/gondolas.
type Gondola struct {
	ID                    string    `json:"id"`
	Quantity              int       `json:"quantity"`
	ShowcaseGroupValue    float64   `json:"showcase_group_value"`
	Bag                   Bag       `json:"bag"`
	AvailabilityStartFrom time.Time `json:"availability_start_from"`
	AvailabilityEndAt     time.Time `json:"availability_end_at"`
}

// Bag corresponds to the bag associated with the gondola.
type Bag struct {
	ID             int     `json:"id"`
	PartnerID      string  `json:"partner_id"`
	Category       string  `json:"category"`
	Type           string  `json:"type"`
	Description    string  `json:"description"`
	Content        string  `json:"content"`
	Price          float64 `json:"price"`
	ReferencePrice float64 `json:"reference_price"`
}
