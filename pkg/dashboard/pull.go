package dashboard

import (
	"context"
	"fmt"

	grafana "github.com/grafana-tools/sdk"
)

func (c *GrafanaDashboardClient) PullByFolderName(ctx context.Context, folderTitle string) error {
	folder, err := c.getFolderByTitle(ctx, folderTitle)
	if err != nil {
		return err
	}
	fmt.Printf("Pulling dashboards in folder: %s [%s]\n", folder.Title, folder.UID)

	foundDashboards, err := c.client.Search(ctx, grafana.SearchFolderID(folder.ID))
	if err != nil {
		return err
	}
	for _, dashboard := range foundDashboards {
		c.PullByUID(ctx, dashboard.UID)
	}

	return nil
}

func (c *GrafanaDashboardClient) PullByTag(ctx context.Context, tag string) error {
	// c.PullByUID(ctx, "cdidyhrdp7e2oe")

	return fmt.Errorf("not implemented")
}

func (c *GrafanaDashboardClient) PullByUID(ctx context.Context, dashboardUID string) error {
	dashboard, metadata, err := c.client.GetDashboardByUID(ctx, dashboardUID)
	if err != nil {
		return err
	}

	var currentVersion int
	if value, ok := c.dashboards[dashboard.UID]; ok {
		currentVersion = value.Metadata.Version
	} else {
		currentVersion = 0
	}

	if metadata.Version <= currentVersion {
		fmt.Printf("Skip to pull: %s (%d -> %d)\n", dashboard.Title, currentVersion, metadata.Version)
		return nil
	}

	c.dashboards[dashboard.UID] = &DashboardData{
		Metadata:  &metadata,
		Dashboard: &dashboard,
		IsPulled:  true,
	}
	fmt.Printf("Found new dashboard: %s [%s]\n", dashboard.Title, dashboard.UID)
	fmt.Println("  Folder:", metadata.FolderTitle)
	fmt.Println("  Version:", metadata.Version)
	fmt.Println("  UpdatedAt:", metadata.Updated)

	return nil
}
