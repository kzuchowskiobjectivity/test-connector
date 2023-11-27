package clickup

import (
	"connectors/pkg/entities"

	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const EntityTypeWorkspace = "workspace"

type Workspaces struct {
	Teams []Team `json:"teams"`
}

type Team struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Color   string   `json:"color"`
	Avatar  string   `json:"avatar"`
	Members []Member `json:"members"`

	Url string `json:"-"`
}

type Member struct {
	User User `json:"user"`
}

type User struct {
	Id             int    `json:"id"`
	Username       string `json:"username"`
	Color          string `json:"color"`
	ProfilePicture string `json:"profilePicture"`
	Initials       string `json:"initials,omitempty"`
}

func (c *Client) GetWorkspaces(ctx context.Context) ([]Team, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+"team", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Add("Authorization", c.apiKey)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve workspace data: %w", err)
	}
	defer resp.Body.Close()

	var workspaces Workspaces
	err = json.NewDecoder(resp.Body).Decode(&workspaces)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return workspaces.Teams, nil
}

func (t Team) ToEntity(ownerId string) entities.Entity {
	return entities.Entity{
		Name:         t.Name,
		EntityUrl:    "",
		ExternalId:   t.Id,
		Type:         EntityTypeWorkspace,
		ContentUrl:   "",
		OwnerId:      ownerId,
		LastModified: time.Time{},
		Data:         t,
	}
}
