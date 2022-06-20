package subcmd

import (
	"errors"
	"nrk-reader"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func NewCmdTrack() *cobra.Command {
	return &cobra.Command{
		Use:   "track",
		Short: "Record news currently on nrk.no",
		RunE: func(cmd *cobra.Command, args []string) error {
			return trackRun()
		},
	}
}

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

func trackRun() error {
	file, err := nrk.DBFile()
	if err != nil {
		return err
	}

	old := nrk.ReadOldNews(file)

	newArticles, err := nrk.FetchArticles()
	if err != nil {
		return errors.New("error fetching news: " + err.Error())
	}

	allNews := addArticles(&old, newArticles)

	err = nrk.StoreNews(file, allNews)
	if err != nil {
		return errors.New("error storing the news: " + err.Error())
	}
	return nil
}
