package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	grafana "github.com/grafana-tools/sdk"
)

func (c *GrafanaDashboardClient) LoadAll(ctx context.Context) error {
	if _, err := os.Stat(c.dashboardDir); err != nil {
		return nil
	}

	files, err := os.ReadDir(c.dashboardDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".dashboard.json") {
			continue
		}

		err := c.Load(ctx, fmt.Sprintf("%s/%s", c.dashboardDir, file.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *GrafanaDashboardClient) Load(ctx context.Context, filename string) error {
	fmt.Printf("Loading: %s\n", filename)

	basename := regexp.MustCompile("\\.dashboard\\.json").ReplaceAllString(filename, "")

	dashboard, err := c.LoadDashboard(ctx, basename+".dashboard.json")
	if err != nil {
		return err
	}
	metadata, err := c.LoadDashboardMetadata(ctx, basename+".metadata.json")
	if err != nil {
		return err
	}

	c.dashboards[dashboard.UID] = &DashboardData{
		Dashboard: dashboard,
		Metadata:  metadata,
		IsPulled:  false,
	}
	// for _, dashboard := range c.dashboards {
	// 	fmt.Println("Dashboard:", dashboard.Dashboard.UID)
	// 	fmt.Println("  Title:", dashboard.Dashboard.Title)
	// 	fmt.Println("  Version:", dashboard.Dashboard.Version)
	// }

	return nil
}

func (c *GrafanaDashboardClient) LoadDashboard(ctx context.Context, filepath string) (*grafana.Board, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var dashboard grafana.Board
	err = json.Unmarshal(buf, &dashboard)
	if err != nil {
		return nil, err
	}

	return &dashboard, nil
}

func (c *GrafanaDashboardClient) LoadDashboardMetadata(ctx context.Context, filepath string) (*grafana.BoardProperties, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var metadata grafana.BoardProperties
	err = json.Unmarshal(buf, &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}
