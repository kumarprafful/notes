package models

import (
	"errors"
	"html"
	"notes/serializers"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Note struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	User      User      `json:"user"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	Title     string    `gorm:"size:255" json:"title"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Content struct {
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Note      Note      `json:"note"`
	NoteID    uint32    `gorm:"not null" json:"note_id"`
	Type      string    `gorm:"size:20" json:"type"`
	Text      string    `json:"text"`
	File      string    `gorm:"size:255" json:"file"`
}

var content_type = map[string]string{
	"text":  "text",
	"image": "image",
	"audio": "audio",
	"video": "video",
}

func (n *Note) Prepare() {
	n.ID = 0
	n.Title = html.EscapeString(strings.TrimSpace(n.Title))
	n.User = User{}
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
}

func (c *Content) Prepare() {
	c.ID = 0
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	c.Note = Note{}
	c.Type = html.EscapeString(strings.TrimSpace(c.Type))
	c.Text = html.EscapeString(strings.TrimSpace(c.Text))
	c.File = html.EscapeString(strings.TrimSpace(c.File))
}

func (n *Note) Validate() error {
	if n.UserID < 1 {
		return errors.New("user is required")
	}
	return nil
}

func (c *Content) Validate() error {
	if c.Type == "" {
		return errors.New("type is required")
	} else {
		if _, ok := content_type[c.Type]; !ok {
			return errors.New("type not valid")
		}
	}
	if c.Type != "text" {
		if c.File == "" {
			return errors.New("file is required")
		}
	}
	if c.NoteID < 1 {
		return errors.New("notes is required")
	} else {
		c.NoteID = uint32(c.NoteID)
	}
	return nil
}

func (n *Note) SaveNote(db *gorm.DB) (*Note, error) {
	err := db.Debug().Model(&Note{}).Create(&n).Error
	if err != nil {
		return &Note{}, err
	}
	if n.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", n.UserID).Take(&n.User).Error
		if err != nil {
			return &Note{}, err
		}
	}
	return n, nil
}

func (c *Content) SaveContent(db *gorm.DB) (*Content, error) {
	err := db.Debug().Model(&Content{}).Create(&c).Error
	if err != nil {
		return &Content{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&Note{}).Where("id = ?", c.NoteID).Take(&c.Note).Error
		if err != nil {
			return &Content{}, nil
		}
	}
	return c, nil
}

func (n *Note) GetNoteByID(db *gorm.DB, note_id uint32) (*serializers.NoteSerializer, error) {
	note := serializers.NoteSerializer{}
	err := db.Debug().Model(&Note{}).Where("id = ?", note_id).Find(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (n *Note) ContentsOfNotes(db *gorm.DB, note_id uint32) (*[]serializers.ContentSerializer, error) {
	contents := []serializers.ContentSerializer{}
	err := db.Debug().Model(&Content{}).Where("note_id=?", note_id).Limit(100).Find(&contents).Error
	if err != nil {
		return nil, err
	}
	return &contents, nil
}
