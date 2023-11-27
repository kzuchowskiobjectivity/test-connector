package sink

import "connectors/pkg/entities"

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

func (s *Sink) Dump() []entities.Entity {
	<-s.readingIsDone
	return s.allEntities
}

func (s *Sink) add(e entities.Entity) {

	s.allEntities = append(s.allEntities, e)
}
