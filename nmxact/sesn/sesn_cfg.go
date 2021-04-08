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

package sesn

import (
	"time"

	"mynewt.apache.org/newtmgr/nmxact/bledefs"
	"mynewt.apache.org/newtmgr/nmxact/lora"
	"mynewt.apache.org/newtmgr/nmxact/nmcoap"
)

type MgmtProto int

const (
	MGMT_PROTO_NMP MgmtProto = iota
	MGMT_PROTO_OMP
	MGMT_PROTO_COAP_SERVER
)

var mgmtProtoMap = map[MgmtProto]string{
	MGMT_PROTO_NMP:         "nmp",
	MGMT_PROTO_OMP:         "omp",
	MGMT_PROTO_COAP_SERVER: "coapserver",
}

func (r MgmtProto) String() string {
	return mgmtProtoMap[r]
}

type OnCloseFn func(s Sesn, err error)

type PeerSpec struct {
	Ble bledefs.BleDev
	Udp string // <addr>:<port>
	Tcp string // <addr>:<port>
}

type SesnCfgBleCentral struct {
	ConnTries   int
	ConnTimeout time.Duration
	// XXX: Missing fields.
}

type SesnCfgBle struct {
	// General configuration.
	OwnAddrType  bledefs.BleAddrType
	EncryptWhen  bledefs.BleEncryptWhen
	CloseTimeout time.Duration
	WriteRsp     bool

	// Central configuration.
	Central SesnCfgBleCentral
}

type SesnCfgLora struct {
	Addr        string
	SegSz       int
	ConfirmedTx bool
	Port        uint8
}

type SesnCfg struct {
	// General configuration.
	MgmtProto MgmtProto
	PeerSpec  PeerSpec
	OnCloseCb OnCloseFn

	// Transport-specific configuration.
	Ble  SesnCfgBle
	Lora SesnCfgLora

	// Filters
	TxFilter nmcoap.TxMsgFilter
	RxFilter nmcoap.RxMsgFilter
}

func NewSesnCfg() SesnCfg {
	return SesnCfg{
		// XXX: For now, assume an own address type of random static.  In the
		// future, there will need to be some global default, or something that
		// gets read from blehostd.
		Ble: SesnCfgBle{
			OwnAddrType:  bledefs.BLE_ADDR_TYPE_RANDOM,
			CloseTimeout: 30 * time.Second,
			WriteRsp:     false,

			Central: SesnCfgBleCentral{
				ConnTries:   5,
				ConnTimeout: 10 * time.Second,
			},
		},
		Lora: SesnCfgLora{
			ConfirmedTx: false,
			Port:        lora.COAP_LORA_PORT,
		},
	}
}
