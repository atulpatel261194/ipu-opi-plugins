// Copyright (c) 2023 Intel Corporation.  All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"bytes"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	vsiToVportOffset = 16
	pbPythonEnvVar   = "PROTOCOL_BUFFERS_PYTHON_IMPLEMENTATION=python"
)

var execCommand = exec.Command

func ExecOsCommand(cmdBin string, params ...string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := execCommand(cmdBin, params...)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.WithFields(log.Fields{
			"params": params,
			"err":    err,
			"stdout": stdout.String(),
			"stderr": stderr.String(),
		}).Errorf("error while executing %s", cmdBin)
		return err
	}

	log.WithField("params", params).Debugf("successfully executed %s", cmdBin)
	return nil
}

func GetVportForVsi(vsi int) int {
	return vsi + vsiToVportOffset
}

func GetMacIntValueFromBytes(macAddr []byte) uint64 {
	// see how this works: https://go.dev/play/p/MZnMiotnew2
	hwAddr := net.HardwareAddr(macAddr)
	macStr := strings.Replace(hwAddr.String(), ":", "", -1)
	macToInt, _ := strconv.ParseUint(macStr, 16, 64)
	return macToInt
}

var p4rtCtlCommand = exec.Command

func RunP4rtCtlCommand(p4RtBin string, params ...string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := p4rtCtlCommand(p4RtBin, params...)

	// Set required env var for python implemented protobuf
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, pbPythonEnvVar)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.WithFields(log.Fields{
			"params": params,
			"err":    err,
			"stdout": stdout.String(),
			"stderr": stderr.String(),
		}).Errorf("error while executing %s", p4RtBin)
		return err
	}

	log.WithField("params", params).Debugf("successfully executed %s", p4RtBin)
	return nil
}
