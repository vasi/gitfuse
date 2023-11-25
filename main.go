package main

import (
	"context"
	"log"
	"os"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

func (f *fsys) Lookup(ctx context.Context, name string) (fs.Node, error) {
	hash := plumbing.NewHash(name)
	commit, err := f.repo.CommitObject(hash)
	if err == plumbing.ErrObjectNotFound {
		return nil, syscall.ENOENT
	}
	if err != nil {
		return nil, err
	}
	return &commitFile{commit}, nil
}

type commitFile struct {
	commit *object.Commit
}

func (f *commitFile) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Mode = 0o444
	a.Size = uint64(len(f.commit.Author.Name))
	return nil
}

func (f *commitFile) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(f.commit.Author.Name), nil
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

	conn, err := fuse.Mount(mountpoint)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	fs.Serve(conn, &fsys{repo})
}
