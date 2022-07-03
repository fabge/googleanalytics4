package googleanalytics4

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/wtfutil/wtf/utils"
	"golang.org/x/oauth2/google"
	gaV4 "google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/option"
)

type websiteReport struct {
	Name   string
	Report *gaV4.RunReportResponse
}

func (widget *Widget) Fetch() []websiteReport {
	secretPath, err := utils.ExpandHomeDir(widget.settings.secretFile)
	if err != nil {
		log.Fatalf("Unable to parse secretFile path")
	}

	serviceV4, err := makeReportServiceV4(secretPath)
	if err != nil {
		log.Fatalf("Unable to create v3 Google Analytics Reporting Service")
	}

	visitorsDataArray := getReports(
		serviceV4, widget.settings.propertyIds, widget.settings.months,
	)
	return visitorsDataArray
}

func buildNetClient(secretPath string) *http.Client {
	clientSecret, err := os.ReadFile(filepath.Clean(secretPath))
	if err != nil {
		log.Fatalf("Unable to read secretPath. %v", err)
	}

	jwtConfig, err := google.JWTConfigFromJSON(clientSecret, gaV4.AnalyticsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to get config from JSON. %v", err)
	}

	return jwtConfig.Client(context.Background())
}

func makeReportServiceV4(secretPath string) (*gaV4.Service, error) {
	client := buildNetClient(secretPath)
	svc, err := gaV4.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Failed to create v4 Google Analytics Reporting Service")
	}

	return svc, err
}

func getReports(
	serviceV4 *gaV4.Service, propertyIds map[string]interface{}, displayedMonths int,
) []websiteReport {
	startDate := fmt.Sprintf("%s-01", time.Now().AddDate(0, -displayedMonths+1, 0).Format("2006-01"))
	var websiteReports []websiteReport

	for website, property := range propertyIds {
		// For custom queries: https://ga-dev-tools.appspot.com/dimensions-metrics-explorer/

		req := &gaV4.RunReportRequest{
			Property: property.(string),
			DateRanges: []*gaV4.DateRange{
				{StartDate: startDate, EndDate: "today"},
			},
			Metrics: []*gaV4.Metric{
				{Name: "activeUsers"},
			},
			Dimensions: []*gaV4.Dimension{
				{Name: "month"},
			},
		}
		response, err := serviceV4.Properties.RunReport("properties/310294414", req).Do()

		if err != nil {
			log.Fatalf("GET request to analyticsreporting/v4 returned error with propertyId: %s", property)
		}
		if response.HTTPStatusCode != 200 {
			log.Fatalf("Did not get expected HTTP response code")
		}

		report := websiteReport{Name: website, Report: response}
		websiteReports = append(websiteReports, report)
	}
	return websiteReports
}
