package models

type UpcLookup struct {
	Type        string `json:"upc"`
	ReferableId string `json:"referable_id"`
}

type Upc struct {
	ID            int64  `json:"id"`
	Upc           string `json:"upc"`
	Referable     string `json:"referable_type"`
	ReferableId   string `json:"referable_id"`
	ReferableName string `json:"referable_name"`
}
