package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type User struct {
	Name string
	// Files
}

func userFiles(path string) {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".md" {
				fmt.Println(path)
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
}

func main() {
	rootdir := "./garden"
	var users []User

	userDirs, err := ioutil.ReadDir(rootdir)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range userDirs {
		if file.IsDir() {
			users = append(users, User{Name: file.Name()})
			fmt.Println(file.Name())
		}
	}
	fmt.Println(users)

	userFiles("garden/vera")
}
