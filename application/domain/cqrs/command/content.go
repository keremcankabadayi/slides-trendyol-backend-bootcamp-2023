package command

type Content struct {
	Id               int64         `json:"id"`
	Name             string        `json:"name"`
	Description      string        `json:"description"`
	Status           ContentStatus `json:"status"`
	Culture          string        `json:"culture"`
	CreationDate     string        `json:"creationDate"`
	CreatedBy        string        `json:"createdBy"`
	LastModifiedDate string        `json:"lastModifiedDate"`
	LastModifiedBy   string        `json:"lastModifiedBy"`
}

type ContentStatus struct {
	Id               int64  `json:"id"`
	LastModifiedDate string `json:"lastModifiedDate"`
	LastModifiedBy   string `json:"lastModifiedBy"`
}
