package models

type Component struct {
	ID         string  `json:"_id"`
	Name       string  `json:"name"`
	CategoryID int     `json:"categoryId"`
	Price      float32 `json:"price"`
	IsDeleted  int
}

type ProtocolComponent struct {
	ProtocolComponentID string `json:"_id"`
	ComponetID          string `json:"¨componentId"`
	NameComponent       string `json:"namecomponent"`
	ProtocolID          string `json:"protocolId"`
	Price               string `json:"price"`
}
