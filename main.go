package main

import (
	"fmt"
	ssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"io"
	"io/ioutil"
	"os"

	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v2"
)

// Config holds the configuration contained in .dotfiles.yml
type Config struct {
	Version      string   `yaml:"version"`
	TrackedFiles []string `yaml:"include"`
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	// fmt.Println("Syncing dotfiles changes to dotfiles repository...")
	repo, err := getDotfilesRepo(home)
	if err != nil {
		panic(err)
	}

	config := readConfig(home)
	for _, file := range config.TrackedFiles {
		_, err := copy(fmt.Sprintf("%s/%s", home, file), fmt.Sprintf("%s/%s", getRepoName(home), file))
		if err != nil {
			panic(err)
		}
		err = os.Chmod(fmt.Sprintf("%s/%s", getRepoName(home), file), 0600)
		if err != nil {
			panic(err)
		}
	}
	err = commitChanges(repo)
	if err != nil {
		panic(err)
	}
	// fmt.Println("Done!")
}

func commitChanges(repo *git.Repository) error {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	_, err = worktree.Add(".")
	if err != nil {
		return fmt.Errorf("could not add files to git %s", err.Error())
	}
	_, err = worktree.Commit("changes to dotfiles", &git.CommitOptions{})
	if err != nil {
		return fmt.Errorf("could not commit files to git %s", err.Error())
	}
	auth, err := publicKey(home)
	if err != nil {
		panic(err)
	}
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil {
		fmt.Printf("could not push changes to remote: %s", err)
	}
	return nil
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func readConfig(home string) Config {
	var config Config
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/.dotfiles.yml", home))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func getRepoName(home string) string {
	return fmt.Sprintf("%s/.dotfiles", home)
}
func getDotfilesRepo(home string) (*git.Repository, error) {
	repoName := getRepoName(home)
	if _, err := os.Stat(fmt.Sprintf("%s/%s", repoName, ".git")); os.IsNotExist(err) {
		fmt.Printf("Repo does not exist, initializing a repository for tracking dotfiles at %s...\n", repoName)
		repo, err := git.PlainInit(repoName, false)
		if err != nil {
			return repo, err
		}
		return repo, nil
	}
	// fmt.Printf("Opening %s repository...\n", repoName)
	repo, err := git.PlainOpen(repoName)
	return repo, err
}

func publicKey(home string) (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys
	sshPath := home + "/.ssh/id_ed25519"
	sshKey, _ := ioutil.ReadFile(sshPath)
	publicKey, err := ssh.NewPublicKeys("git", []byte(sshKey), "")
	if err != nil {
		return nil, err
	}
	return publicKey, err
}
