package main

import (
	"context"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type fsys string

func (f fsys) Root() (fs.Node, error) {
	return rootNode(f), nil
}

type rootNode string

func (n rootNode) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Mode = os.ModeDir | 0o555
	return nil
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("Usage: gitfuse REPO MOUNTPOINT")
	}
	repo := os.Args[1]
	mountpoint := os.Args[2]

	conn, err := fuse.Mount(mountpoint)
	if err != nil {
		log.Fatalln(err)
	}

	fs.Serve(conn, fsys(repo))
}
