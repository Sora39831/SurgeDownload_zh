package cmd

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/SurgeDM/Surge/internal/i18n"
	"github.com/SurgeDM/Surge/internal/engine/state"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:     "rm <ID>",
	Aliases: []string{"kill"},
	Short:   i18n.T("Remove a download"),
	Long:    i18n.T("Remove a download by its ID. Use --clean to remove all completed downloads."),
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initializeGlobalState(); err != nil {
			return err
		}

		clean, _ := cmd.Flags().GetBool("clean")

		if !clean && len(args) == 0 {
			return errors.New(i18n.T("provide a download ID or use --clean"))
		}

		if clean {
			// Remove completed downloads from DB
			count, err := state.RemoveCompletedDownloads()
			if err != nil {
				return fmt.Errorf("error cleaning downloads: %w", err)
			}
			fmt.Printf("Removed %d completed downloads.\n", count)
			return nil
		}

		return ExecuteAPIAction(args[0], "/delete", http.MethodPost, i18n.T("Removed download"))
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().Bool("clean", false, i18n.T("Remove all completed downloads"))
}
