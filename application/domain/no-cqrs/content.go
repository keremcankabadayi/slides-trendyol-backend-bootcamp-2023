package no_cqrs

type Content struct {
	Id               int64  `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Culture          string `json:"culture"`
	CreationDate     string `json:"creationDate"`
	LastModifiedDate string `json:"lastModifiedDate"`
	CreatedBy        string `json:"createdBy"`
	LastModifiedBy   string `json:"lastModifiedBy"`
	IndexedAt        string `json:"indexedAt"`
}
