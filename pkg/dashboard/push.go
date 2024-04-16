package dashboard

import (
	"context"
	"fmt"

	grafana "github.com/grafana-tools/sdk"
)

func (c *GrafanaDashboardClient) Push(ctx context.Context, dashboard *DashboardData) error {
	folder, err := c.getFolderByTitle(ctx, dashboard.Metadata.FolderTitle)
	if err != nil {
		newFolder, err := c.client.CreateFolder(ctx, grafana.Folder{
			Title: dashboard.Metadata.FolderTitle,
		})
		if err != nil {
			return err
		}

		folder = &newFolder
	}

	_, remoteMetadata, err := c.client.GetDashboardByUID(ctx, dashboard.Dashboard.UID)
	// if err != nil {
	// 	return err
	// }
	if remoteMetadata.Version >= dashboard.Metadata.Version {
		fmt.Printf("Skipping to push: %s [%s]\n", dashboard.Dashboard.Title, dashboard.Dashboard.UID)
		return nil
	}

	fmt.Printf(
		"Pushing dashboard: %s [%s]\n  Folder: %s\n  Version: %d -> %d\n",
		dashboard.Dashboard.Title, dashboard.Dashboard.UID,
		folder.Title,
		remoteMetadata.Version, dashboard.Metadata.Version,
	)

	// Use dashboard without ID to avoid ID/UID mismatch error
	// c.f. https://github.com/grafana/grafana/blob/8f0f0387b8d316d73a279e30d6ad264070047441/pkg/services/dashboards/database/database.go#L288
	tmpDashboard := *dashboard.Dashboard
	tmpDashboard.ID = 0

	_, err = c.client.SetDashboard(ctx, tmpDashboard, grafana.SetDashboardParams{
		FolderID:  folder.ID,
		Overwrite: (remoteMetadata.Version > 0),
	})
	if err != nil {
		// return err
		fmt.Println("Error:", err)
	}

	return nil
}

func (c *GrafanaDashboardClient) PushAll(ctx context.Context) error {
	for _, dashboard := range c.dashboards {
		err := c.Push(ctx, dashboard)
		if err != nil {
			return err
		}
	}

	return nil
}
