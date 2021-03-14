package main

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/swexbe/bagop/internal/pkg/file"
	"github.com/swexbe/bagop/internal/pkg/utility"
)

func printArchives() {

	archives, err := file.GetArchiveIDs(utility.ArchiveIDLocation)
	panicIfErr(err)

	tw := table.NewWriter()

	tw.AppendHeader(table.Row{"TIMESTAMP", "ARCHIVE_ID", "EXPIRES"})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Name: "TIMESTAMP"},
		{Name: "ARCHIVE_ID", AlignHeader: text.AlignCenter},
		{Name: "EXPIRES"},
	})

	for _, archive := range archives {
		timeLayout := "2006-01-02"
		timestamp := archive.Timestamp.Local().Format(timeLayout)
		expires := "Never"
		if archive.Expires {
			expires = archive.ExpiresTimestamp.Local().Format(timeLayout)
		}
		tw.AppendRow(table.Row{timestamp, archive.ArchiveID, expires})
	}
	tw.SetIndexColumn(1)
	tw.SetTitle("Glacier Archives")
	tw.Style().Title.Align = text.AlignCenter

	fmt.Println(tw.Render())

}
