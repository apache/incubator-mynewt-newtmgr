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
	"os"

	"github.com/spf13/cobra"

	"mynewt.apache.org/newt/util"
	"mynewt.apache.org/newtmgr/newtmgr/nmutil"
	"mynewt.apache.org/newtmgr/nmxact/nmp"
	"mynewt.apache.org/newtmgr/nmxact/xact"
)

func memfaultDownloadCmd(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		nmUsage(cmd, nil)
	}

	tmpName := args[0] + ".tmp"
	progressBytes := 0
	file, err := os.OpenFile(tmpName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		nmUsage(cmd, util.NewNewtError(fmt.Sprintf(
			"Cannot open file %s - %s", tmpName, err.Error())))
	}
	defer os.Remove(tmpName)
	defer file.Close()

	s, err := GetSesn()
	if err != nil {
		nmUsage(nil, err)
	}

	c := xact.NewMemfaultPullCmd()
	c.SetTxOptions(nmutil.TxOptions())
	c.ProgressCb = func(c *xact.MemfaultPullCmd, rsp *nmp.MemfaultPullRsp) {
		progressBytes += len(rsp.Data)
		fmt.Printf("Pulled %d bytes\n", progressBytes)
		if _, err := file.Write(rsp.Data); err != nil {
			nmUsage(nil, util.ChildNewtError(err))
		}
	}

	res, err := c.Run(s)
	if err != nil {
		nmUsage(nil, util.ChildNewtError(err))
	}

	sres := res.(*xact.MemfaultPullResult)
	if sres.Status() != 0 {
		fmt.Printf("Error: %d\n", sres.Status())
		return
	}

	err = file.Sync()
	if err != nil {
		nmUsage(nil, util.ChildNewtError(err))
	}

	err = file.Close()
	if err != nil {
		nmUsage(nil, util.ChildNewtError(err))
	}

	os.Rename(tmpName, args[0])
	fmt.Printf("Done writing memfault file to %s\n", args[0])
}

func memfaultCmd() *cobra.Command {
	memfaultCmd := &cobra.Command{
		Use:   "memfault",
		Short: "Manage Memfault data on a device",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	memfaultEx := "  " + nmutil.ToolInfo.ExeName +
		" -c olimex memfault pull\n"

	memfaultDownloadCmd := &cobra.Command{
		Use:     "pull <memfault-file> -c <conn_profile>",
		Short:   "Pull memfault data from a device",
		Example: memfaultEx,
		Run:     memfaultDownloadCmd,
	}
	memfaultCmd.AddCommand(memfaultDownloadCmd)

	return memfaultCmd
}
