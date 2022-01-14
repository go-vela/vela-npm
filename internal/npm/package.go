// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.
package npm

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	log "github.com/sirupsen/logrus"
)

type packageJSON struct {
	Name          string        `json:"name"`
	Version       string        `json:"version"`
	PublishConfig publishConfig `json:"publishConfig"`
	Workspaces    []string      `json:"workspaces"`
}

type publishConfig struct {
	Registry string `json:"registry"`
}

// Validate makes sure basic package information is present
func (p *packageJSON) Validate(registry string) error {
	if len(p.Name) == 0 {
		return errors.New("Name not found in package.json")
	}

	if len(p.Version) == 0 {
		return errors.New("Version not found in package.json")
	}

	if _, err := semver.NewConstraint(p.Version); err != nil {
		return fmt.Errorf("Package version error: %w", err)
	}

	// make sure given registry matches what's in "publishConfig"
	// https://docs.npmjs.com/files/package.json#publishconfig
	if len(p.PublishConfig.Registry) != 0 && len(registry) != 0 {
		if p.PublishConfig.Registry != registry {
			return fmt.Errorf("PublishConfig registry %s does not match given registry %s", p.PublishConfig.Registry, registry)
		}
		log.Trace("Registry matches the registry parameter")
	}

	return nil
}
