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

package nmp

//////////////////////////////////////////////////////////////////////////////
// $memfaultpull                                                            //
//////////////////////////////////////////////////////////////////////////////

type MemfaultPullReq struct {
NmpBase `codec:"-"`
}

type MemfaultPullRsp struct {
NmpBase
Rc     int  `codec:"rc"`
Status int  `codec:"status"`
Data []byte `codec:"data"`
}

func NewMemfaultPullReq() *MemfaultPullReq {
r := &MemfaultPullReq{}
fillNmpReq(r, NMP_OP_READ, NMP_GROUP_MEMFAULT, NMP_ID_MEMFAULT_PULL)
return r
}

func (r *MemfaultPullReq) Msg() *NmpMsg { return MsgFromReq(r) }

func NewMemfaultPullRsp() *MemfaultPullRsp {
return &MemfaultPullRsp{}
}

func (r *MemfaultPullRsp) Msg() *NmpMsg { return MsgFromReq(r) }

