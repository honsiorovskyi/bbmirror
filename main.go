package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
)

type UpdateMessage struct {
	Repository struct {
		FullName string `json:"full_name"`
		SCM      string `json:"scm"`
	}
}

func updateRepository(w http.ResponseWriter, r *http.Request) {
	// read payload
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading payload")
		return
	}

	// decode payload
	message := UpdateMessage{}
	if err := json.Unmarshal(body, &message); err != nil {
		log.Println("Error parsing JSON")
		return
	}

	// check for supported SCMs
	if message.Repository.SCM != "git" {
		log.Printf("SCM %q is not supported yet", message.Repository.SCM)
	}

	// check if directory exists
	basePath := config("REPOSITORY_PATH", "/var/lib/bbmirror/repository")
	repoPath := path.Join(basePath, message.Repository.FullName)

	if _, err := os.Stat(repoPath); err != nil {
		if os.IsNotExist(err) {
			// git clone
			log.Printf("Clonning %q...", message.Repository.FullName)
			cmd := exec.Command(
				"git", "clone", "--bare",
				fmt.Sprintf("git@bitbucket.org:%s.git", message.Repository.FullName),
				repoPath,
			)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			err := cmd.Run()

			if err != nil {
				log.Println("Clonning %q failed:", err)
				return
			}

			log.Printf("Successfully cloned %q.", message.Repository.FullName)
			return
		} else {
			log.Printf("Error accessing directory %q: %v", repoPath, err)
			return
		}
	}

	// git pull
	log.Printf("Fetching %q...", message.Repository.FullName)
	cmd := exec.Command("git", "fetch", "--all")
	cmd.Dir = repoPath
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Println("Fetching %q failed:", err)
		return
	}

	log.Printf("Successfully fetched %q.", message.Repository.FullName)
}

func config(name string, default_value string) string {
	val := os.Getenv(name)
	if val == "" {
		val = default_value
	}

	return val
}

func main() {
	http.HandleFunc("/update", updateRepository)
	err := http.ListenAndServe(config("LISTEN", "127.0.0.1:5678"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
