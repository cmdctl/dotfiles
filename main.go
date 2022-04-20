package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Version      string   `yaml:"version"`
	TrackedFiles []string `yaml:"include"`
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	fmt.Println("Syncing dotfiles changes to dotfiles repository...")
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
	fmt.Println("Done!")
}

func commitChanges(repo *git.Repository) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	_, _ = worktree.Add(".")
	_, _ = worktree.Commit("changes to dotfiles", &git.CommitOptions{})
	remote, err := repo.Remote("origin")
	if err == nil {
		remote.Push(&git.PushOptions{})
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
	fmt.Printf("Opening %s repository...\n", repoName)
	repo, err := git.PlainOpen(repoName)
	return repo, err
}
