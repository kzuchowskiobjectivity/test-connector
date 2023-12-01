package connector

import (
	"log"
	"miroconnector/pkg/miro"
	"miroconnector/pkg/sink"
	"sync"
)

type Connector struct {
	miroClient *miro.Client
	sink       *sink.Sink
	userId     string
}

func NewConnector(client *miro.Client, sink *sink.Sink, userId string) *Connector {
	return &Connector{miroClient: client, sink: sink, userId: userId}
}

func (c *Connector) QuickSync() {
	defer c.sink.Dump(c.userId)
	defer c.sink.Close()
	c.syncData()
}

func (c *Connector) syncData() {
	boards := c.miroClient.GetBoards()
	boardItems := make(chan miro.BoardItem, 50)

	go func() {
		var wg sync.WaitGroup
		defer close(boardItems)
		defer wg.Wait()
		for boardResult := range boards {
			if boardResult.Error != nil {
				log.Print(boardResult.Error)
			}
			board := boardResult.Board
			c.sink.Push(board.ToEntity(c.userId))
			wg.Add(1)
			go func() {
				c.streamBoardItems(board.ID, boardItems, &wg)
			}()
		}
	}()

	for boardItem := range boardItems {
		c.sink.Push(boardItem.ToEntity(c.userId))
	}
}

func (c *Connector) streamBoardItems(boardId string, items chan<- miro.BoardItem, wg *sync.WaitGroup) {
	boardItemResults := c.miroClient.GetBoardItems(boardId)
	for itemResult := range boardItemResults {
		if itemResult.Error != nil {
			log.Print(itemResult.Error)
		}
		items <- itemResult.Result
	}
	wg.Done()
}
