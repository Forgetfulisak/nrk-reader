package nrk

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	DataDir    = ".config/nrktracker"
	TargetFile = "news.json.gz"
)

type StoredArticle struct {
	Article
	Seen []time.Time `json:"seen"`
}

func (sa *StoredArticle) Print() {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	encoder.Encode(sa)
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
	os.Rename(file, backup)

	f, err := os.Create(file)
	if err != nil {
		// Could not create new file.
		// Replace backup and abort
		os.Rename(backup, file)
		return err
	}
	defer f.Close()

	fmt.Printf("storing to %s. Currently %d unique articles\n", file, len(news))
	err = encode(f, news)
	if err != nil {
		// Could not write new data.
		// Replace file with backup and abort
		os.Rename(backup, file)
		return err
	}

	// Only remove backup if everything succeeded
	os.Remove(backup)
	return nil
}
