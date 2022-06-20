package main

import (
	"log"
	"nrk-reader"
	"time"

	"golang.org/x/exp/slices"
)

// Mutates oldNews
func addArticles(oldNews *nrk.StoredNews, newArticles []nrk.Article) nrk.StoredNews {
	slices.SortFunc(*oldNews, func(a, b nrk.StoredArticle) bool {
		return a.Title < b.Title
	})

	curTime := time.Now().Truncate(24 * time.Hour)
	for _, article := range newArticles {
		newArticle := nrk.StoredArticle{
			Article: article,
			Seen:    []time.Time{curTime},
		}
		dupIdx, found := slices.BinarySearchFunc(*oldNews, newArticle, func(a, b nrk.StoredArticle) int {
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
	file, err := nrk.DBFile()
	if err != nil {
		panic(err)
	}

	old := nrk.ReadOldNews(file)

	newArticles, err := nrk.FetchArticles()
	if err != nil {
		log.Fatalln("error fetching news:", err)
	}

	allNews := addArticles(&old, newArticles)

	err = nrk.StoreNews(file, allNews)
	if err != nil {
		log.Fatalln("error storing the news:", err)
	}
}
