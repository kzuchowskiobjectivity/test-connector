package clickup

import (
	"connectors/pkg/entities"

	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const ApiUrl = "https://api.clickup.com/api/v2/"

type Client struct {
	client  *http.Client
	url     string
	ownerId string
	apiKey  string
}

func NewClient(url string, ownerId, apiKey string) *Client {
	return &Client{
		client:  http.DefaultClient,
		url:     url,
		ownerId: ownerId,
		apiKey:  apiKey,
	}
}

func (c *Client) GetEntities(ctx context.Context, entityChannel chan entities.Entity) {
	log.Println("Getting workspaces")
	workspaces, err := c.GetWorkspaces(ctx)
	if err != nil {
		log.Print(err)
		return
	}

	log.Println("Getting spaces")
	spaces := c.getSpacesForWorkspaces(ctx, workspaces, entityChannel)

	log.Println("Getting folders")
	folders := c.getFoldersForSpaces(ctx, spaces, entityChannel)

	log.Println("Getting lists")
	lists := c.getListsForFolders(folders, entityChannel)

	log.Println("Getting tasks")
	_ = c.getTasksForLists(ctx, lists, entityChannel)
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do http request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		responseBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		return nil, fmt.Errorf("failed to retrieve workspaces: %s", string(responseBytes))
	}

	return resp, nil
}

func (c *Client) getSpacesForWorkspaces(ctx context.Context, workspaces []Team, entityChannel chan entities.Entity) []Space {
	wg := sync.WaitGroup{}
	wg.Add(len(workspaces))

	var spaces []Space
	for _, w := range workspaces {
		w := w
		go func() {
			defer wg.Done()
			entityChannel <- w.ToEntity(c.ownerId)
			newSpaces, err := c.GetSpaces(ctx, w.Id)
			if err != nil {
				log.Println(err)
				return
			}
			spaces = append(spaces, newSpaces...)
		}()
	}
	wg.Wait()

	return spaces
}

func (c *Client) getFoldersForSpaces(ctx context.Context, spaces []Space, entityChannel chan entities.Entity) []Folder {
	wg := sync.WaitGroup{}
	wg.Add(len(spaces))
	var folders []Folder
	for _, s := range spaces {
		s := s
		go func() {
			defer wg.Done()
			entityChannel <- s.ToEntity(c.ownerId)
			newFolders, err := c.GetFolders(ctx, s.Id)
			if err != nil {
				log.Println(err)
				return
			}
			folders = append(folders, newFolders...)
		}()
	}
	wg.Wait()
	return folders
}

func (c *Client) getListsForFolders(folders []Folder, entityChannel chan entities.Entity) []List {
	wg := sync.WaitGroup{}
	wg.Add(len(folders))
	var lists []List
	for _, f := range folders {
		f := f
		go func() {
			defer wg.Done()
			entityChannel <- f.ToEntity(c.ownerId)
			lists = append(lists, f.Lists...)
		}()
	}
	wg.Wait()
	return lists
}

func (c *Client) getTasksForLists(ctx context.Context, lists []List, entityChannel chan entities.Entity) []Task {
	wg := sync.WaitGroup{}
	wg.Add(len(lists))
	var tasks []Task
	for _, l := range lists {
		l := l
		go func() {
			defer wg.Done()
			entityChannel <- l.ToEntity(c.ownerId)
			newTasks, err := c.GetTasks(ctx, l.Id)
			if err != nil {
				log.Println(err)
				return
			}
			tasks = append(tasks, newTasks...)
			for _, nt := range newTasks {
				entityChannel <- nt.ToEntity(c.ownerId)
			}
		}()
	}
	wg.Wait()
	return tasks
}
