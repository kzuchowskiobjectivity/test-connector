package sink_test

import (
	"connectors/pkg/entities"
	"connectors/pkg/sink"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSink(t *testing.T) {
	testCases := []struct {
		name       string
		dataToPush []entities.Entity
	}{
		{
			name:       "no data",
			dataToPush: []entities.Entity{},
		},
		{
			name:       "single data",
			dataToPush: makePrefixedEntities([]string{"single"}),
		},
		{
			name:       "multiple data",
			dataToPush: makePrefixedEntities([]string{"first", "second", "third"}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := sink.New(1)

			for _, e := range tc.dataToPush {
				s.Push(e)
			}
			s.Close()

			allEntities := s.Dump()
			assert.Equal(t, len(tc.dataToPush), len(allEntities))
			assert.Equal(t, tc.dataToPush, allEntities)
		})
	}

}

func makePrefixedEntities(prefixes []string) []entities.Entity {

	result := make([]entities.Entity, 0, len(prefixes))
	for _, p := range prefixes {
		e := entities.Entity{
			Name:       p + "_Name",
			EntityUrl:  p + "_EntityUrl",
			ExternalId: p + "_ExternalId",
			Type:       p + "_Type",
			ContentUrl: p + "_ContentUrl",
			OwnerId:    p + "_OwnerId",
		}
		result = append(result, e)
	}

	return result
}
