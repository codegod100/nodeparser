package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type Edge struct {
	bun.BaseModel `bun:"table:edges,alias:e"`
	Incoming      string
	Outgoing      string
}
type Graph struct {
	Edges []*Edge
}

type File struct {
	ID            int64 `bun:",pk,autoincrement"`
	bun.BaseModel `bun:"table:files,alias:f"`
	NodeName      string
	Path          string
	Content       []byte
	Outlinks      []string
	User          *User `bun:"rel:belongs-to,join:user_id=id"`
	UserID        int64
}
type User struct {
	ID            int64 `bun:",pk,autoincrement"`
	bun.BaseModel `bun:"table:users,alias:u"`
	Name          string
}

var GRAPH = &Graph{}

func UserFiles(path string, user *User, db *bun.DB) (files []*File) {
	ctx := context.Background()
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
				nodeName := strings.ToLower(strings.TrimSuffix(info.Name(), ".md"))
				uq, err := db.NewInsert().Model(user).Exec(ctx)

				fmt.Println(uq, err, user)
				file := &File{
					NodeName: nodeName,
					Path:     path,
					Content:  content,
					Outlinks: outlinks,
					UserID:   user.ID,
				}
				_, err = db.NewInsert().Model(file).Exec(ctx)
				files = append(files, file)
				for _, outlink := range outlinks {
					edges = append(edges, &Edge{
						Incoming: nodeName,
						Outgoing: strings.ToLower(outlink),
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

func Users(db *bun.DB) (users []*User) {
	rootdir := "./garden"
	userDirs, err := ioutil.ReadDir(rootdir)

	if err != nil {
		log.Fatal(err)
	}
	for _, file := range userDirs {
		if file.IsDir() {
			// fmt.Println(paths)
			user := &User{
				Name: file.Name(),
			}
			// _, err := db.NewInsert().Model(user).Exec(ctx)
			// if err != nil {
			// 	panic(err)
			// }
			UserFiles(fmt.Sprintf("%s/%s", rootdir, file.Name()), user, db)
			users = append(users, user)
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
	ctx := context.Background()
	sqldb, err := sql.Open(sqliteshim.ShimName, "agora.db")
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())
	_, err = db.NewCreateTable().Model((*File)(nil)).Exec(ctx)
	_, err = db.NewCreateTable().Model((*User)(nil)).Exec(ctx)
	if err != nil {
		fmt.Println(err)
	}
	Users(db)

	var people []User
	err = db.NewSelect().Model(&people).OrderExpr("id ASC").Limit(10).Scan(ctx)
	if err != nil {
		panic(err)
	}
	for _, person := range people {
		fmt.Println(person.Name)
	}
}
