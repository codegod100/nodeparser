package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Edge struct {
	Incoming string
	Outgoing string
}
type Graph struct {
	Edges []*Edge
}

type File struct {
	NodeName string
	Path     string
	Content  []byte
	Outlinks []string
}
type User struct {
	Name  string
	Files []*File
}

var GRAPH = &Graph{}

func UserFiles(path string) (files []*File) {
	edges := GRAPH.Edges
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".md" {
				// fmt.Println(path)
				content, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}
				outlinks := ParseLinks(string(content))
				nodeName := strings.TrimSuffix(info.Name(), ".md")
				files = append(files, &File{
					NodeName: nodeName,
					Path:     path,
					Content:  content,
					Outlinks: outlinks,
				})
				for _, outlink := range outlinks {
					edges = append(edges, &Edge{
						Incoming: nodeName,
						Outgoing: outlink,
					})
				}
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
	GRAPH.Edges = edges
	return files
}

func Users() (users []*User) {
	rootdir := "./garden"
	userDirs, err := ioutil.ReadDir(rootdir)

	if err != nil {
		log.Fatal(err)
	}
	for _, file := range userDirs {
		if file.IsDir() {
			files := UserFiles(fmt.Sprintf("%s/%s", rootdir, file.Name()))
			// fmt.Println(paths)
			users = append(users, &User{
				Name:  file.Name(),
				Files: files,
			})
			// fmt.Println(file.Name())
		}
	}
	return users
}

func ParseLinks(content string) []string {
	pattern := `\[\[(\w+)\]\]`
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	submatches := re.FindAllStringSubmatch(content, -1)

	// Create a slice to store the words
	words := make([]string, 0)

	// Loop through the submatches and append the words to the slice
	for _, submatch := range submatches {
		words = append(words, submatch[1])
	}
	return words
}

func main() {
	users := Users()
	for _, user := range users {
		if user.Name == "vera" {
			for _, file := range user.Files {
				fmt.Println("NODE:", file.NodeName, "OUTLINKS", file.Outlinks)
			}
		}
	}

	for _, edge := range GRAPH.Edges {
		fmt.Println("INCOMING", edge.Incoming, "OUTGOING", edge.Outgoing)
	}
}
