package web

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var distFS embed.FS

// DistFS returns the embedded file system rooted at dist.
func DistFS() (fs.FS, error) {
	return fs.Sub(distFS, "dist")
}
