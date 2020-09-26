package runner

import (
	"log"
)

const (
	playgroundImage = "golang/playground"
	playgroundTag   = "latest"
)

func init() {
	repoTag := playgroundImage + ":" + playgroundTag
	if !hasImage(repoTag) {
		log.Fatalf("Missing Docker image: %s\n", repoTag)
	}
}
