//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package command

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/teramoby/speedle-plus/cmd/spctl/client"

	"github.com/teramoby/speedle-plus/api/pms"

	"github.com/spf13/cobra"
)

var (
	all         bool
	serviceName string
)

var (
	getExample = `
		# List all services
		spctl get service --all 

		# List services "foo"
		spctl get service foo
		
		# List all policies in service "foo"
		spctl get policy --all --service-name=foo
		
		# List the policy with id "1" in service "foo"
		spctl get policy 1 --service-name=foo
		
		# List all functions
		spctl get function --all
		
		# List function "foo"
		spctl get function foo`
)

func newGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get (service | policy | rolepolicy | function) (--all | NAME | ID) [--service-name=NAME]",
		Short:   "Get one or many services | policies | role-policies",
		Example: getExample,
		Run:     getCommandFunc,
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Get all elements")
	cmd.Flags().StringVar(&serviceName, "service-name", "", "Service name")
	return cmd
}

func getCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		printHelpAndExit(cmd)
	}

	hc, err := httpClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cli := &client.Client{
		PMSEndpoint: globalFlags.PMSEndpoint,
		HTTPClient:  hc,
	}
	var res []byte
	var output = []byte{}

	switch strings.ToLower(args[0]) {
	case "service":
		if all {
			res, err = cli.Get([]string{"service"}, nil, "")
			if err == nil {
				services := []pms.Service{}
				if json.Unmarshal(res, &services) == nil {
					output, _ = json.MarshalIndent(&services, "", strings.Repeat(" ", 4))
				}
			}
		} else {
			if len(args[1:]) == 0 {
				printHelpAndExit(cmd)
			}
			for _, name := range args[1:] {
				service := pms.Service{}
				res, err = cli.Get([]string{"service", name}, nil, "")
				if err != nil {
					break
				}
				if json.Unmarshal(res, &service) == nil {
					s, _ := json.MarshalIndent(&service, "", strings.Repeat(" ", 4))
					output = append(output, s...)
					output = append(output, byte('\n'))
				}
			}
		}
	case "policy", "rolepolicy":
		if serviceName == "" {
			printHelpAndExit(cmd)
		}
		var kind string
		if "policy" == strings.ToLower(args[0]) {
			kind = "policy"
		} else {
			kind = "role-policy"
		}
		if all {
			res, err = cli.Get([]string{"service", serviceName, kind}, nil, "")

			if err == nil {
				var policies interface{}
				if kind == "policy" {
					policies = []pms.Policy{}
				} else {
					policies = []pms.RolePolicy{}
				}

				if json.Unmarshal(res, &policies) == nil {
					output, _ = json.MarshalIndent(&policies, "", strings.Repeat(" ", 4))
				}
			}
		} else {
			if len(args[1:]) == 0 {
				printHelpAndExit(cmd)
			}
			for _, name := range args[1:] {
				var policy interface{}
				if kind == "policy" {
					policy = pms.Policy{}
				} else {
					policy = pms.RolePolicy{}
				}
				res, err = cli.Get([]string{"service", serviceName, kind, name}, nil, "")

				if err != nil {
					break
				}
				if json.Unmarshal(res, &policy) == nil {
					s, _ := json.MarshalIndent(&policy, "", strings.Repeat(" ", 4))
					output = append(output, s...)
					output = append(output, byte('\n'))
				}
			}
		}
	case "function":
		if all {
			res, err = cli.Get([]string{"function"}, nil, "")
			if err == nil {
				functions := []pms.Function{}
				if json.Unmarshal(res, &functions) == nil {
					output, _ = json.MarshalIndent(&functions, "", strings.Repeat(" ", 4))
				}
			}
		} else {
			if len(args[1:]) == 0 {
				printHelpAndExit(cmd)
			}
			for _, name := range args[1:] {
				function := pms.Function{}
				res, err = cli.Get([]string{"function", name}, nil, "")
				if err != nil {
					break
				}
				if json.Unmarshal(res, &function) == nil {
					s, _ := json.MarshalIndent(&function, "", strings.Repeat(" ", 4))
					output = append(output, s...)
					output = append(output, byte('\n'))
				}
			}
		}
	default:
		printHelpAndExit(cmd)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}
