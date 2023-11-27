package clickup

import (
	"connectors/pkg/entities"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const EntityTypeFolder = "folder"

type Folders struct {
	Folders []Folder `json:"folders"`
}

type Folder struct {
	Id               string      `json:"id"`
	Name             string      `json:"name"`
	OrderIndex       int         `json:"orderindex"`
	OverrideStatuses bool        `json:"override_statuses"`
	Hidden           bool        `json:"hidden"`
	Space            FolderSpace `json:"space"`
	TaskCount        string      `json:"task_count"`
	Lists            []List      `json:"lists"`
}

type FolderSpace struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Access bool   `json:"access"`
}

func (c *Client) GetFolders(ctx context.Context, spaceId string) ([]Folder, error) {
	url := fmt.Sprintf("%sspace/%s/folder?archived=false", c.url, spaceId)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Add("Authorization", c.apiKey)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve folder data for spaceId %s: %w", spaceId, err)
	}
	defer resp.Body.Close()

	var spaces Folders
	err = json.NewDecoder(resp.Body).Decode(&spaces)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response for spaceId %s: %w", spaceId, err)
	}

	return spaces.Folders, nil
}

func (f Folder) ToEntity(ownerId string) entities.Entity {
	return entities.Entity{
		Name:         f.Name,
		EntityUrl:    "",
		ExternalId:   f.Id,
		Type:         EntityTypeFolder,
		ContentUrl:   "",
		OwnerId:      ownerId,
		LastModified: time.Time{},
		Data:         f,
	}
}
