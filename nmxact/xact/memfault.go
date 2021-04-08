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

package xact

import (
	"mynewt.apache.org/newtmgr/nmxact/nmp"
	"mynewt.apache.org/newtmgr/nmxact/sesn"
)

//////////////////////////////////////////////////////////////////////////////
// $memfaultpull                                                            //
//////////////////////////////////////////////////////////////////////////////

const (
	MEMFAULT_PACKETIZER_STATUS_NO_MORE_DATA = 0
	MEMFAULT_PACKETIZER_STATUS_END_OF_CHUNK = 1
	MEMFAULT_PACKETIZER_STATUS_MORE_DATA_FOR_CHUNK = 2
)

type MemfaultPullProgressFn func(c *MemfaultPullCmd, r *nmp.MemfaultPullRsp)
type MemfaultPullCmd struct {
	CmdBase
	ProgressCb MemfaultPullProgressFn
}

type MemfaultPullResult struct {
	Rsps []*nmp.MemfaultPullRsp
}

func NewMemfaultPullCmd() *MemfaultPullCmd {
	return &MemfaultPullCmd{
		CmdBase: NewCmdBase(),
	}
}

func newMemfaultPullResult() *MemfaultPullResult {
	return &MemfaultPullResult{}
}

func (r *MemfaultPullResult) Status() int {
	rsp := r.Rsps[len(r.Rsps)-1]
	return rsp.Rc
}

func (c *MemfaultPullCmd) Run(s sesn.Sesn) (Result, error) {
	res := newMemfaultPullResult()

	for {
		r := nmp.NewMemfaultPullReq()

		rsp, err := txReq(s, r.Msg(), &c.CmdBase)
		if err != nil {
			return nil, err
		}
		irsp := rsp.(*nmp.MemfaultPullRsp)

		if c.ProgressCb != nil {
			c.ProgressCb(c, irsp)
		}

		res.Rsps = append(res.Rsps, irsp)
		if irsp.Status == MEMFAULT_PACKETIZER_STATUS_NO_MORE_DATA ||
			irsp.Status == MEMFAULT_PACKETIZER_STATUS_END_OF_CHUNK ||
			len(irsp.Data) == 0 {
			// Download complete.
			break
		}
	}

	return res, nil
}

