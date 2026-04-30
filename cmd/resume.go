package cmd

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/SurgeDM/Surge/internal/i18n"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume <ID>",
	Short: i18n.T("Resume a paused download"),
	Long:  i18n.T("Resume a paused download by its ID. Use --all to resume all paused downloads."),
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initializeGlobalState(); err != nil {
			return err
		}

		all, _ := cmd.Flags().GetBool("all")

		if !all && len(args) == 0 {
			return errors.New(i18n.T("provide a download ID or use --all"))
		}

		if all {
			fmt.Println(i18n.T("Resuming all downloads is not yet implemented for running server."))
			return nil
		}

		return ExecuteAPIAction(args[0], "/resume", http.MethodPost, i18n.T("Resumed download"))
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)
	resumeCmd.Flags().Bool("all", false, i18n.T("Resume all paused downloads"))
}
