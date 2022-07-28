package subcmd

import (
	"encoding/json"
	"fmt"
	"nrk-reader"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func NewCmdAnalyze() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Display news recorded by track",
		RunE: func(cmd *cobra.Command, args []string) error {
			return analyzeRun()
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "durations",
		Short: "Displays how many articles are displayed for N days",
		RunE: func(cmd *cobra.Command, args []string) error {
			return analyzeDurationRun()
		},
	})

	return cmd
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

func analyzeDurationRun() error {
	file, err := nrk.DBFile()
	if err != nil {
		return err
	}

	news := nrk.ReadOldNews(file)

	durations := make([]int, 0, 0)
	for _, article := range news {
		l := len(article.Seen)
		if l >= len(durations) {
			durations = append(durations, make([]int, l-len(durations)+1)...)
		}
		durations[l] += 1
	}

	for duration, count := range durations {
		if count > 0 {
			fmt.Printf("%d: %d\n", duration, count)
		}
	}

	return nil
}
