package dashboard

import (
	"context"
	"fmt"
	"slices"

	grafana "github.com/grafana-tools/sdk"

	"github.com/ks-yuzu/grafana-dashboard-cli/pkg/config"
)

type GrafanaDashboardClient struct {
	client *grafana.Client

	// save/load directory in local storage
	dashboardDir string

	// UID -> Dashboard
	dashboards map[string]*DashboardData
}

type DashboardData struct {
	Metadata  *grafana.BoardProperties
	Dashboard *grafana.Board
	IsPulled  bool
}

var clients = map[string]*GrafanaDashboardClient{}

func NewGrafanaDashboardClient(grafanaUrl string) (*GrafanaDashboardClient, error) {
	c := &GrafanaDashboardClient{
		dashboardDir: "./dashboards",
		dashboards:   map[string]*DashboardData{},
	}

	apiToken := config.GetConfig().Grafana[grafanaUrl]

	if clients[grafanaUrl] == nil {
		client, err := grafana.NewClient("https://"+grafanaUrl, apiToken, grafana.DefaultHTTPClient)
		if err != nil {
			return nil, err
		}

		c.client = client
		clients[grafanaUrl] = c
	}

	return c, nil
}

func (c *GrafanaDashboardClient) getFolderByTitle(ctx context.Context, folderTitle string) (*grafana.Folder, error) {
	folders, err := c.client.GetAllFolders(ctx)
	if err != nil {
		return nil, err
	}

	for _, folder := range folders {
		if folder.Title == folderTitle {
			return &folder, nil
		}
	}

	return nil, fmt.Errorf("Failed to find folder: %s", folderTitle)
}

func (c *GrafanaDashboardClient) AddCommonTag(tag string) {
	for _, d := range c.dashboards {
		if !slices.Contains(d.Dashboard.Tags, tag) {
			d.Dashboard.Tags = append(d.Dashboard.Tags, tag)
		}
	}
}
