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

package tcp

import (
	"fmt"

	"mynewt.apache.org/newtmgr/nmxact/nmxutil"
	"mynewt.apache.org/newtmgr/nmxact/sesn"
)

type TcpXport struct {
	started bool
}

func NewTcpXport() *TcpXport {
	return &TcpXport{}
}

func (x *TcpXport) BuildSesn(cfg sesn.SesnCfg) (sesn.Sesn, error) {
	return NewTcpSesn(cfg)
}

func (x *TcpXport) Start() error {
	if x.started {
		return nmxutil.NewXportError("TCP xport started twice")
	}
	x.started = true
	return nil
}

func (x *TcpXport) Stop() error {
	if !x.started {
		return nmxutil.NewXportError("TCP xport stopped twice")
	}
	x.started = false
	return nil
}

func (x *TcpXport) Tx(bytes []byte) error {
	return fmt.Errorf("unsupported")
}
