package main

import (
	"connectors/pkg/clickup"
	"connectors/pkg/entities"
	"connectors/pkg/idstorage"

	"context"
	"log"
	"time"
)

func main() {
	ctx := context.Background()

	userIds, err := idstorage.Load(idstorage.FromFile("path-to-config-file"))
	if err != nil {
		log.Fatal(err)
	}

	var connectors []*clickup.Client
	for _, id := range userIds {
		connectors = append(connectors, clickup.NewClient(clickup.ApiUrl, id.Owner, id.ApiKey))
	}

	entityChannel := make(chan entities.Entity, 64)
	defer close(entityChannel)
	for _, connector := range connectors {
		connector.GetEntities(ctx, entityChannel)
	}

	ticker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ticker.C:
			for _, connector := range connectors {
				go connector.GetEntities(ctx, entityChannel)
			}
		case e := <-entityChannel:
			//Placeholder for storage
			log.Printf("GOT ENTITY %v %v %v\n", e.Type, e.ExternalId, e.Name)
		}
	}
}
