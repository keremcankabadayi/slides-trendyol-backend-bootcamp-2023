package query

type Content struct {
	Id               int64  `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	LastModifiedDate string `json:"lastModifiedDate"`
}
