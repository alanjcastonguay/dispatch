///////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
///////////////////////////////////////////////////////////////////////

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/go-openapi/swag"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	"github.com/vmware/dispatch/pkg/dispatchcli/i18n"
	client "github.com/vmware/dispatch/pkg/event-manager/gen/client/drivers"
	models "github.com/vmware/dispatch/pkg/event-manager/gen/models"
)

var (
	createEventDriverLong = i18n.T(
		`Create dispatch event driver

Types and Settings:
* vcenter
	- vcenterurl 	string (required) (e.g. <user>:<password>@<vcenter-host>:<vcenter-port> )
		`)
	// TODO: add examples
	createEventDriverExample = i18n.T(``)

	configs = []string{}
)

// NewCmdCreateEventDriver creates command responsible for dispatch function eventDriver creation.
func NewCmdCreateEventDriver(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "event-driver DRIVER_NAME DRIVER_TYPE [--set KEY=VALUE]",
		Short:   i18n.T("Create event driver"),
		Long:    createEventDriverLong,
		Example: createEventDriverExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			err := createEventDriver(out, errOut, cmd, args)
			CheckErr(err)
		},
	}
	cmd.Flags().StringArrayVarP(&configs, "set", "s", []string{}, "set event driver configurations, default: empty")
	return cmd
}

func createEventDriver(out, errOut io.Writer, cmd *cobra.Command, args []string) error {

	driverName := args[0]
	driverType := args[1]

	driverConfig := models.DriverConfig{}
	for _, conf := range configs {
		result := strings.Split(conf, "=")
		if len(result) != 2 {
			fmt.Fprint(errOut, "Invalid Configuration Format, should be --config key=value")
		}
		driverConfig = append(driverConfig, &models.Config{
			Key:   result[0],
			Value: result[1],
		})
	}

	eventDriver := &models.Driver{
		Name:   swag.String(driverName),
		Type:   swag.String(driverType),
		Config: driverConfig,
	}

	params := &client.AddDriverParams{
		Body:    eventDriver,
		Context: context.Background(),
	}
	client := eventManagerClient()

	created, err := client.Drivers.AddDriver(params, GetAuthInfoWriter())
	if err != nil {
		return formatAPIError(err, params)
	}
	if dispatchConfig.Json {
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "    ")
		return encoder.Encode(*created.Payload)
	}
	fmt.Fprintf(out, "Created eventDriver: %s\n", *created.Payload.Name)
	return nil
}
