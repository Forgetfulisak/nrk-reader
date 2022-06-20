package subcmd

import (
	"encoding/json"
	"nrk-reader"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func NewCmdAnalyze() *cobra.Command {
	return &cobra.Command{
		Use:   "analyze",
		Short: "Display news recorded by track",
		RunE: func(cmd *cobra.Command, args []string) error {
			return analyzeRun()
		},
	}
}

func display(news []nrk.StoredArticle) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	encoder.Encode(news)
}

func analyzeRun() error {
	file, err := nrk.DBFile()
	if err != nil {
		return err
	}

	news := nrk.ReadOldNews(file)

	slices.SortFunc(news, func(a, b nrk.StoredArticle) bool {
		if len(a.Seen) == len(b.Seen) && len(a.Seen) > 0 {
			return a.Seen[0].Before(b.Seen[0])
		}

		return len(a.Seen) > len(b.Seen)
	})

	display(news)

	return nil
}
