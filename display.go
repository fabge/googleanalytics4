package googleanalytics4

import (
	"fmt"
	"strings"
	"time"
)

func (widget *Widget) createTable(websiteReports []websiteReport) string {
	content := ""

	if len(websiteReports) == 0 {
		return content
	}

	content += widget.createHeader()

	for _, websiteReport := range websiteReports {
		websiteRow := ""

		for _, row := range websiteReport.Report.Rows {
			websiteRow += fmt.Sprintf(" %-20s", websiteReport.Name)
			noDataMonth := widget.settings.months - len(websiteReport.Report.Rows)

			// Fill in requested months with no data from query
			if noDataMonth > 0 {
				websiteRow += strings.Repeat("-         ", noDataMonth)
			}

			if row == nil {
				websiteRow += "No data found for given PropertyId"
			} else {
				// for _, row := range reportRows {
				// 	metrics := row.Metrics

				// 	for _, metric := range metrics {
				// 		websiteRow += fmt.Sprintf("%-10s", metric.Values[0])
				// 	}
				// }
				websiteRow += fmt.Sprintf("%-10s", row.MetricValues[0].Value)
			}

			content += websiteRow + "\n"
		}
	}

	return content
}

func (widget *Widget) createHeader() string {
	// Creates the table header of consisting of Months
	currentMonth := int(time.Now().Month())
	widgetStartMonth := currentMonth - widget.settings.months + 1
	header := "                     "

	for i := widgetStartMonth; i < currentMonth+1; i++ {
		header += fmt.Sprintf("%-10s", time.Month(i))
	}
	header += "\n"

	return header
}
