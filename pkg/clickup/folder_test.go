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

const correctFoldersJson = "{\n  \"folders\": [\n    {\n      \"id\": \"457\",\n      \"name\": \"Updated Folder Name\",\n      \"orderindex\": 0,\n      \"override_statuses\": false,\n      \"hidden\": false,\n      \"space\": {\n        \"id\": \"789\",\n        \"name\": \"Space Name\",\n        \"access\": true\n      },\n      \"task_count\": \"0\",\n      \"lists\": []\n    },\n    {\n      \"id\": \"321\",\n      \"name\": \"Folder Name\",\n      \"orderindex\": 0,\n      \"override_statuses\": false,\n      \"hidden\": false,\n      \"space\": {\n        \"id\": \"789\",\n        \"name\": \"Space Name\",\n        \"access\": true\n      },\n      \"task_count\": \"0\",\n      \"lists\": []\n    }\n  ]\n}"

func TestFolder_ToEntity(t *testing.T) {
	folder := Folder{
		Id:               "FolderId",
		Name:             "Folder name 1",
		OrderIndex:       4,
		OverrideStatuses: true,
		Hidden:           false,
		Space: FolderSpace{
			Id:     "32423523",
			Name:   "Universe",
			Access: true,
		},
		TaskCount: "3",
		Lists: []List{
			{
				Id:         "ListInFolder",
				Name:       "ListName1",
				Orderindex: 6,
				Content:    "This is a content",
				Status: ListStatus{
					Status:    "useful",
					Color:     "#456789",
					HideLabel: false,
				},
				Priority: ListPriority{
					Priority: "High",
					Color:    "#00FF00",
				},
				DueDate: "",
				Folder: ListFolder{
					Id:     "List Folder",
					Name:   "Folder name",
					Hidden: true,
					Access: true,
				},
				Space: ListSpace{
					Id:     "Universe",
					Name:   "List in universe",
					Access: true,
				},
				PermissionLevel: "access",
			},
		},
	}

	expectedEntity := entities.Entity{
		Name:         folder.Name,
		EntityUrl:    "",
		ExternalId:   folder.Id,
		Type:         EntityTypeFolder,
		ContentUrl:   "",
		OwnerId:      ownerId,
		LastModified: time.Time{},
		Data:         folder,
	}

	currentEntity := folder.ToEntity(ownerId)

	assert.Equal(t, expectedEntity, currentEntity)
}

func TestClient_GetFolders(t *testing.T) {
	testCases := []struct {
		name              string
		spaceId           string
		responseJson      string
		responseHttpCode  int
		expectedResult    []Folder
		shouldReturnError bool
	}{
		{
			name:             "with correct response",
			responseHttpCode: http.StatusOK,
			responseJson:     correctFoldersJson,
			expectedResult: []Folder{
				{
					Id:               "457",
					Name:             "Updated Folder Name",
					OrderIndex:       0,
					OverrideStatuses: false,
					Hidden:           false,
					Space: FolderSpace{
						Id:     "789",
						Name:   "Space Name",
						Access: true,
					},
					TaskCount: "0",
					Lists:     []List{},
				},
				{
					Id:               "321",
					Name:             "Folder Name",
					OrderIndex:       0,
					OverrideStatuses: false,
					Hidden:           false,
					Space: FolderSpace{
						Id:     "789",
						Name:   "Space Name",
						Access: true,
					},
					TaskCount: "0",
					Lists:     []List{},
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
			currentResult, err := clickupClient.GetFolders(context.Background(), tc.spaceId)
			if tc.shouldReturnError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, currentResult)
		})
	}
}
