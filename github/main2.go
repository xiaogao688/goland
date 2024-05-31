package main

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
)

func main() {
	// Clone the given repository to the given directory
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/emn178/js-md5.git",
		//Auth: &http.BasicAuth{
		//	Username: "your-username", // can be anything except an empty string
		//	Password: "your-token",    // GitHub personal access token
		//},
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatalf("Failed to clone repository: %s", err)
	}

	// Access the references
	refs, err := repo.References()
	if err != nil {
		log.Fatalf("Failed to get references: %s", err)
	}

	// Iterate over the references
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		fmt.Println(ref.String())
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to iterate over references: %s", err)
	}
}
