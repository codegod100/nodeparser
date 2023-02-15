package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	NodeName string
	Path     string
	Content  []byte
}
type User struct {
	Name  string
	Files []*File
}

func UserFiles(path string) (files []*File) {
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
				files = append(files, &File{
					NodeName: strings.TrimSuffix(info.Name(), ".md"),
					Path:     path,
					Content:  content,
				})
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
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

func main() {
	users := Users()
	for _, user := range users {
		if user.Name == "vera" {
			for _, file := range user.Files {
				fmt.Println(file.NodeName)
			}
		}
	}
}
