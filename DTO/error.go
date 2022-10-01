package dto

type Error struct {
	Status Status `json:"status"`
	Error  string `json:"error"`
}
