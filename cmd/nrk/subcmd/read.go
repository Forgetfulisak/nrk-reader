package subcmd

import (
	"errors"
	"fmt"
	"nrk-reader"

	"github.com/spf13/cobra"
)

func NewCmdRead() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Read news currently on nrk.no",
		RunE: func(cmd *cobra.Command, args []string) error {
			return readRun()
		},
	}
	return cmd
}

func readRun() error {
	fmt.Println("fetching...")

	articles, err := nrk.FetchArticles()
	if err != nil {
		return errors.New("error fetching articles: " + err.Error())
	}

	fmt.Printf("%v articles found\n\n", len(articles))
	fmt.Println("press enter to read the next article")
	for _, article := range articles {
		fmt.Scanln()
		article.Print()
	}
	return nil
}
