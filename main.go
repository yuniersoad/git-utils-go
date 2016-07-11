package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v3"
	"gopkg.in/src-d/go-git.v3/utils/fs"
	"io"
	"os"
	"regexp"
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

	index := -1
	for i, tag := range tags {
		if tag.Name == tagName {
			index = i
			break
		}
	}
	fmt.Println(index)
	co1, _ := tags[index].Commit()
	co2, _ := tags[index-1].Commit()

	for {
		regexAuthyTicket := regexp.MustCompile(`(AUTHYC-\d+)`)
		ticketId := regexAuthyTicket.FindStringSubmatch(co1.Message)[1]
		fmt.Println(ticketId)
		piter := co1.Parents()
		co1, _ = piter.Next()
		if co1.ID() == co2.ID() {
			break
		}
	}
}
