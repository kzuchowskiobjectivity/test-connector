package clickup

import (
	"connectors/pkg/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestList_ToEntity(t *testing.T) {
	list := List{
		Id:         "4567865678",
		Name:       "Lista pierwsza",
		Orderindex: 1,
		Content:    "I am a content",
		Status: ListStatus{
			Status:    "done",
			Color:     "#00FF00",
			HideLabel: false,
		},
		Priority: ListPriority{
			Priority: "High",
			Color:    "#FF0000",
		},
		Assignee:  nil,
		TaskCount: nil,
		DueDate:   "",
		StartDate: "qweqwrqweq",
		Folder: ListFolder{
			Id:     "376823",
			Name:   "MyBasicFolder",
			Hidden: false,
			Access: true,
		},
		Space: ListSpace{
			Id:     "241151",
			Name:   "IAmSpace",
			Access: true,
		},
		Archived:         false,
		OverrideStatuses: true,
		PermissionLevel:  "",
	}

	expectedEntity := entities.Entity{
		Name:         list.Name,
		EntityUrl:    "",
		ExternalId:   list.Id,
		Type:         EntityTypeList,
		ContentUrl:   "",
		OwnerId:      ownerId,
		LastModified: time.Time{},
		Data:         list,
	}

	currentEntity := list.ToEntity(ownerId)
	assert.Equal(t, expectedEntity, currentEntity)
}
