package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/ks-yuzu/grafana-dashboard-cli/pkg/dashboard"
)

var (
	// baseDir = kingpin.Flag("dir", "specify base directory").Short('d').Default(".").String()

	params = kingpin.New("grafana-dashboard-cli", "A CLI tool for pull/push Grafana dashboards")
)

var commands = map[string]*kingpin.CmdClause{
	"pull":    params.Command("pull", "pull dashboard from grafana"),
	"pullAll": params.Command("pull-all", "pull all dashboards from grafana"),
	"push":    params.Command("push", "push dashboard to grafana"),
	"pushAll": params.Command("push-all", "push all dashboards to grafana"),
}

var args = map[string]map[string]*string{
	"pull": map[string]*string{
		"url":    commands["pull"].Arg("url", "grafana URL").Required().String(),
		"tag":    commands["pull"].Flag("tag", "find dashboards by tag").String(),
		"folder": commands["pull"].Flag("folder", "find dashboards in folder").String(),
	},
	"pullAll": map[string]*string{
		"url":    commands["pullAll"].Arg("url", "grafana URL").Required().String(),
		"tag":    commands["pullAll"].Flag("tag", "find dashboards by tag").String(),
		"folder": commands["pullAll"].Flag("folder", "find dashboards in folder").String(),
	},
	"push": map[string]*string{
		"url": commands["push"].Arg("url", "grafana URL").Required().String(),
		"tag": commands["push"].Flag("tag", "add tag to all dashboard").Default("auto_synced").String(),
	},
	"pushAll": map[string]*string{
		"url": commands["pushAll"].Arg("url", "grafana URL").Required().String(),
		"tag": commands["pushAll"].Flag("tag", "add tag to all dashboard").Default("auto_synced").String(),
	},
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	ctx := context.Background()

	switch kingpin.MustParse(params.Parse(os.Args[1:])) {
	case commands["pull"].FullCommand():
		return pull(ctx)
	case commands["push"].FullCommand():
		return push(ctx)

	case commands["pullAll"].FullCommand():
	case commands["pushAll"].FullCommand():
	}

	return nil
}

func pull(ctx context.Context) error {
	client, err := dashboard.NewGrafanaDashboardClient(*args["pull"]["url"])
	if err != nil {
		return err
	}

	if *args["pull"]["folder"] != "" {
		err = client.LoadAll(ctx)
		if err != nil {
			return err
		}

		err = client.PullByFolderName(ctx, *args["pull"]["folder"])
		if err != nil {
			return err
		}
	} else if *args["pull"]["tag"] != "" {
		err = client.LoadAll(ctx)
		if err != nil {
			return err
		}

		err = client.PullByTag(ctx, *args["pull"]["tag"])
		if err != nil {
			return err
		}
	} else {
		params.FatalUsage("folder or tag must be specified\n")
	}

	err = client.SaveAll(ctx)
	if err != nil {
		return err
	}

	return nil
}

func push(ctx context.Context) error {
	client, err := dashboard.NewGrafanaDashboardClient(*args["push"]["url"])
	if err != nil {
		return err
	}

	err = client.LoadAll(ctx)
	if err != nil {
		return err
	}

	if *args["push"]["tag"] != "" {
		client.AddCommonTag(*args["push"]["tag"])
	}

	err = client.PushAll(ctx)
	if err != nil {
		return err
	}

	return nil
}
