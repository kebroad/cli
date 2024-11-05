package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newReposBlockIngestCmd() *cobra.Command {
	var allRepos bool

	cmd := &cobra.Command{
		Use:   "block-ingest [flags] [repository] <seconds>",
		Short: "Block ingest for repositories",
		Long: `Block ingest for one or all repositories for a specified number of seconds.

Examples:
  # Block ingest for a specific repository
  humioctl repos block-ingest myrepo 3600

  # Block ingest for all repositories
  humioctl repos block-ingest --all 3600`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewApiClient(cmd)

			if !allRepos && len(args) != 2 {
				return fmt.Errorf("requires repository name and seconds, or --all flag with seconds")
			}

			var seconds int
			var err error

			if allRepos {
				if len(args) != 1 {
					return fmt.Errorf("when using --all, only specify seconds")
				}
				seconds, err = strconv.Atoi(args[0])
			} else {
				seconds, err = strconv.Atoi(args[1])
			}

			if err != nil {
				return fmt.Errorf("invalid duration: %w", err)
			}

			if allRepos {
				repos, err := client.Repositories().List()
				if err != nil {
					return fmt.Errorf("error listing repositories: %w", err)
				}

				for _, repo := range repos {
					err := client.Repositories().BlockIngest(repo.Name, seconds)
					if err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "Error blocking ingest for repository %q: %v\n", repo.Name, err)
						continue
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Successfully blocked ingest for repository %q\n", repo.Name)
				}
				return nil
			}

			err = client.Repositories().BlockIngest(args[0], seconds)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully blocked ingest for repository %q\n", args[0])
			return nil
		},
	}

	cmd.Flags().BoolVar(&allRepos, "all", false, "Block ingest for all repositories")

	return cmd
} 