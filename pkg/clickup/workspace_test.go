package clickup

import (
	"connectors/pkg/entities"

	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	ownerId = "1ce19795-4563-41e8-a745-f602e0ec6f65"
	apiKey  = "HERE_I_HAVE_YOUR_API_KEY"

	correctWorkspaceJson = "{\"teams\":[{\"id\":\"1234\",\"name\":\"My ClickUp Workspace\",\"color\":\"#000000\",\"avatar\":\"https://clickup.com/avatar.jpg\",\"members\":[{\"user\":{\"id\":123,\"username\":\"John Doe\",\"color\":\"#000000\",\"profilePicture\":\"https://clickup.com/avatar.jpg\"}}]}]}"
	tokenInvalidJson     = "{\"err\":\"Token invalid\",\"ECODE\":\"OAUTH_025\"}"
)

func TestTeam_ToEntity(t *testing.T) {
	testCases := []struct {
		name           string
		team           Team
		expectedEntity entities.Entity
	}{
		{
			name: "without members",
			team: Team{
				Id:      "54529837431",
				Name:    "W1Name",
				Color:   "#AABBCC",
				Avatar:  "https://example.com/url-to-avatar.jpg",
				Members: nil,
				Url:     "",
			},
			expectedEntity: entities.Entity{
				Name:         "W1Name",
				EntityUrl:    "",
				ExternalId:   "54529837431",
				Type:         EntityTypeWorkspace,
				ContentUrl:   "",
				OwnerId:      ownerId,
				LastModified: time.Time{},
				Data: Team{
					Id:      "54529837431",
					Name:    "W1Name",
					Color:   "#AABBCC",
					Avatar:  "https://example.com/url-to-avatar.jpg",
					Members: nil,
					Url:     "",
				},
			},
		},
		{
			name: "with members",
			team: Team{
				Id:     "54529837431",
				Name:   "W1Name",
				Color:  "#AABBCC",
				Avatar: "https://example.com/url-to-avatar.jpg",
				Members: []Member{
					{
						User: User{
							Id:             2142141532,
							Username:       "Xantos Santos",
							Color:          "#123456",
							ProfilePicture: "https://example.com/url-to-picture.jpg",
							Initials:       "XS",
						},
					},
				},
				Url: "",
			},
			expectedEntity: entities.Entity{
				Name:         "W1Name",
				EntityUrl:    "",
				ExternalId:   "54529837431",
				Type:         EntityTypeWorkspace,
				ContentUrl:   "",
				OwnerId:      ownerId,
				LastModified: time.Time{},
				Data: Team{
					Id:     "54529837431",
					Name:   "W1Name",
					Color:  "#AABBCC",
					Avatar: "https://example.com/url-to-avatar.jpg",
					Members: []Member{
						{
							User: User{
								Id:             2142141532,
								Username:       "Xantos Santos",
								Color:          "#123456",
								ProfilePicture: "https://example.com/url-to-picture.jpg",
								Initials:       "XS",
							},
						},
					},
					Url: "",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			currentEntity := tc.team.ToEntity(ownerId)
			assert.Equal(t, tc.expectedEntity, currentEntity)
		})
	}
}

func TestClient_GetWorkspaces(t *testing.T) {
	testCases := []struct {
		name              string
		responseJson      string
		responseHttpCode  int
		expectedResult    []Team
		shouldReturnError bool
	}{
		{
			name:             "with correct response",
			responseHttpCode: http.StatusOK,
			responseJson:     correctWorkspaceJson,
			expectedResult: []Team{
				{
					Id:     "1234",
					Name:   "My ClickUp Workspace",
					Color:  "#000000",
					Avatar: "https://clickup.com/avatar.jpg",
					Members: []Member{
						{User: User{
							Id:             123,
							Username:       "John Doe",
							Color:          "#000000",
							ProfilePicture: "https://clickup.com/avatar.jpg",
							Initials:       ""}},
					},
					Url: "",
				},
			},
		},
		{
			name:              "with incorrect json response",
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
			currentResult, err := clickupClient.GetWorkspaces(context.Background())

			if tc.shouldReturnError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, currentResult)
		})
	}
}
