package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all problems in the Jutge.org archive",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			limit := a.effectiveLimit(0)
			problems, err := a.client.ListProblems(cmd.Context(), limit)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(problems, len(problems))
		},
	}
	return cmd
}

func (a *App) searchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search problems by title or code",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			limit := a.effectiveLimit(0)
			results, err := a.client.Search(cmd.Context(), args[0], limit)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(results, len(results))
		},
	}
	return cmd
}
