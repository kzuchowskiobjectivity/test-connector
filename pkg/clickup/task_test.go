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

const correctTaskJson = "{\"tasks\":[{\"id\":\"9hx\",\"custom_item_id\":null,\"name\":\"New Task Name\",\"status\":{\"status\":\"Open\",\"color\":\"#d3d3d3\",\"orderindex\":0,\"type\":\"open\"},\"markdown_description\":\"Task description\",\"orderindex\":\"1.00000000000000000000000000000000\",\"date_created\":\"1567780450202\",\"date_updated\":\"1567780450202\",\"date_closed\":null,\"date_done\":null,\"creator\":{\"id\":183,\"username\":\"John Doe\",\"color\":\"#827718\",\"profilePicture\":\"https://attachments-public.clickup.com/profilePictures/183_abc.jpg\"},\"assignees\":[],\"checklists\":[],\"tags\":[],\"parent\":null,\"priority\":null,\"due_date\":null,\"start_date\":null,\"time_estimate\":null,\"time_spent\":null,\"list\":{\"id\":\"123\"},\"folder\":{\"id\":\"456\"},\"space\":{\"id\":\"789\"},\"url\":\"https://app.clickup.com/t/9hx\"},{\"id\":\"9hz\",\"custom_item_id\":null,\"name\":\"Second task\",\"status\":{\"status\":\"Open\",\"color\":\"#d3d3d3\",\"orderindex\":0,\"type\":\"open\"},\"orderindex\":\"2.00000000000000000000000000000000\",\"date_created\":\"1567780450202\",\"date_updated\":\"1567780450202\",\"date_closed\":null,\"date_done\":null,\"creator\":{\"id\":183,\"username\":\"John Doe\",\"color\":\"#827718\",\"profilePicture\":\"https://attachments-public.clickup.com/profilePictures/183_abc.jpg\"},\"assignees\":[],\"checklists\":[],\"tags\":[],\"parent\":null,\"priority\":null,\"due_date\":null,\"start_date\":null,\"time_estimate\":null,\"time_spent\":null,\"list\":{\"id\":\"123\"},\"folder\":{\"id\":\"456\"},\"space\":{\"id\":\"789\"},\"url\":\"https://app.clickup.com/t/9hz\"}]}"

func TestTask_ToEntity(t *testing.T) {
	task := Task{
		Id:   "AaaaaA",
		Name: "ZzzzZZzzz",
		Status: Status{
			Status:     "done",
			Type:       "done",
			Orderindex: 0,
			Color:      "#22FF22",
		},
		Creator: User{
			Id:             143253215,
			Username:       "Andrew Santos",
			Color:          "#988765",
			ProfilePicture: "https://path-to-pic/",
			Initials:       "AS",
		},
		List:   StringId{Id: "ewrweiuwe"},
		Folder: StringId{Id: "wqf3r98235"},
		Space:  StringId{Id: "fqwrqrqr"},
		Url:    "https://url-to-task/",
	}

	expectedEntity := entities.Entity{
		Name:         task.Name,
		EntityUrl:    "",
		ExternalId:   task.Id,
		Type:         EntityTypeTask,
		ContentUrl:   task.Url,
		OwnerId:      ownerId,
		LastModified: time.Time{},
		Data:         task,
	}

	currentEntity := task.ToEntity(ownerId)
	assert.Equal(t, expectedEntity, currentEntity)
}

func TestClient_GetTasks(t *testing.T) {
	testCases := []struct {
		name              string
		responseJson      string
		responseHttpCode  int
		listId            string
		expectedResult    []Task
		shouldReturnError bool
	}{
		{
			name:             "with correct response",
			responseHttpCode: http.StatusOK,
			responseJson:     correctTaskJson,
			expectedResult: []Task{
				{
					Id:   "9hx",
					Name: "New Task Name",
					Status: Status{
						Status:     "Open",
						Type:       "open",
						Orderindex: 0,
						Color:      "#d3d3d3",
					},
					MarkdownDescription: "Task description",
					Orderindex:          "1.00000000000000000000000000000000",
					DateCreated:         "1567780450202",
					DateUpdated:         "1567780450202",
					Creator: User{
						Id:             183,
						Username:       "John Doe",
						Color:          "#827718",
						ProfilePicture: "https://attachments-public.clickup.com/profilePictures/183_abc.jpg",
						Initials:       "",
					},
					Assignees:  []interface{}{},
					Checklists: []interface{}{},
					Tags:       []interface{}{},
					List:       StringId{Id: "123"},
					Folder:     StringId{Id: "456"},
					Space:      StringId{Id: "789"},
					Url:        "https://app.clickup.com/t/9hx",
				},
				{
					Id:   "9hz",
					Name: "Second task",
					Status: Status{
						Status:     "Open",
						Type:       "open",
						Orderindex: 0,
						Color:      "#d3d3d3",
					},
					MarkdownDescription: "",
					Orderindex:          "2.00000000000000000000000000000000",
					DateCreated:         "1567780450202",
					DateUpdated:         "1567780450202",
					Creator: User{
						Id:             183,
						Username:       "John Doe",
						Color:          "#827718",
						ProfilePicture: "https://attachments-public.clickup.com/profilePictures/183_abc.jpg",
						Initials:       "",
					},
					Assignees:  []interface{}{},
					Checklists: []interface{}{},
					Tags:       []interface{}{},
					List:       StringId{Id: "123"},
					Folder:     StringId{Id: "456"},
					Space:      StringId{Id: "789"},
					Url:        "https://app.clickup.com/t/9hz",
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
			currentResult, err := clickupClient.GetTasks(context.Background(), tc.listId)

			if tc.shouldReturnError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, currentResult)
		})
	}
}
