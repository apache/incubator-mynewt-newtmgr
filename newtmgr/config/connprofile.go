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

package config

import (
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"mynewt.apache.org/newt/util"
	"mynewt.apache.org/newtmgr/newtmgr/nmutil"
)

type ConnProfileMgr struct {
	profiles map[string]*ConnProfile
}

type ConnType int

type ConnProfile struct {
	Name       string   `json:"MyName"`
	Type       ConnType `json:"MyType"`
	ConnString string   `json:"MyConnString"`
}

func (p *ConnProfile) String() string {
	return fmt.Sprintf("name=%s type=%s connstring=%s",
		p.Name, ConnTypeToString(p.Type), p.ConnString)
}

const (
	CONN_TYPE_NONE ConnType = iota
	CONN_TYPE_SERIAL_PLAIN
	CONN_TYPE_SERIAL_OIC
	CONN_TYPE_BLL_PLAIN
	CONN_TYPE_BLL_OIC
	CONN_TYPE_BLE_PLAIN
	CONN_TYPE_BLE_OIC
	CONN_TYPE_UDP_PLAIN
	CONN_TYPE_UDP_OIC
	CONN_TYPE_TCP_PLAIN
	CONN_TYPE_TCP_OIC
	CONN_TYPE_MTECH_LORA_OIC
)

var connTypeNameMap = map[ConnType]string{
	CONN_TYPE_SERIAL_PLAIN:   "serial",
	CONN_TYPE_SERIAL_OIC:     "oic_serial",
	CONN_TYPE_BLL_PLAIN:      "ble",
	CONN_TYPE_BLL_OIC:        "oic_ble",
	CONN_TYPE_BLE_PLAIN:      "bhd",
	CONN_TYPE_BLE_OIC:        "oic_bhd",
	CONN_TYPE_UDP_PLAIN:      "udp",
	CONN_TYPE_UDP_OIC:        "oic_udp",
	CONN_TYPE_TCP_PLAIN:      "tcp",
	CONN_TYPE_TCP_OIC:        "oic_tcp",
	CONN_TYPE_MTECH_LORA_OIC: "oic_mtech",
	CONN_TYPE_NONE:           "???",
}

func ConnTypeToString(ct ConnType) string {
	return connTypeNameMap[ct]
}

func ConnTypeFromString(s string) (ConnType, error) {
	for k, v := range connTypeNameMap {
		if s == v {
			return k, nil
		}
	}

	return ConnType(0), util.FmtNewtError("Invalid connection type: %s", s)
}

func (t *ConnType) MarshalJSON() ([]byte, error) {
	return json.Marshal(ConnTypeToString(*t))
}

func (ct *ConnType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	var err error
	*ct, err = ConnTypeFromString(s)
	if err != nil {
		*ct = CONN_TYPE_NONE
	}
	return nil
}

func NewConnProfileMgr() (*ConnProfileMgr, error) {
	cpm := &ConnProfileMgr{
		profiles: map[string]*ConnProfile{},
	}

	if err := cpm.Init(); err != nil {
		return nil, err
	}

	return cpm, nil
}

func connProfileCfgFilename() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", util.NewNewtError(err.Error())
	}

	return filepath.Join(dir, nmutil.ToolInfo.CfgFilename), nil
}

func (cpm *ConnProfileMgr) Init() error {
	filename, err := connProfileCfgFilename()
	if err != nil {
		return err
	}

	log.Debugf("Reading connection profiles from %s", filename)
	blob, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return util.ChildNewtError(err)
		}
	}

	var profiles []*ConnProfile
	if err := json.Unmarshal(blob, &profiles); err != nil {
		return util.FmtNewtError("error reading connection profile "+
			"config (%s): %s", filename, err.Error())
	}

	for _, p := range profiles {
		cpm.profiles[p.Name] = p
	}

	return nil
}

type connProfSorter struct {
	cps []*ConnProfile
}

func (s connProfSorter) Len() int {
	return len(s.cps)
}
func (s connProfSorter) Swap(i, j int) {
	s.cps[i], s.cps[j] = s.cps[j], s.cps[i]
}
func (s connProfSorter) Less(i, j int) bool {
	return s.cps[i].Name < s.cps[j].Name
}

func SortConnProfs(cps []*ConnProfile) []*ConnProfile {
	sorter := connProfSorter{
		cps: make([]*ConnProfile, 0, len(cps)),
	}

	for _, p := range cps {
		sorter.cps = append(sorter.cps, p)
	}

	sort.Sort(sorter)
	return sorter.cps
}

func (cpm *ConnProfileMgr) GetConnProfileList() ([]*ConnProfile, error) {
	log.Debugf("Getting list of connection profiles")

	cpList := make([]*ConnProfile, 0, len(cpm.profiles))
	for _, p := range cpm.profiles {
		cpList = append(cpList, p)
	}

	return SortConnProfs(cpList), nil
}

func (cpm *ConnProfileMgr) save() error {
	list, _ := cpm.GetConnProfileList()
	b, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		return util.NewNewtError(err.Error())
	}

	filename, err := connProfileCfgFilename()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		return util.ChildNewtError(err)
	}

	return nil
}

func (cpm *ConnProfileMgr) DeleteConnProfile(name string) error {
	if cpm.profiles[name] == nil {
		return util.FmtNewtError("connection profile \"%s\" doesn't exist",
			name)
	}

	delete(cpm.profiles, name)

	err := cpm.save()
	if err != nil {
		return err
	}

	return nil
}

func (cpm *ConnProfileMgr) AddConnProfile(cp *ConnProfile) error {
	cpm.profiles[cp.Name] = cp

	err := cpm.save()
	if err != nil {
		return err
	}

	return nil
}

func (cpm *ConnProfileMgr) GetConnProfile(pName string) (*ConnProfile, error) {
	// Each section is a connection profile, key values are the contents
	// of that section.
	p := cpm.profiles[pName]
	if p == nil {
		return nil, util.FmtNewtError("connection profile \"%s\" doesn't "+
			"exist", pName)
	}

	return p, nil
}

func NewConnProfile() *ConnProfile {
	return &ConnProfile{}
}

var globalConnProfileMgr *ConnProfileMgr

func GlobalConnProfileMgr() *ConnProfileMgr {
	if globalConnProfileMgr == nil {
		panic("connection profile manager not initialized")
	}
	return globalConnProfileMgr
}

func InitGlobalConnProfileMgr() error {
	if globalConnProfileMgr != nil {
		return util.NewNewtError("connection profile manager initialized twice")
	}

	var err error
	globalConnProfileMgr, err = NewConnProfileMgr()
	if err != nil {
		return err
	}

	return nil
}
