// Package main provide miguser command
package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/urfave/cli/v2"
)

const (
	name    = "miguser"
	version = "1.0.1"
)

var CSVHEADER = []string{"UserName", "Password", "Description"}

// outputToCsv returns 2d slices for csv encoding using encoding/csv
func outputToCsv(output *computing.DescribeRemoteAccessVpnGatewaysOutput) [][]string {
	users := [][]string{CSVHEADER}
	for _, ravgw := range output.RemoteAccessVpnGatewaySet {
		for _, user := range ravgw.RemoteUserSet {
			users = append(users, []string{*user.UserName, "", *user.Description})
		}
	}
	return users
}

type emptyUserNameError []string

func (e emptyUserNameError) Error() string {
	return fmt.Sprintf("userName is empty: %#v", e)
}

type emptyPasswordError []string

func (e emptyPasswordError) Error() string {
	return fmt.Sprintf("password is empty: %#v", e)
}

func csvToInput(csvData [][]string, ravgwId string) (*computing.CreateRemoteAccessVpnGatewayUsersInput, error) {
	users := make([]types.RequestRemoteUser, len(csvData)-1)
	for i, line := range csvData[1:] {
		userName := line[0]
		if userName == "" {
			return nil, emptyUserNameError(line)
		}

		password := line[1]
		if password == "" {
			return nil, emptyPasswordError(line)
		}

		// allow empty in description
		description := line[2]

		users[i] = types.RequestRemoteUser{
			UserName:    &userName,
			Description: &description,
			Password:    &password,
		}
	}

	return &computing.CreateRemoteAccessVpnGatewayUsersInput{
		RemoteUser:               users,
		RemoteAccessVpnGatewayId: &ravgwId,
	}, nil
}

func exportAction(ctx *cli.Context) error {
	accessKey := ctx.String("access-key")
	privateAccessKey := ctx.String("secret-access-key")
	region := ctx.String("region")
	ravgwId := ctx.String("ravgwid")

	cfg := nifcloud.NewConfig(accessKey, privateAccessKey, region)
	svc := computing.NewFromConfig(cfg)
	inp := &computing.DescribeRemoteAccessVpnGatewaysInput{
		RemoteAccessVpnGatewayId: []string{ravgwId},
	}

	resp, err := svc.DescribeRemoteAccessVpnGateways(context.TODO(), inp)
	if err != nil {
		return err
	}

	userData := outputToCsv(resp)

	f, err := os.Create(fmt.Sprintf("%s.csv", ravgwId))
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	err = w.WriteAll(userData)
	if err != nil {
		return err
	}

	return nil
}

func importAction(ctx *cli.Context) error {
	accessKey := ctx.String("access-key")
	privateAccessKey := ctx.String("secret-access-key")
	region := ctx.String("region")
	ravgwId := ctx.String("ravgwid")
	srcPath := ctx.String("src")
	cfg := nifcloud.NewConfig(accessKey, privateAccessKey, region)
	svc := computing.NewFromConfig(cfg)

	csvFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	csvData, err := r.ReadAll()
	if err != nil {
		return err
	}

	req, err := csvToInput(csvData, ravgwId)
	if err != nil {
		return err
	}

	_, err = svc.CreateRemoteAccessVpnGatewayUsers(context.TODO(), req)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	app := cli.NewApp()

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "access-key",
			Aliases:  []string{"a"},
			Usage:    "set `ACCESS_KEY`",
			Value:    "",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "secret-access-key",
			Aliases:  []string{"s"},
			Usage:    "set `SECRET_ACCESS_KEY`",
			Value:    "",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "region",
			Aliases:  []string{"r"},
			Usage:    "set `NIFCLOUD_REGION`",
			Value:    "",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "ravgwid",
			Usage:    "`RAVGWID` is ID of your RemoteAccessVpnGateway",
			Value:    "",
			Required: true,
		},
	}

	app.Name = "miguser"
	app.Usage = "migrate RAVGW user to another RAVGW"

	app.Commands = []*cli.Command{
		{
			Name:   "export",
			Usage:  "Export user information from RAVGW",
			Action: exportAction,
			Flags:  flags,
		},
		{
			Name:   "import",
			Usage:  "Import user information to RAVGW",
			Action: importAction,
			Flags: append(flags, &cli.StringFlag{
				Name:     "src",
				Usage:    "`SRC_CSV_FILE` of users information",
				Value:    "",
				Required: true,
			}),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
