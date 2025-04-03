package model

import "time"

const (
	ProjectStatusCreated = "created"
	ProjectStatusDeleted = "deleted"
)

// Project is the model for the project table
type Project struct {
	ID              string    `json:"-"`
	UUID            string    `json:"uuid"`
	Name            string    `json:"name"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	AvatarURL       string    `json:"avatar_url"`
	CreatedAt       time.Time `json:"created_at"`
	Status          string    `json:"status"`
	AuthorName      string    `json:"author_name"`
	AuthorAvatarURL string    `json:"author_avatar_url"`
	Tags            string    `json:"tags"`
	Category        string    `json:"category"`
	IsFeatured      bool      `json:"is_featured"`
	Sort            int       `json:"sort"`
	URL             string    `json:"url"`
	Type            string    `json:"type"`
	UserUUID        string    `json:"user_uuid"`
	Tools           string    `json:"tools"`
	SSEURL          string    `json:"sse_url"`
	SSEProvider     string    `json:"sse_provider"`
	SSEParams       string    `json:"sse_params"`
	ServerCommand   string    `json:"server_command"`
	ServerParams    string    `json:"server_params"`
	ServerConfig    string    `json:"server_config"`
	AllowCall       bool      `json:"allow_call"`
}

// TableName returns the table name for the project model
func (p *Project) TableName() string {
	return "projects"
}

// FindProjectByUUID finds a project by its UUID
func FindProjectByUUID(uuid string) (*Project, error) {
	project := &Project{}
	err := db().
		Where("uuid = ?", uuid).
		Where("status = ?", ProjectStatusCreated).
		First(project).Error

	if err != nil {
		return nil, err
	}

	return project, nil
}
