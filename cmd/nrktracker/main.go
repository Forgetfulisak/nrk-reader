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
	Seen time.Time `json:"seen"`
}

type StoredNews = []StoredArticle

func DBFile() (string, error) {

	home, err := os.UserHomeDir()
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
	if err == nil {
		defer os.Remove(backup)
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	// fmt.Println("storing to...", file, len(news))
	fmt.Println("storing to...", file, len(news))
	err = encode(f, news)
	if err != nil {
		return err
	}

	return nil
}

func toStoredNews(articles []nrk.Article) StoredNews {
	out := make(StoredNews, 0)

	for _, article := range articles {
		stored := StoredArticle{
			Article: article,
			Seen:    time.Now().Truncate(24 * time.Hour),
		}
		out = append(out, stored)
	}

	return out
}

func removeDuplicateNews(news StoredNews) StoredNews {

	slices.SortFunc(news, func(a, b StoredArticle) bool {
		if a.Title == b.Title {
			return a.Seen.Before(b.Seen)
		}
		return a.Title < b.Title
	})

	out := make(StoredNews, 0, len(news))

	var previous StoredArticle
	for _, article := range news {
		// Ignore duplcate articles
		if article.Equal(&previous.Article) {
			continue
		}
		out = append(out, article)
		previous = article
	}

	return out
}

func main() {
	file, err := DBFile()
	if err != nil {
		panic(err)
	}

	old := ReadOldNews(file)

	newArticles, err := nrk.FetchArticles()
	newNews := toStoredNews(newArticles)
	allNews := append(old, newNews...)

	withoutDuplicates := removeDuplicateNews(allNews)

	err = StoreNews(file, withoutDuplicates)

	if err != nil {
		log.Fatalln("error storing the news. ", err)
	}
}
