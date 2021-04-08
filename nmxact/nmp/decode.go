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

import (
	"fmt"

	"github.com/ugorji/go/codec"
)

// These aliases just allow the ctor map to fit within 79 columns.
const op_wr = NMP_OP_WRITE_RSP
const op_rr = NMP_OP_READ_RSP
const gr_def = NMP_GROUP_DEFAULT
const gr_img = NMP_GROUP_IMAGE
const gr_sta = NMP_GROUP_STAT
const gr_cfg = NMP_GROUP_CONFIG
const gr_log = NMP_GROUP_LOG
const gr_cra = NMP_GROUP_CRASH
const gr_run = NMP_GROUP_RUN
const gr_fil = NMP_GROUP_FS
const gr_she = NMP_GROUP_SHELL
const gr_mft = NMP_GROUP_MEMFAULT

// Op-Group-Id
type Ogi struct {
	Op    uint8
	Group uint16
	Id    uint8
}

type rspCtor func() NmpRsp

func echoRspCtor() NmpRsp          { return NewEchoRsp() }
func taskStatRspCtor() NmpRsp      { return NewTaskStatRsp() }
func mpStatRspCtor() NmpRsp        { return NewMempoolStatRsp() }
func dateTimeReadRspCtor() NmpRsp  { return NewDateTimeReadRsp() }
func dateTimeWriteRspCtor() NmpRsp { return NewDateTimeWriteRsp() }
func resetRspCtor() NmpRsp         { return NewResetRsp() }
func imageUploadRspCtor() NmpRsp   { return NewImageUploadRsp() }
func imageStateRspCtor() NmpRsp    { return NewImageStateRsp() }
func coreListRspCtor() NmpRsp      { return NewCoreListRsp() }
func coreLoadRspCtor() NmpRsp      { return NewCoreLoadRsp() }
func coreEraseRspCtor() NmpRsp     { return NewCoreEraseRsp() }
func imageEraseRspCtor() NmpRsp    { return NewImageEraseRsp() }
func statReadRspCtor() NmpRsp      { return NewStatReadRsp() }
func statListRspCtor() NmpRsp      { return NewStatListRsp() }
func logReadRspCtor() NmpRsp       { return NewLogShowRsp() }
func logListRspCtor() NmpRsp       { return NewLogListRsp() }
func logModuleListRspCtor() NmpRsp { return NewLogModuleListRsp() }
func logLevelListRspCtor() NmpRsp  { return NewLogLevelListRsp() }
func logClearRspCtor() NmpRsp      { return NewLogClearRsp() }
func crashRspCtor() NmpRsp         { return NewCrashRsp() }
func runTestRspCtor() NmpRsp       { return NewRunTestRsp() }
func runListRspCtor() NmpRsp       { return NewRunListRsp() }
func fsDownloadRspCtor() NmpRsp    { return NewFsDownloadRsp() }
func fsUploadRspCtor() NmpRsp      { return NewFsUploadRsp() }
func configReadRspCtor() NmpRsp    { return NewConfigReadRsp() }
func configWriteRspCtor() NmpRsp   { return NewConfigWriteRsp() }
func shellExecRspCtor() NmpRsp     { return NewShellExecRsp() }
func memfaultPullRspCtor() NmpRsp     { return NewMemfaultPullRsp() }

var rspCtorMap = map[Ogi]rspCtor{
	{op_wr, gr_def, NMP_ID_DEF_ECHO}:         echoRspCtor,
	{op_rr, gr_def, NMP_ID_DEF_TASKSTAT}:     taskStatRspCtor,
	{op_rr, gr_def, NMP_ID_DEF_MPSTAT}:       mpStatRspCtor,
	{op_rr, gr_def, NMP_ID_DEF_DATETIME_STR}: dateTimeReadRspCtor,
	{op_wr, gr_def, NMP_ID_DEF_DATETIME_STR}: dateTimeWriteRspCtor,
	{op_wr, gr_def, NMP_ID_DEF_RESET}:        resetRspCtor,
	{op_wr, gr_img, NMP_ID_IMAGE_UPLOAD}:     imageUploadRspCtor,
	{op_rr, gr_img, NMP_ID_IMAGE_STATE}:      imageStateRspCtor,
	{op_wr, gr_img, NMP_ID_IMAGE_STATE}:      imageStateRspCtor,
	{op_rr, gr_img, NMP_ID_IMAGE_CORELIST}:   coreListRspCtor,
	{op_rr, gr_img, NMP_ID_IMAGE_CORELOAD}:   coreLoadRspCtor,
	{op_wr, gr_img, NMP_ID_IMAGE_CORELOAD}:   coreEraseRspCtor,
	{op_wr, gr_img, NMP_ID_IMAGE_ERASE}:      imageEraseRspCtor,
	{op_rr, gr_sta, NMP_ID_STAT_READ}:        statReadRspCtor,
	{op_rr, gr_sta, NMP_ID_STAT_LIST}:        statListRspCtor,
	{op_rr, gr_log, NMP_ID_LOG_SHOW}:         logReadRspCtor,
	{op_rr, gr_log, NMP_ID_LOG_LIST}:         logListRspCtor,
	{op_rr, gr_log, NMP_ID_LOG_MODULE_LIST}:  logModuleListRspCtor,
	{op_rr, gr_log, NMP_ID_LOG_LEVEL_LIST}:   logLevelListRspCtor,
	{op_wr, gr_log, NMP_ID_LOG_CLEAR}:        logClearRspCtor,
	{op_wr, gr_cra, NMP_ID_CRASH_TRIGGER}:    crashRspCtor,
	{op_wr, gr_run, NMP_ID_RUN_TEST}:         runTestRspCtor,
	{op_rr, gr_run, NMP_ID_RUN_LIST}:         runListRspCtor,
	{op_rr, gr_fil, NMP_ID_FS_FILE}:          fsDownloadRspCtor,
	{op_wr, gr_fil, NMP_ID_FS_FILE}:          fsUploadRspCtor,
	{op_rr, gr_cfg, NMP_ID_CONFIG_VAL}:       configReadRspCtor,
	{op_wr, gr_cfg, NMP_ID_CONFIG_VAL}:       configWriteRspCtor,
	{op_wr, gr_she, NMP_ID_SHELL_EXEC}:       shellExecRspCtor,
	{op_rr, gr_mft, NMP_ID_MEMFAULT_PULL}:    memfaultPullRspCtor,
}

func DecodeRspBody(hdr *NmpHdr, body []byte) (NmpRsp, error) {
	cb := rspCtorMap[Ogi{hdr.Op, hdr.Group, hdr.Id}]
	if cb == nil {
		return nil, fmt.Errorf("Unrecognized NMP op+group+id: %d, %d, %d",
			hdr.Op, hdr.Group, hdr.Id)
	}

	r := cb()
	cborCodec := new(codec.CborHandle)
	dec := codec.NewDecoderBytes(body, cborCodec)

	if err := dec.Decode(r); err != nil {
		return nil, fmt.Errorf("Invalid response: %s", err.Error())
	}

	r.SetHdr(hdr)
	return r, nil
}
