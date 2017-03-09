package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func main() {
	topLevel()
}

func topLevel() {
	matched, err := filepath.Glob("src/*")
	check(err)
	repos := []repo{}
	for _, f := range matched {
		if isDir(f) {
			repo := newRepoForPath(f)
			repos = append(repos, repo)
		}
	}
	sort.Sort(reposByFullSize(repos))
	for _, repo := range repos {
		fmt.Println(repo.asString())
	}
}

func newRepoForPath(dirPath string) repo {
	fullSize := DirSizeMB(dirPath)
	gitSize := DirSizeMB(dirPath + "/.git")
	codeSize := fullSize - gitSize
	return repo{
		name:     dirPath,
		fullSize: fullSize,
		gitSize:  gitSize,
		codeSize: codeSize,
	}
}

/*
  REPO logic
*/
type repo struct {
	name     string
	fullSize float64
	gitSize  float64
	codeSize float64
}

func (r repo) asString() string {
	return fmt.Sprintf("%s: %s MB \n  (%s git / %s code)",
		r.name,
		floatAsString(r.fullSize),
		floatAsString(r.gitSize),
		floatAsString(r.codeSize),
	)
}

func floatAsString(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

// DirSizeMB returns the MB size of a folder
func DirSizeMB(path string) float64 {
	var dirSize int64

	readSize := func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			dirSize += file.Size()
		}
		return nil
	}

	filepath.Walk(path, readSize)
	sizeMB := float64(dirSize) / 1024.0 / 1024.0
	return sizeMB
}

func isDir(filePath string) bool {
	fi, err := os.Stat(filePath)
	check(err)
	if fi.IsDir() {
		return true
	}
	return false
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type reposByFullSize []repo

func (ris reposByFullSize) Len() int           { return len(ris) }
func (ris reposByFullSize) Less(i, j int) bool { return ris[i].fullSize > ris[j].fullSize }
func (ris reposByFullSize) Swap(i, j int)      { ris[i], ris[j] = ris[j], ris[i] }
