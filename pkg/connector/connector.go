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
	boards := c.miroClient.GetBoards()
	boardItems := make(chan miro.BoardItem, 50)
	done := make(chan bool)

	// move to function
	// sync boards
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

			// trigger board items sync
			wg.Add(1)
			go func() {
				defer wg.Done()
				c.streamBoardItems(board.ID, boardItems)
			}()
		}
	}()

	// move to function
	// sync board items
	go func() {
		for boardItem := range boardItems {
			c.sink.Push(boardItem.ToEntity(c.userId))
		}
		done <- true
	}()

	<-done
}

func (c *Connector) streamBoardItems(boardId string, items chan<- miro.BoardItem) {
	boardItemResults := c.miroClient.GetBoardItems(boardId)
	for itemResult := range boardItemResults {
		if itemResult.Error != nil {
			log.Print(itemResult.Error)
		}
		items <- itemResult.Result
	}
}
