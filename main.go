package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v3"
	"gopkg.in/src-d/go-git.v3/utils/fs"
	"io"
	"os"
	"sort"
)

type TagList []*git.Tag

func (tl TagList) Len() int           { return len(tl) }
func (tl TagList) Swap(i, j int)      { tl[i], tl[j] = tl[j], tl[i] }
func (tl TagList) Less(i, j int) bool { return tl[i].Tagger.When.Before(tl[j].Tagger.When) }

func main() {
	repo := "./repo/.git"
	arg := os.Args[1]
	getcommits(repo, arg)
	fmt.Println(arg)
}

func getcommits(repoPath string, tagName string) {
	fs := fs.NewOS() // a simple proxy for the local host filesystem
	path := repoPath

	repo, err := git.NewRepositoryFromFS(fs, path)
	if err != nil {
		panic(err)
	}

	iter, err := repo.Tags()
	if err != nil {
		panic(err)
	}
	defer iter.Close()
	tags := TagList{}
	for {
		tag, err := iter.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		tags = append(tags, tag)

	}
	sort.Sort(tags)
	fmt.Println(tags)
}
