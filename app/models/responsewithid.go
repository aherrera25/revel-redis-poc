package models

type ResponseWithId struct {
	Response
	ContentId string `json:"content_id"`
}
