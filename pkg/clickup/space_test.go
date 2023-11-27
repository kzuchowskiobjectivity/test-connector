package clickup

import (
	"connectors/pkg/entities"
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const correctSpacesJson = "{\"spaces\":[{\"id\":\"90150340588\",\"name\":\"Team Space\",\"color\":\"#03A2FD\",\"private\":false,\"avatar\":null,\"admin_can_manage\":true,\"statuses\":[{\"id\":\"p90150340588_p90100401420_p32280339_DCnaeiSB\",\"status\":\"to do\",\"type\":\"open\",\"orderindex\":0,\"color\":\"#87909e\"},{\"id\":\"p90150340588_p90100401420_p32280339_g9uxhsQM\",\"status\":\"in progress\",\"type\":\"custom\",\"orderindex\":1,\"color\":\"#1090e0\"},{\"id\":\"p90150340588_p90100401420_p32280339_syqLtYOY\",\"status\":\"complete\",\"type\":\"closed\",\"orderindex\":2,\"color\":\"#008844\"}],\"multiple_assignees\":true,\"features\":{\"due_dates\":{\"enabled\":true,\"start_date\":true,\"remap_due_dates\":false,\"remap_closed_due_date\":false},\"sprints\":{\"enabled\":false},\"time_tracking\":{\"enabled\":true,\"harvest\":false,\"rollup\":false},\"points\":{\"enabled\":false},\"custom_items\":{\"enabled\":false},\"priorities\":{\"enabled\":true,\"priorities\":[{\"color\":\"#f50000\",\"id\":\"1\",\"orderindex\":\"1\",\"priority\":\"urgent\"},{\"color\":\"#f8ae00\",\"id\":\"2\",\"orderindex\":\"2\",\"priority\":\"high\"},{\"color\":\"#6fddff\",\"id\":\"3\",\"orderindex\":\"3\",\"priority\":\"normal\"},{\"color\":\"#d8d8d8\",\"id\":\"4\",\"orderindex\":\"4\",\"priority\":\"low\"}]},\"tags\":{\"enabled\":true},\"check_unresolved\":{\"enabled\":true,\"subtasks\":null,\"checklists\":null,\"comments\":null},\"zoom\":{\"enabled\":true},\"milestones\":{\"enabled\":false},\"custom_fields\":{\"enabled\":true},\"dependency_warning\":{\"enabled\":true},\"status_pies\":{\"enabled\":false},\"multiple_assignees\":{\"enabled\":true}},\"archived\":false}]}"

func TestSpace_ToEntity(t *testing.T) {
	space := Space{
		Id:             "45678965",
		Name:           "Universe",
		Private:        false,
		Color:          nil,
		Avatar:         "https://example.com/path-to-avatar.jpg",
		AdminCanManage: false,
		Archived:       false,
		Members: []Member{
			{
				User: User{
					Id:             4567823,
					Username:       "Andrew Santos",
					Color:          "#098765",
					ProfilePicture: "https://example.com/path-to-picture.jpg",
					Initials:       "AS",
				},
			},
		},
		Statuses:          nil,
		MultipleAssignees: false,
		Features:          Features{},
	}

	expectedEntity := entities.Entity{
		Name:         space.Name,
		EntityUrl:    "",
		ExternalId:   space.Id,
		Type:         EntityTypeSpace,
		ContentUrl:   "",
		OwnerId:      ownerId,
		LastModified: time.Time{},
		Data:         space,
	}

	currentEntity := space.ToEntity(ownerId)
	assert.Equal(t, expectedEntity, currentEntity)
}

func TestClient_GetSpaces(t *testing.T) {
	testCases := []struct {
		name              string
		spaceId           string
		responseJson      string
		responseHttpCode  int
		expectedResult    []Space
		shouldReturnError bool
	}{
		{
			name:             "with correct response",
			responseHttpCode: http.StatusOK,
			responseJson:     correctSpacesJson,
			expectedResult: []Space{
				{
					Id:             "90150340588",
					Name:           "Team Space",
					Private:        false,
					Color:          "#03A2FD",
					Avatar:         "",
					AdminCanManage: true,
					Archived:       false,
					Members:        []Member(nil),
					Statuses: []Status{
						{Status: "to do", Type: "open", Orderindex: 0, Color: "#87909e"},
						{Status: "in progress", Type: "custom", Orderindex: 1, Color: "#1090e0"},
						{Status: "complete", Type: "closed", Orderindex: 2, Color: "#008844"},
					},
					MultipleAssignees: true,
					Features: Features{
						DueDates: DueDates{
							Enabled:            true,
							StartDate:          true,
							RemapDueDates:      false,
							RemapClosedDueDate: false,
						},
						TimeTracking:      EnabledFeature{Enabled: true},
						Tags:              EnabledFeature{Enabled: true},
						TimeEstimates:     EnabledFeature{Enabled: false},
						Checklists:        EnabledFeature{Enabled: false},
						CustomFields:      EnabledFeature{Enabled: true},
						RemapDependencies: EnabledFeature{Enabled: false},
						DependencyWarning: EnabledFeature{Enabled: true},
						Portfolios:        EnabledFeature{Enabled: false},
					},
				},
			},
		},
		{
			name:              "with incorrect json",
			responseHttpCode:  http.StatusOK,
			responseJson:      "[[[[[[[[[[",
			shouldReturnError: true,
		},
		{
			name:              "with token invalid response",
			responseHttpCode:  http.StatusUnauthorized,
			responseJson:      tokenInvalidJson,
			shouldReturnError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clickupMockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, apiKey, r.Header.Get("Authorization"))

				w.WriteHeader(tc.responseHttpCode)
				w.Write([]byte(tc.responseJson))
			}))
			defer clickupMockServer.Close()

			clickupClient := NewClient(clickupMockServer.URL+"/", ownerId, apiKey)
			currentResult, err := clickupClient.GetSpaces(context.Background(), tc.spaceId)
			if tc.shouldReturnError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, currentResult)
		})
	}
}
