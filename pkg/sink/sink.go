package sink

import (
	"encoding/json"
	"fmt"
	"miroconnector/pkg/entities"
	"os"
)

func New(bufferSize uint64) *Sink {
	s := &Sink{
		clientChan:    make(chan entities.Entity, bufferSize),
		readingIsDone: make(chan struct{}),
		allEntities:   make([]entities.Entity, 0, bufferSize),
	}

	go func() {
		defer close(s.readingIsDone)
		for e := range s.clientChan {
			e := e
			s.add(e)
		}
	}()

	return s
}

type Sink struct {
	clientChan    chan entities.Entity
	readingIsDone chan struct{}

	allEntities []entities.Entity
}

func (s *Sink) Close() {
	close(s.clientChan)
	<-s.readingIsDone
}

func (s *Sink) Push(e entities.Entity) {
	s.clientChan <- e
}

// consider using io.ReadCloser or function returning io.ReadCloser, error
// add proper test
// similar concept to https://github.com/gciezkowskiobjectivity/connectors/blob/main/pkg/idstorage/idstorage.go
func (s *Sink) Dump(ownerId string) ([]entities.Entity, error) {
	<-s.readingIsDone
	fileBytes, err := json.Marshal(s.allEntities)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entities: %w", err)
	}

	// close file/io.ReadCloser in defer
	file, err := os.Create(fmt.Sprintf("%s.json", ownerId))
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	_, err = file.Write(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return s.allEntities, nil
}

func (s *Sink) add(e entities.Entity) {

	s.allEntities = append(s.allEntities, e)
}
