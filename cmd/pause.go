package cmd

import (
	"github.com/SurgeDM/Surge/internal/i18n"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var pauseCmd = &cobra.Command{
	Use:   "pause <ID>",
	Short: i18n.T("Pause a download"),
	Long:  i18n.T("Pause a download by its ID. Use --all to pause all downloads."),
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initializeGlobalState(); err != nil {
			return err
		}

		all, _ := cmd.Flags().GetBool("all")

		if !all && len(args) == 0 {
			return fmt.Errorf(i18n.T("provide a download ID or use --all"))
		}

		if all {
			// TODO: Implement /pause-all endpoint or iterate
			fmt.Println(i18n.T("Pausing all downloads is not yet implemented for running server."))
			return nil
		}

		return ExecuteAPIAction(args[0], "/pause", http.MethodPost, i18n.T("Paused download"))
	},
}

func init() {
	rootCmd.AddCommand(pauseCmd)
	pauseCmd.Flags().Bool("all", false, i18n.T("Pause all downloads"))
}
