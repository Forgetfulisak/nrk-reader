package main

import (
	"nrk-reader/cmd/nrk/subcmd"

	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use: "nrk",
	}

	cmd.AddCommand(subcmd.NewCmdRead())
	cmd.AddCommand(subcmd.NewCmdTrack())
	cmd.AddCommand(subcmd.NewCmdAnalyze())
	cmd.Execute()
}
