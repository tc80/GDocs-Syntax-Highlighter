package runner

import (
	"context"
	"log"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

var (
	client = newClient()
)

// Gets a new Docker client, exiting on error.
func newClient() *docker.Client {
	c, err := docker.NewEnvClient()
	if err != nil {
		log.Fatalf("Failed to initialize Docker client: %v\n", err)
	}
	return c
}

// Check if Docker has an image by checking if it
// contains a particular repo tag.
func hasImage(repoTag string) bool {
	images, err := client.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Fatalf("Failed to fetch Docker images: %v\n", err)
	}
	for _, img := range images {
		for _, t := range img.RepoTags {
			if t == repoTag {
				return true // found repo tag
			}
		}
	}
	return false
}

func runContainer() {
	client.ContainerStart(context.Background(), "", types.ContainerStartOptions{})
}
