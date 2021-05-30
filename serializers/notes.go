package serializers

import "time"

type NoteSerializer struct {
	ID        uint32    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
}

type ContentSerializer struct {
	ID   uint32 `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
	File string `json:"file"`
}
