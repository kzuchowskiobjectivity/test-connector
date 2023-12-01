package miro

import (
	"encoding/json"
	"fmt"
	"miroconnector/pkg/entities"
	"net/http"
	"time"
)

type BoardItemResult struct {
	Result BoardItem
	Error  error
}

func (c *Client) GetBoardItems(boardId string) <-chan BoardItemResult {
	size := 50
	ch := make(chan BoardItemResult, size)
	go func() {
		defer close(ch)
		var cursor string
		for {
			itemsResponse, err := c.getBoardItems(boardId, cursor)
			if err != nil {
				ch <- BoardItemResult{Error: err}
				return
			}
			for _, item := range itemsResponse.Data {
				ch <- BoardItemResult{Result: item}
			}
			if len(itemsResponse.Cursor) == 0 {
				return
			}
			cursor = itemsResponse.Cursor
		}
	}()
	return ch
}

func (c *Client) getBoardItems(boardId string, cursor string) (BoardItemsResponse, error) {
	url := fmt.Sprintf("%s/boards/%s/items?cursor=%s", c.url, boardId, cursor)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	if err != nil {
		return BoardItemsResponse{}, fmt.Errorf("failed to create http request: %w", err)
	}
	response, err := c.doRequest(req)
	if err != nil {
		return BoardItemsResponse{}, fmt.Errorf("failed to perform http request: %w", err)
	}
	defer response.Body.Close()
	var itemsResponse BoardItemsResponse
	err = json.NewDecoder(response.Body).Decode(&itemsResponse)
	if err != nil {
		return BoardItemsResponse{}, fmt.Errorf("error when decoding json: %w", err)
	}
	return itemsResponse, nil
}

func (i BoardItem) ToEntity(ownerId string) entities.Entity {
	item := entities.Entity{
		Name:         i.Type,
		ExternalId:   i.ID,
		Type:         "boardItem",
		OwnerId:      ownerId,
		LastModified: i.ModifiedAt,
		Data:         i,
	}
	return item
}

type BoardItem struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data struct {
		Content string `json:"content"`
		Shape   string `json:"shape"`
	} `json:"data"`
	Style struct {
		FillColor         string `json:"fillColor"`
		FillOpacity       string `json:"fillOpacity"`
		FontFamily        string `json:"fontFamily"`
		FontSize          string `json:"fontSize"`
		BorderColor       string `json:"borderColor"`
		BorderWidth       string `json:"borderWidth"`
		BorderOpacity     string `json:"borderOpacity"`
		BorderStyle       string `json:"borderStyle"`
		TextAlign         string `json:"textAlign"`
		TextAlignVertical string `json:"textAlignVertical"`
		Color             string `json:"color"`
	} `json:"style"`
	Geometry struct {
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	} `json:"geometry"`
	Position struct {
		X          float64 `json:"x"`
		Y          float64 `json:"y"`
		Origin     string  `json:"origin"`
		RelativeTo string  `json:"relativeTo"`
	} `json:"position"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"createdBy"`
	ModifiedAt time.Time `json:"modifiedAt"`
	ModifiedBy struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"modifiedBy"`
}

type BoardItemsResponse struct {
	Size   int         `json:"size"`
	Limit  int         `json:"limit"`
	Total  int         `json:"total"`
	Cursor string      `json:"cursor"`
	Data   []BoardItem `json:"data"`
	Links  struct {
		Self string `json:"self"`
	} `json:"links"`
	Type string `json:"type"`
}
