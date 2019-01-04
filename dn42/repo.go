package dn42

import (
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"strings"
)

var RepoURL = "https://git.dn42.us/dn42/registry"
var repos map[string]billy.Filesystem = make(map[string]billy.Filesystem)

func GetRepoFS(path string) (billy.Filesystem, error) {
	if strings.HasPrefix(path, "http") || strings.HasPrefix(path, "git") {
		return GetRepoFSGIT(path)
	} else {
		return GetRepoFSLocal(path)
	}
}

func GetRepoFSGIT(url string) (billy.Filesystem, error) {
	if val, ok := repos[url]; ok {
		return val, nil
	}
	fs := memfs.New()
	repos[url] = fs
	storer := memory.NewStorage()
	_, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return nil, err
	}
	return fs, nil
}

func GetRepoFSLocal(path string) (billy.Filesystem, error) {
	if val, ok := repos[path]; ok {
		return val, nil
	}
	fs := osfs.New(path)
	repos[path] = fs
	return fs, nil
}

func GetRoutes(filterPath string, routePath string) ([]Route, error) {
	fs, err := GetRepoFS(RepoURL)
	if err != nil {
		return nil, err
	}
	filterFile, err := fs.Open(filterPath)
	if err != nil {
		return nil, err
	}
	defer filterFile.Close()
	filters, err := ParseFilter(filterFile)
	if err != nil {
		return nil, err
	}
	files, err := fs.ReadDir(routePath)
	if err != nil {
		return nil, err
	}
	routes := make([]Route, 0)
	for _, file := range files {
		raw, err := fs.Open(fs.Join(routePath, file.Name()))
		if err != nil {
			raw.Close()
			return nil, err
		}
		rs, err := ParseRoutes(raw, filters)
		raw.Close()
		if err != nil {
			return nil, err
		}
		routes = append(routes, rs...)
	}
	return routes, nil
}
