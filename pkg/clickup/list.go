package clickup

import (
	"connectors/pkg/entities"
	"time"
)

const EntityTypeList = "list"

type Lists struct {
	Lists []List `json:"lists"`
}

type List struct {
	Id               string       `json:"id"`
	Name             string       `json:"name"`
	Orderindex       int          `json:"orderindex"`
	Content          string       `json:"content"`
	Status           ListStatus   `json:"status"`
	Priority         ListPriority `json:"priority"`
	Assignee         interface{}  `json:"assignee"`
	TaskCount        interface{}  `json:"task_count"`
	DueDate          string       `json:"due_date"`
	StartDate        interface{}  `json:"start_date"`
	Folder           ListFolder   `json:"folder"`
	Space            ListSpace    `json:"space"`
	Archived         bool         `json:"archived"`
	OverrideStatuses bool         `json:"override_statuses"`
	PermissionLevel  string       `json:"permission_level"`
}

type ListStatus struct {
	Status    string `json:"status"`
	Color     string `json:"color"`
	HideLabel bool   `json:"hide_label"`
}

type ListPriority struct {
	Priority string `json:"priority"`
	Color    string `json:"color"`
}

type ListFolder struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Hidden bool   `json:"hidden"`
	Access bool   `json:"access"`
}

type ListSpace struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Access bool   `json:"access"`
}

func (l List) ToEntity(ownerId string) entities.Entity {
	return entities.Entity{
		Name:         l.Name,
		EntityUrl:    "",
		ExternalId:   l.Id,
		Type:         EntityTypeList,
		ContentUrl:   "",
		OwnerId:      ownerId,
		LastModified: time.Time{},
		Data:         l,
	}
}
