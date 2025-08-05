// SPDX-License-Identifier: Apache-2.0

package shell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"os/user"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// ErrorStruct defines contents of NPMErrorResponse.
type ErrorStruct struct {
	Code    string `json:"code"`
	Summary string `json:"summary"`
	Detail  string `json:"detail"`
}

// NPMErrorResponse is the JSON response for an npm command when it errors.
type NPMErrorResponse struct {
	ErrorBlock ErrorStruct `json:"error"`
}

// OSContext interface for running CLI commands.
type OSContext interface {
	RunCommand(name string, args ...string) (bytes.Buffer, error)
	RunCommandBytes(name string, args ...string) ([]byte, error)
	RunCommandString(name string, args ...string) (string, error)
	GetHomeDir() (string, error)
}

// OSContextImpl an implementation for OSContext using os/exec.
type OSContextImpl struct {
}

// NewOSContext creates an instance of OSContextImpl.
func NewOSContext() OSContext {
	return &OSContextImpl{}
}

// RunCommand used instead of exec.Command.
func (os *OSContextImpl) RunCommand(name string, args ...string) (bytes.Buffer, error) {
	cmd := exec.CommandContext(context.Background(), name, args...)

	var outBuffer, errorBuffer bytes.Buffer

	cmd.Stdout = &outBuffer
	cmd.Stderr = &errorBuffer

	log.WithFields(log.Fields{"cmd": cmd.String()}).Debug("running command")

	err := cmd.Run()

	log.WithFields(log.Fields{"output": "stdout"}).Trace(outBuffer.String())
	log.WithFields(log.Fields{"output": "stderr"}).Trace(errorBuffer.String())

	if err != nil {
		// if command goes to std error it should follow error block format
		if errorBuffer.Len() > 0 {
			var errResp NPMErrorResponse
			// sanitize error output when it's not silent
			re := regexp.MustCompile("(?m)^.*npm ERR+.*")
			cleanBuffer := re.ReplaceAllString(errorBuffer.String(), "")

			if err := json.Unmarshal([]byte(cleanBuffer), &errResp); err != nil {
				log.Trace("Failed to convert npm error response: %w", err)
			} else {
				log.WithFields(log.Fields{
					"code": errResp.ErrorBlock.Code,
				}).Debug(errResp.ErrorBlock.Summary + ":" + errResp.ErrorBlock.Detail)
			}

			return errorBuffer, fmt.Errorf("command failed (%w)", err)
		}

		return outBuffer, fmt.Errorf("command failed(%w)", err)
	}

	return outBuffer, nil
}

// RunCommandBytes used instead of exec.Command.
func (os *OSContextImpl) RunCommandBytes(name string, args ...string) ([]byte, error) {
	o, err := os.RunCommand(name, args...)

	if err != nil {
		return o.Bytes(), err
	}

	return o.Bytes(), nil
}

// RunCommandString used instead of exec.Command.
func (os *OSContextImpl) RunCommandString(name string, args ...string) (string, error) {
	o, err := os.RunCommand(name, args...)

	if err != nil {
		return o.String(), err
	}

	return o.String(), nil
}

// GetHomeDir fetches current user and gets its home directory.
func (os *OSContextImpl) GetHomeDir() (string, error) {
	var home string
	// capture current user running commands
	u, err := user.Current()
	if err == nil {
		// set home directory to current user
		home = u.HomeDir
	}

	return home, nil
}
