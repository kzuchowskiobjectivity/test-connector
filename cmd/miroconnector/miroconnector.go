package main

import (
	"log"
	"miroconnector/internal/config"
	"miroconnector/pkg/connector"
	"miroconnector/pkg/miro"
	"miroconnector/pkg/sink"
	"sync"
)

func main() {
	configs, err := config.Load()
	if err != nil {
		log.Fatal("error when loading config file", err)
	}

	connectors := make([]connector.Connector, len(configs))
	for i, config := range configs {
		client := miro.NewClient(config.ApiUrl, config.ClientId, config.ApiToken)
		sink := sink.New(uint64(config.BufferSize))
		connector := connector.NewConnector(client, sink, config.ClientId)
		connectors[i] = *connector
	}

	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(len(connectors))
	for _, connector := range connectors {
		defer wg.Done()
		connector.QuickSync()
	}
}
