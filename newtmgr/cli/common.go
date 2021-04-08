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

	"mynewt.apache.org/newt/util"
	"mynewt.apache.org/newtmgr/newtmgr/bll"
	"mynewt.apache.org/newtmgr/newtmgr/config"
	"mynewt.apache.org/newtmgr/newtmgr/nmutil"
	"mynewt.apache.org/newtmgr/nmxact/mtech_lora"
	"mynewt.apache.org/newtmgr/nmxact/nmble"
	"mynewt.apache.org/newtmgr/nmxact/nmcoap"
	"mynewt.apache.org/newtmgr/nmxact/nmserial"
	"mynewt.apache.org/newtmgr/nmxact/sesn"
	"mynewt.apache.org/newtmgr/nmxact/tcp"
	"mynewt.apache.org/newtmgr/nmxact/udp"
	"mynewt.apache.org/newtmgr/nmxact/xport"
)

var globalSesn sesn.Sesn
var globalXport xport.Xport
var globalP *config.ConnProfile

// This keeps track of whether the global interface has been assigned.  This
// is necessary to accommodate golang's nil-interface semantics.
var globalXportSet bool
var globalTxFilter nmcoap.TxMsgFilter
var globalRxFilter nmcoap.RxMsgFilter

func initConnProfile() error {
	var p *config.ConnProfile

	if nmutil.ConnProfile == "" {
		p = config.NewConnProfile()
		p.Name = "unnamed"
	} else {
		var err error

		cpm := config.GlobalConnProfileMgr()
		p, err = cpm.GetConnProfile(nmutil.ConnProfile)
		if err != nil {
			return err
		}
	}

	if nmutil.ConnType != "" {
		t, err := config.ConnTypeFromString(nmutil.ConnType)
		if err != nil {
			return util.FmtNewtError("invalid conntype: \"%s\"",
				nmutil.ConnType)
		}

		p.Type = t
	}

	if nmutil.ConnString != "" {
		p.ConnString = nmutil.ConnString
	}

	if nmutil.ConnExtra != "" {
		if p.ConnString != "" {
			p.ConnString += ","
		}
		p.ConnString += nmutil.ConnExtra
	}

	if p.Type == config.CONN_TYPE_NONE {
		return util.FmtNewtError("No connection type specified")
	}

	log.Debugf("Using connection profile: %v", p)
	globalP = p

	return nil
}

func getConnProfile() (*config.ConnProfile, error) {
	if globalP == nil {
		if err := initConnProfile(); err != nil {
			return nil, err
		}
	}

	return globalP, nil
}

func GetXport() (xport.Xport, error) {
	if globalXport != nil {
		return globalXport, nil
	}

	cp, err := getConnProfile()
	if err != nil {
		return nil, err
	}

	switch cp.Type {
	case config.CONN_TYPE_SERIAL_PLAIN, config.CONN_TYPE_SERIAL_OIC:
		sc, err := config.ParseSerialConnString(cp.ConnString)
		if err != nil {
			return nil, err
		}

		globalXport = nmserial.NewSerialXport(sc)

	case config.CONN_TYPE_BLL_PLAIN, config.CONN_TYPE_BLL_OIC:
		bc, err := config.ParseBllConnString(cp.ConnString)
		if err != nil {
			return nil, err
		}

		cfg := bll.NewXportCfg()
		if bc.CtlrName != "" {
			cfg.CtlrName = bc.CtlrName
		}
		cfg.OwnAddrType = bc.OwnAddrType
		globalXport = bll.NewBllXport(cfg, bc.HciIdx)

	case config.CONN_TYPE_BLE_PLAIN, config.CONN_TYPE_BLE_OIC:
		bc, err := config.ParseBleConnString(cp.ConnString)
		if err != nil {
			return nil, err
		}
		globalXport, err = config.BuildBleXport(bc)
		if err != nil {
			return nil, err
		}

	case config.CONN_TYPE_UDP_PLAIN, config.CONN_TYPE_UDP_OIC:
		globalXport = udp.NewUdpXport()

	case config.CONN_TYPE_TCP_PLAIN, config.CONN_TYPE_TCP_OIC:
		globalXport = tcp.NewTcpXport()

	case config.CONN_TYPE_MTECH_LORA_OIC:
		cfg := mtech_lora.NewXportCfg()
		globalXport = mtech_lora.NewLoraXport(cfg)

	default:
		return nil, util.FmtNewtError("Unknown connection type: %s (%d)",
			config.ConnTypeToString(cp.Type), int(cp.Type))
	}

	globalXportSet = true

	if err := globalXport.Start(); err != nil {
		return nil, util.ChildNewtError(err)
	}

	return globalXport, nil
}

func GetXportIfOpen() (xport.Xport, error) {
	if !globalXportSet {
		return nil, fmt.Errorf("xport not initailized")
	}

	return globalXport, nil
}

func buildSesnCfg() (sesn.SesnCfg, error) {
	sc := sesn.NewSesnCfg()

	cp, err := getConnProfile()
	if err != nil {
		return sc, err
	}

	switch cp.Type {
	case config.CONN_TYPE_SERIAL_PLAIN:
		sc.MgmtProto = sesn.MGMT_PROTO_NMP
		return sc, nil

	case config.CONN_TYPE_SERIAL_OIC:
		sc.MgmtProto = sesn.MGMT_PROTO_OMP
		return sc, nil

	case config.CONN_TYPE_BLE_PLAIN:
		bc, err := config.ParseBleConnString(cp.ConnString)
		if err != nil {
			return sc, err
		}

		x, err := GetXport()
		if err != nil {
			return sc, err
		}
		bx := x.(*nmble.BleXport)

		sc.MgmtProto = sesn.MGMT_PROTO_NMP
		if err := config.FillSesnCfg(bx, bc, &sc); err != nil {
			return sc, err
		}

		return sc, nil

	case config.CONN_TYPE_BLE_OIC:
		bc, err := config.ParseBleConnString(cp.ConnString)
		if err != nil {
			return sc, err
		}

		x, err := GetXport()
		if err != nil {
			return sc, err
		}
		bx := x.(*nmble.BleXport)

		sc.MgmtProto = sesn.MGMT_PROTO_OMP
		if err := config.FillSesnCfg(bx, bc, &sc); err != nil {
			return sc, err
		}

		return sc, nil

	case config.CONN_TYPE_UDP_PLAIN:
		sc.MgmtProto = sesn.MGMT_PROTO_NMP
		sc.PeerSpec.Udp = cp.ConnString

		return sc, nil

	case config.CONN_TYPE_UDP_OIC:
		sc.MgmtProto = sesn.MGMT_PROTO_OMP
		sc.PeerSpec.Udp = cp.ConnString

		return sc, nil

	case config.CONN_TYPE_TCP_PLAIN:
		sc.MgmtProto = sesn.MGMT_PROTO_NMP
		sc.PeerSpec.Tcp = cp.ConnString

		return sc, nil

	case config.CONN_TYPE_TCP_OIC:
		sc.MgmtProto = sesn.MGMT_PROTO_OMP
		sc.PeerSpec.Tcp = cp.ConnString

		return sc, nil

	case config.CONN_TYPE_MTECH_LORA_OIC:
		mc, err := config.ParseMtechLoraConnString(cp.ConnString)
		if err != nil {
			return sc, err
		}
		sc.MgmtProto = sesn.MGMT_PROTO_OMP
		err = config.FillMtechLoraSesnCfg(mc, &sc)
		return sc, err

	default:
		return sc, util.FmtNewtError("Unknown connection type: %s (%d)",
			config.ConnTypeToString(cp.Type), int(cp.Type))
	}

}

func buildBllSesn(cp *config.ConnProfile) (sesn.Sesn, error) {
	bc, err := config.ParseBllConnString(cp.ConnString)
	if err != nil {
		return nil, err
	}

	sc, err := config.BuildBllSesnCfg(bc)
	if err != nil {
		return nil, err
	}

	x, err := GetXport()
	if err != nil {
		return nil, err
	}
	bx := x.(*bll.BllXport)

	switch cp.Type {
	case config.CONN_TYPE_BLL_PLAIN:
		sc.MgmtProto = sesn.MGMT_PROTO_NMP

	case config.CONN_TYPE_BLL_OIC:
		sc.MgmtProto = sesn.MGMT_PROTO_OMP

	default:
		return nil, util.NewNewtError("ERROR")
	}

	sc.TxFilter = globalTxFilter
	sc.RxFilter = globalRxFilter

	s, err := bx.BuildBllSesn(sc)
	if err != nil {
		return nil, util.ChildNewtError(err)
	}

	return s, nil
}

func GetSesn() (sesn.Sesn, error) {
	if globalSesn != nil {
		return globalSesn, nil
	}

	cp, err := getConnProfile()
	if err != nil {
		return nil, err
	}

	var s sesn.Sesn
	if cp.Type == config.CONN_TYPE_BLL_PLAIN ||
		cp.Type == config.CONN_TYPE_BLL_OIC {

		s, err = buildBllSesn(cp)
		if err != nil {
			return nil, util.ChildNewtError(err)
		}
	} else {
		sc, err := buildSesnCfg()
		if err != nil {
			return nil, err
		}
		sc.TxFilter = globalTxFilter
		sc.RxFilter = globalRxFilter

		x, err := GetXport()
		if err != nil {
			return nil, err
		}

		s, err = x.BuildSesn(sc)
		if err != nil {
			return nil, util.ChildNewtError(err)
		}
	}

	globalSesn = s
	if err := globalSesn.Open(); err != nil {
		return nil, util.ChildNewtError(err)
	}

	return globalSesn, nil
}

func GetSesnIfOpen() (sesn.Sesn, error) {
	if globalSesn == nil {
		return nil, fmt.Errorf("sesn not initailized")
	}

	return globalSesn, nil
}

func SetFilters(txFilter nmcoap.TxMsgFilter, rxFilter nmcoap.RxMsgFilter) {
	globalTxFilter = txFilter
	globalRxFilter = rxFilter
}
