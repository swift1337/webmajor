package web

import (
	"embed"
	"io/fs"
)

//go:embed dashboard/*
var dashboard embed.FS

func DashboardFiles() fs.FS {
	sub, _ := fs.Sub(dashboard, "dashboard")

	return sub
}
