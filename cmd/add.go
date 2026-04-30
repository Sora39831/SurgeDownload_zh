package cmd

import (
	"errors"
	"fmt"

	"github.com/SurgeDM/Surge/internal/utils"
	"github.com/SurgeDM/Surge/internal/i18n"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add [url]...",
	Aliases: []string{"get"},
	Short:   i18n.T("Add a new download to the running Surge instance"),
	Long:    i18n.T("Add one or more URLs to the download queue of a running Surge instance."),
	RunE: func(cmd *cobra.Command, args []string) error {
		//initializeGlobally is required to ensure that the config and logger are set up before we attempt to resolve the API connection or read the batch file.
		if err := initializeGlobalState(); err != nil {
			return err
		}

		batchFile, _ := cmd.Flags().GetString("batch")
		output, _ := cmd.Flags().GetString("output")

		var urls []string
		urls = append(urls, args...)

		// 2. URLs from batch file
		if batchFile != "" {
			fileUrls, err := utils.ReadURLsFromFile(batchFile)
			if err != nil {
				return fmt.Errorf("error reading batch file: %w", err)
			}
			urls = append(urls, fileUrls...)
		}

		if len(urls) == 0 {
			_ = cmd.Help()
			return nil
		}

		baseURL, token, err := resolveAPIConnection(true)
		if err != nil {
			return err
		}
		resolvedOutput := resolveClientOutputPath(output)

		// Send downloads to server
		count := 0
		attempted := 0
		for _, arg := range urls {
			url, mirrors := ParseURLArg(arg)
			if url == "" {
				continue
			}
			attempted++
			if err := sendToServer(url, mirrors, resolvedOutput, baseURL, token); err != nil {
				fmt.Printf(i18n.T("Error adding %s: %v\n"), url, err)
				continue
			}
			count++
		}

		if count > 0 {
			fmt.Printf(i18n.T("Successfully added %d downloads.\n"), count)
			return nil
		}

		if attempted > 0 {
			return errors.New(i18n.T("failed to add any downloads"))
		}

		return errors.New(i18n.T("no valid URLs to add"))
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("batch", "b", "", i18n.T("File containing URLs to download (one per line)"))
	addCmd.Flags().StringP("output", "o", "", i18n.T("Output directory (defaults to current working directory)"))
}
