package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	grafana "github.com/grafana-tools/sdk"
)

func (c *GrafanaDashboardClient) SaveAll(ctx context.Context) error {
	for _, dashboard := range c.dashboards {
		err := c.Save(ctx, dashboard)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *GrafanaDashboardClient) Save(ctx context.Context, dashboard *DashboardData) error {
	if !dashboard.IsPulled {
		fmt.Println("Skip to save:", dashboard.Dashboard.Title)
		return nil
	}

	saveDir := fmt.Sprintf("%s", c.dashboardDir)
	err := os.MkdirAll(saveDir, 0755)
	if err != nil {
		return err
	}

	basename := fmt.Sprintf(
		"%s/%s - %s",
		saveDir,
		regexp.MustCompile("/").ReplaceAllString(dashboard.Dashboard.Title, " "),
		dashboard.Dashboard.UID,
	)

	if err = c.SaveDashboard(ctx, basename+".dashboard.json", dashboard.Dashboard); err != nil {
		return err
	}
	if err = c.SaveDashboardMetadata(ctx, basename+".metadata.json", dashboard.Metadata); err != nil {
		return err
	}

	return nil
}

func (c *GrafanaDashboardClient) SaveDashboard(ctx context.Context, filepath string, dashboard *grafana.Board) error {
	fmt.Println("Saving:", filepath)
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	json, err := json.Marshal(dashboard)
	if err != nil {
		return err
	}

	_, err = file.Write(json)
	if err != nil {
		return err
	}

	return nil
}

func (c *GrafanaDashboardClient) SaveDashboardMetadata(ctx context.Context, filepath string, metadata *grafana.BoardProperties) error {
	fmt.Println("Saving:", filepath)
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	json, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	_, err = file.Write(json)
	if err != nil {
		return err
	}

	return nil
}
