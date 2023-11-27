package clickup

import (
	"connectors/pkg/entities"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const EntityTypeTask = "task"

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Id                  string        `json:"id"`
	Name                string        `json:"name"`
	Status              Status        `json:"status"`
	MarkdownDescription string        `json:"markdown_description,omitempty"`
	Orderindex          string        `json:"orderindex"`
	DateCreated         string        `json:"date_created"`
	DateUpdated         string        `json:"date_updated"`
	DateClosed          interface{}   `json:"date_closed"`
	DateDone            interface{}   `json:"date_done"`
	Creator             User          `json:"creator"`
	Assignees           []interface{} `json:"assignees"`
	Checklists          []interface{} `json:"checklists"`
	Tags                []interface{} `json:"tags"`
	Parent              interface{}   `json:"parent"`
	Priority            interface{}   `json:"priority"`
	DueDate             interface{}   `json:"due_date"`
	StartDate           interface{}   `json:"start_date"`
	TimeEstimate        interface{}   `json:"time_estimate"`
	TimeSpent           interface{}   `json:"time_spent"`
	List                StringId      `json:"list"`
	Folder              StringId      `json:"folder"`
	Space               StringId      `json:"space"`
	Url                 string        `json:"url"`
}

type StringId struct {
	Id string `json:"id"`
}

func (t Task) ToEntity(ownerId string) entities.Entity {
	return entities.Entity{
		Name:         t.Name,
		EntityUrl:    "",
		ExternalId:   t.Id,
		Type:         EntityTypeTask,
		ContentUrl:   t.Url,
		OwnerId:      ownerId,
		LastModified: time.Time{},
		Data:         t,
	}
}

func (c *Client) GetTasks(ctx context.Context, listId string) ([]Task, error) {
	url := fmt.Sprintf("%slist/%s/task?archive=true", c.url, listId)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Add("Authorization", c.apiKey)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do http request: %w", err)
	}
	defer resp.Body.Close()

	var spaces Tasks
	err = json.NewDecoder(resp.Body).Decode(&spaces)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response for listId %s: %w", listId, err)
	}

	return spaces.Tasks, nil
}
