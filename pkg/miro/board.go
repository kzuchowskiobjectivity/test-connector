package miro

import (
	"encoding/json"
	"fmt"
	"miroconnector/pkg/entities"
	"net/http"
	"time"
)

type BoardResult struct {
	Board Board
	Error error
}

// add tests
func (c *Client) GetBoards() <-chan BoardResult {
	size := 50
	ch := make(chan BoardResult, size)
	go func() {
		defer close(ch)
		for i := 0; ; i = i + 1 {
			boards, err := c.getBoards(i*size, size)
			if err != nil {
				ch <- BoardResult{Error: fmt.Errorf("error when getting boards: %w", err)}
				return
			}
			if len(boards) == 0 {
				return
			}
			for _, board := range boards {
				ch <- BoardResult{Board: board}
			}
		}
	}()
	return ch
}

func (c *Client) getBoards(offset int, limit int) ([]Board, error) {
	url := fmt.Sprintf("%s/boards?limit=%d&offset=%d", c.url, limit, offset)

	// should you add to header before error check?
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	response, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform http request: %w", err)
	}

	defer response.Body.Close()
	var boards BoardsResponse
	err = json.NewDecoder(response.Body).Decode(&boards)
	if err != nil {
		return nil, fmt.Errorf("error when decoding json: %w", err)
	}

	return boards.Data, nil
}

func (b Board) ToEntity(ownerId string) entities.Entity {
	entity := entities.Entity{
		Name:         b.Name,
		EntityUrl:    b.ViewLink,
		ExternalId:   b.ID,
		Type:         "board",
		OwnerId:      ownerId,
		Data:         b,
		ContentUrl:   "",
		LastModified: b.ModifiedAt,
	}
	return entity
}

type BoardsResponse struct {
	Size   int     `json:"size"`
	Offset int     `json:"offset"`
	Limit  int     `json:"limit"`
	Total  int     `json:"total"`
	Data   []Board `json:"data"`
}

type Board struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Links       struct {
		Self    string `json:"self"`
		Related string `json:"related"`
	} `json:"links"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"createdBy"`
	CurrentUserMembership Membership `json:"currentUserMembership"`
	LastOpenedAt          time.Time  `json:"lastOpenedAt"`
	LastOpenedBy          struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"lastOpenedBy"`
	ModifiedAt time.Time `json:"modifiedAt"`
	ModifiedBy struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"modifiedBy"`
	Owner struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"owner"`
	PermissionsPolicy struct {
		CollaborationToolsStartAccess string `json:"collaborationToolsStartAccess"`
		CopyAccess                    string `json:"copyAccess"`
		CopyAccessLevel               string `json:"copyAccessLevel"`
		SharingAccess                 string `json:"sharingAccess"`
	} `json:"permissionsPolicy"`
	Picture struct {
		ID       int64  `json:"id"`
		Type     string `json:"type"`
		ImageURL string `json:"imageURL"`
	} `json:"picture"`
	Policy struct {
		PermissionsPolicy struct {
			CollaborationToolsStartAccess string `json:"collaborationToolsStartAccess"`
			CopyAccess                    string `json:"copyAccess"`
			CopyAccessLevel               string `json:"copyAccessLevel"`
			SharingAccess                 string `json:"sharingAccess"`
		} `json:"permissionsPolicy"`
		SharingPolicy struct {
			Access                            string `json:"access"`
			InviteToAccountAndBoardLinkAccess string `json:"inviteToAccountAndBoardLinkAccess"`
			OrganizationAccess                string `json:"organizationAccess"`
			TeamAccess                        string `json:"teamAccess"`
		} `json:"sharingPolicy"`
	} `json:"policy"`
	SharingPolicy struct {
		Access                            string `json:"access"`
		InviteToAccountAndBoardLinkAccess string `json:"inviteToAccountAndBoardLinkAccess"`
		OrganizationAccess                string `json:"organizationAccess"`
		TeamAccess                        string `json:"teamAccess"`
	} `json:"sharingPolicy"`
	Team struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"team"`
	ViewLink string `json:"viewLink"`
}

type Membership struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}
