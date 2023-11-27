package clickup

import (
	"connectors/pkg/entities"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const EntityTypeSpace = "space"

type Spaces struct {
	Spaces []Space `json:"spaces"`
}

type Space struct {
	Id                string      `json:"id"`
	Name              string      `json:"name"`
	Private           bool        `json:"private"`
	Color             interface{} `json:"color"`
	Avatar            string      `json:"avatar,omitempty"`
	AdminCanManage    bool        `json:"admin_can_manage,omitempty"`
	Archived          bool        `json:"archived,omitempty"`
	Members           []Member    `json:"members,omitempty"`
	Statuses          []Status    `json:"statuses"`
	MultipleAssignees bool        `json:"multiple_assignees"`
	Features          Features    `json:"features"`
}

type Status struct {
	Status     string `json:"status"`
	Type       string `json:"type"`
	Orderindex int    `json:"orderindex"`
	Color      string `json:"color"`
}

type Features struct {
	DueDates          DueDates       `json:"due_dates"`
	TimeTracking      EnabledFeature `json:"time_tracking"`
	Tags              EnabledFeature `json:"tags"`
	TimeEstimates     EnabledFeature `json:"time_estimates"`
	Checklists        EnabledFeature `json:"checklists"`
	CustomFields      EnabledFeature `json:"custom_fields,omitempty"`
	RemapDependencies EnabledFeature `json:"remap_dependencies,omitempty"`
	DependencyWarning EnabledFeature `json:"dependency_warning,omitempty"`
	Portfolios        EnabledFeature `json:"portfolios,omitempty"`
}

type DueDates struct {
	Enabled            bool `json:"enabled"`
	StartDate          bool `json:"start_date"`
	RemapDueDates      bool `json:"remap_due_dates"`
	RemapClosedDueDate bool `json:"remap_closed_due_date"`
}

type EnabledFeature struct {
	Enabled bool `json:"enabled"`
}

func (s Space) ToEntity(ownerId string) entities.Entity {
	return entities.Entity{
		Name:         s.Name,
		EntityUrl:    "",
		ExternalId:   s.Id,
		Type:         EntityTypeSpace,
		ContentUrl:   "",
		OwnerId:      ownerId,
		LastModified: time.Time{},
		Data:         s,
	}
}

func (c *Client) GetSpaces(ctx context.Context, workspaceID string) ([]Space, error) {
	url := fmt.Sprintf("%steam/%s/space?archived=false", c.url, workspaceID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Add("Authorization", c.apiKey)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data for workspace %s: %w", workspaceID, err)
	}
	defer resp.Body.Close()

	var spaces Spaces
	err = json.NewDecoder(resp.Body).Decode(&spaces)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return spaces.Spaces, nil
}
