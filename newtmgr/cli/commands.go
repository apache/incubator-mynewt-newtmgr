/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package cli

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"mynewt.apache.org/newt/util"
	"mynewt.apache.org/newtmgr/newtmgr/nmutil"
	"mynewt.apache.org/newtmgr/nmxact/nmxutil"
)

var NewtmgrLogLevel log.Level
var NewtmgrHelp bool

func Commands() *cobra.Command {
	logLevelStr := ""
	nmCmd := &cobra.Command{
		Use:   nmutil.ToolInfo.ExeName,
		Short: nmutil.ToolInfo.ShortName + " helps you manage remote devices",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			NewtmgrLogLevel, err = log.ParseLevel(logLevelStr)
			if err != nil {
				nmUsage(nil, util.ChildNewtError(err))
			}

			err = util.Init(NewtmgrLogLevel, "", util.VERBOSITY_DEFAULT)
			if err != nil {
				nmUsage(nil, err)
			}
			nmxutil.SetLogLevel(NewtmgrLogLevel)

			// Set cbgo log level if we're using macOS.
			OSSpecificInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	nmCmd.PersistentFlags().StringVarP(&nmutil.ConnProfile, "conn", "c", "",
		"connection profile to use")

	nmCmd.PersistentFlags().Float64VarP(&nmutil.Timeout, "timeout", "t", 10.0,
		"timeout in seconds (partial seconds allowed)")

	nmCmd.PersistentFlags().IntVarP(&nmutil.Tries, "tries", "r", 1,
		"total number of tries in case of timeout")

	nmCmd.PersistentFlags().StringVarP(&logLevelStr, "loglevel", "l", "info",
		"log level to use")

	nmCmd.PersistentFlags().StringVar(&nmutil.DeviceName, "name",
		"", "name of target BLE device; overrides profile setting")

	nmCmd.PersistentFlags().BoolVar(&nmutil.BleWriteRsp, "write-rsp", false,
		"Send BLE acked write requests instead of unacked write commands")

	nmCmd.PersistentFlags().StringVar(&nmutil.ConnType, "conntype", "",
		"Connection type to use instead of using the profile's type")

	nmCmd.PersistentFlags().StringVar(&nmutil.ConnString, "connstring", "",
		"Connection key-value pairs to use instead of using the profile's "+
			"connstring")

	nmCmd.PersistentFlags().StringVar(&nmutil.ConnExtra, "connextra", "",
		"Additional key-value pair to append to the connstring")

	nmCmd.PersistentFlags().StringVar(&nmxutil.OmpRes, "ompres", "/omgr",
		"Use this CoAP resource instead of /omgr")

	versCmd := &cobra.Command{
		Use:     "version",
		Short:   "Display the " + nmutil.ToolInfo.ShortName + " version number",
		Example: "  " + nmutil.ToolInfo.ExeName + " version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s %s\n",
				nmutil.ToolInfo.LongName,
				nmutil.ToolInfo.VersionString)
		},
	}
	nmCmd.AddCommand(versCmd)

	nmCmd.PersistentFlags().IntVarP(&nmutil.HciIdx, "hci", "i",
		0, "HCI index for the controller on Linux machine")

	nmCmd.AddCommand(crashCmd())
	nmCmd.AddCommand(dateTimeCmd())
	nmCmd.AddCommand(fsCmd())
	nmCmd.AddCommand(imageCmd())
	nmCmd.AddCommand(logCmd())
	nmCmd.AddCommand(mempoolStatCmd())
	nmCmd.AddCommand(resetCmd())
	nmCmd.AddCommand(runCmd())
	nmCmd.AddCommand(statsCmd())
	nmCmd.AddCommand(taskStatCmd())
	nmCmd.AddCommand(configCmd())
	nmCmd.AddCommand(connProfileCmd())
	nmCmd.AddCommand(echoCmd())
	nmCmd.AddCommand(resCmd())
	nmCmd.AddCommand(interactiveCmd())
	nmCmd.AddCommand(shellCmd())
	nmCmd.AddCommand(memfaultCmd())

	return nmCmd
}
