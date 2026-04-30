package cmd

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/SurgeDM/Surge/internal/engine/events"
	"github.com/SurgeDM/Surge/internal/utils"
	"github.com/SurgeDM/Surge/internal/i18n"
)

// StartHeadlessConsumer starts a goroutine to consume progress messages and log to stdout
func StartHeadlessConsumer() {
	go func() {
		if GlobalService == nil {
			return
		}
		stream, cleanup, err := GlobalService.StreamEvents(context.Background())
		if err != nil {
			utils.Debug("Failed to start event stream: %v", err)
			return
		}
		defer cleanup()

		for msg := range stream {
			switch m := msg.(type) {
			case events.DownloadStartedMsg:
				fmt.Printf(i18n.T("Started: %s [%s]\n"), m.Filename, truncateID(m.DownloadID))
			case events.DownloadCompleteMsg:
				atomic.AddInt32(&activeDownloads, -1)
				fmt.Printf(i18n.T("Completed: %s [%s] (in %s)\n"), m.Filename, truncateID(m.DownloadID), m.Elapsed)
			case events.DownloadErrorMsg:
				atomic.AddInt32(&activeDownloads, -1)
				fmt.Printf(i18n.T("Error: %s [%s]: %v\n"), m.Filename, truncateID(m.DownloadID), m.Err)
			case events.DownloadQueuedMsg:
				fmt.Printf(i18n.T("Queued: %s [%s]\n"), m.Filename, truncateID(m.DownloadID))
			case events.DownloadPausedMsg:
				fmt.Printf(i18n.T("Paused: %s [%s]\n"), m.Filename, truncateID(m.DownloadID))
			case events.DownloadResumedMsg:
				fmt.Printf(i18n.T("Resumed: %s [%s]\n"), m.Filename, truncateID(m.DownloadID))
			case events.DownloadRemovedMsg:
				fmt.Printf(i18n.T("Removed: %s [%s]\n"), m.Filename, truncateID(m.DownloadID))
			}
		}
	}()
}

// truncateID shortens a UUID to its first 8 characters for display
func truncateID(id string) string {
	if len(id) > 8 {
		return id[:8]
	}
	return id
}
