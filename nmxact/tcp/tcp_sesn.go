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
	"net"
	"time"

	"github.com/runtimeco/go-coap"

	"mynewt.apache.org/newtmgr/nmxact/mgmt"
	"mynewt.apache.org/newtmgr/nmxact/nmcoap"
	"mynewt.apache.org/newtmgr/nmxact/nmp"
	"mynewt.apache.org/newtmgr/nmxact/nmxutil"
	"mynewt.apache.org/newtmgr/nmxact/omp"
	"mynewt.apache.org/newtmgr/nmxact/sesn"
)

const MAX_PACKET_SIZE = 512 // Completely arbitrary.

type TcpSesn struct {
	cfg    sesn.SesnCfg
	conn   net.Conn
	txvr   *mgmt.Transceiver
	isOpen bool
}

func NewTcpSesn(cfg sesn.SesnCfg) (*TcpSesn, error) {
	s := &TcpSesn{
		cfg: cfg,
	}

	var err error

	s.txvr, err = mgmt.NewTransceiver(cfg.TxFilterCb, cfg.RxFilterCb, true,
		cfg.MgmtProto, 3)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *TcpSesn) listen() {
	go func() {
		b := make([]byte, 2048)

		for {
			n, err := s.conn.Read(b)
			if err != nil {
				s.Close()
				return
			}

			s.txvr.DispatchNmpRsp(b[:n])
		}
	}()
}

func NewConnectedTcpSesn(cfg sesn.SesnCfg, conn net.Conn) (*TcpSesn, error) {
	s, err := NewTcpSesn(cfg)
	if err != nil {
		return nil, err
	}

	s.conn = conn
	s.isOpen = true

	s.listen()

	return s, nil
}

func (s *TcpSesn) Open() error {
	if s.isOpen {
		return nmxutil.NewSesnAlreadyOpenError(
			"Attempt to open an already-open TCP session")
	}

	var err error

	s.conn, err = net.Dial("tcp", s.cfg.PeerSpec.Tcp)
	if err != nil {
		return err
	}

	s.listen()

	return nil
}

func (s *TcpSesn) Close() error {
	if !s.isOpen {
		return nmxutil.NewSesnClosedError(
			"Attempt to close an unopened TCP session")
	}

	s.isOpen = false
	s.conn.Close()
	s.txvr.ErrorAll(nmxutil.NewSesnClosedError("closed"))
	s.txvr.Stop()
	return nil
}

func (s *TcpSesn) IsOpen() bool {
	return s.isOpen
}

func (s *TcpSesn) MtuIn() int {
	return MAX_PACKET_SIZE -
		omp.OMP_MSG_OVERHEAD -
		nmp.NMP_HDR_SIZE
}

func (s *TcpSesn) MtuOut() int {
	return MAX_PACKET_SIZE -
		omp.OMP_MSG_OVERHEAD -
		nmp.NMP_HDR_SIZE
}

func (s *TcpSesn) TxRxMgmt(m *nmp.NmpMsg,
	timeout time.Duration) (nmp.NmpRsp, error) {

	if !s.IsOpen() {
		return nil, fmt.Errorf("Attempt to transmit over closed TCP session")
	}

	txRaw := func(b []byte) error {
		_, err := s.conn.Write(b)
		return err
	}
	return s.txvr.TxRxMgmt(txRaw, m, s.MtuOut(), timeout)
}

func (s *TcpSesn) AbortRx(seq uint8) error {
	s.txvr.ErrorAll(fmt.Errorf("Rx aborted"))
	return nil
}

func (s *TcpSesn) TxCoap(m coap.Message) error {
	txRaw := func(b []byte) error {
		_, err := s.conn.Write(b)
		return err
	}

	return s.txvr.TxCoap(txRaw, m, s.MtuOut())
}

func (s *TcpSesn) MgmtProto() sesn.MgmtProto {
	return s.cfg.MgmtProto
}

func (s *TcpSesn) ListenCoap(mc nmcoap.MsgCriteria) (*nmcoap.Listener, error) {
	return s.txvr.ListenCoap(mc)
}

func (s *TcpSesn) StopListenCoap(mc nmcoap.MsgCriteria) {
	s.txvr.StopListenCoap(mc)
}

func (s *TcpSesn) CoapIsTcp() bool {
	return true
}

func (s *TcpSesn) RxAccept() (sesn.Sesn, *sesn.SesnCfg, error) {
	return nil, nil, fmt.Errorf("Op not implemented yet")
}

func (s *TcpSesn) RxCoap(opt sesn.TxOptions) (coap.Message, error) {
	return nil, fmt.Errorf("Op not implemented yet")
}

func (s *TcpSesn) Filters() (nmcoap.MsgFilter, nmcoap.MsgFilter) {
	return s.txvr.Filters()
}

func (s *TcpSesn) SetFilters(txFilter nmcoap.MsgFilter,
	rxFilter nmcoap.MsgFilter) {

	s.txvr.SetFilters(txFilter, rxFilter)
}
