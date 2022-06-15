package main

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"nrk-reader"
	"os"
	"time"

	"golang.org/x/exp/slices"
)

const (
	DataDir    = ".config/nrktracker"
	TargetFile = "news.json.gz"
)

type StoredArticle struct {
	nrk.Article
	Seen []time.Time `json:"seen"`
}

type StoredNews = []StoredArticle

func DBFile() (string, error) {

	home, err := os.UserHomeDir()
	home = "."
	if err != nil {
		return "", err
	}

	dir := fmt.Sprintf("%s/%s", home, DataDir)
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/%s", dir, TargetFile)
	return path, nil
}

func encode(f *os.File, v any) error {
	w := gzip.NewWriter(f)
	defer w.Close()

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(v)
	return err
}

func decode(f *os.File, v any) error {
	w, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	err = json.NewDecoder(w).Decode(v)
	return err
}

func ReadOldNews(file string) StoredNews {
	f, err := os.Open(file)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}

		return nil
	}
	var news StoredNews
	err = decode(f, &news)
	if err != nil {
		log.Fatalln("error reading old news: ", err)
	}

	return news
}

func StoreNews(file string, news StoredNews) error {

	// In case something goes wrong while
	// writing new news to disk
	backup := file + ".bak"
	err := os.Rename(file, backup)

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Printf("storing to %s. Currently %d unique articles\n", file, len(news))
	err = encode(f, news)
	if err != nil {
		return err
	}

	// Only remove backup if everything succeeded
	os.Remove(backup)
	return nil
}

// Mutates oldNews
func addArticles(oldNews *StoredNews, newArticles []nrk.Article) StoredNews {
	slices.SortFunc(*oldNews, func(a, b StoredArticle) bool {
		return a.Title < b.Title
	})

	curTime := time.Now().Truncate(24 * time.Hour)
	for _, article := range newArticles {
		newArticle := StoredArticle{
			Article: article,
			Seen:    []time.Time{curTime},
		}
		dupIdx, found := slices.BinarySearchFunc(*oldNews, newArticle, func(a, b StoredArticle) int {
			if a.Title == b.Title {
				return 0
			} else if a.Title < b.Title {
				return -1
			} else {
				return 1
			}
		})

		if found {
			if !slices.Contains((*oldNews)[dupIdx].Seen, curTime) {
				(*oldNews)[dupIdx].Seen = append((*oldNews)[dupIdx].Seen, curTime)
			}
		} else {
			(*oldNews) = append((*oldNews), newArticle)
		}
	}

	return *oldNews
}

func main() {
	file, err := DBFile()
	if err != nil {
		panic(err)
	}

	old := ReadOldNews(file)

	newArticles, err := nrk.FetchArticles()
	allNews := addArticles(&old, newArticles)

	err = StoreNews(file, allNews)

	if err != nil {
		log.Fatalln("error storing the news. ", err)
	}
}
