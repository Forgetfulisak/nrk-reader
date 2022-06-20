package main

import (
	"nrk-reader"

	"golang.org/x/exp/slices"
)

func main() {
	file, err := nrk.DBFile()
	if err != nil {
		panic(err)
	}

	old := nrk.ReadOldNews(file)

	slices.SortFunc(old, func(a, b nrk.StoredArticle) bool {
		if len(a.Seen) == len(b.Seen) && len(a.Seen) > 0 {
			return a.Seen[0].Before(b.Seen[0])
		}

		return len(a.Seen) > len(b.Seen)
	})

	for _, storedArticle := range old {
		storedArticle.Print()
	}
}
