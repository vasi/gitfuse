package main

import (
	"context"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type fsys struct {
	repo *git.Repository
}

func (f *fsys) Root() (fs.Node, error) {
	return f, nil
}

func (f *fsys) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Mode = os.ModeDir | 0o555
	return nil
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("Usage: gitfuse REPO MOUNTPOINT")
	}
	repopath := os.Args[1]
	mountpoint := os.Args[2]

	repo, err := git.PlainOpen(repopath)
	if err != nil {
		log.Fatalln(err)
	}

	commitId := "a5417c87d3db74901a10eae3f74af230811d2886"
	hash := plumbing.NewHash(commitId)
	commit, err := repo.CommitObject(hash)
	if err != nil {
		log.Fatalln(err)
	}
	println(commit.Author.Name)

	conn, err := fuse.Mount(mountpoint)
	if err != nil {
		log.Fatalln(err)
	}

	fs.Serve(conn, &fsys{repo})
}
