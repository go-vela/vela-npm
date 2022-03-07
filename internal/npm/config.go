// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package npm

import (
	"errors"
	"strings"

	"github.com/Masterminds/semver/v3"
	log "github.com/sirupsen/logrus"
)

// Config inputs.
type Config struct {
	Token           string
	UserName        string
	Password        string
	Registry        string
	Email           string
	StrictSSL       bool
	IsStrictSSLSet  bool
	AlwaysAuth      bool
	IsAlwaysAuthSet bool
	SkipPing        bool
	DryRun          bool
	Tag             string
	AuditLevel      string
	Access          string
	Workspaces      bool
	Workspace       string
}

const (
	// Low audit-level low.
	Low = "low"
	// Moderate audit-level moderate.
	Moderate = "moderate"
	// High audit-level high.
	High = "high"
	// Critical audit-level critical.
	Critical = "critical"
	// None used for skipping audit.
	None = "none"
)

// DefaultRegistry is the default URL for npm.
const DefaultRegistry = "https://registry.npmjs.org"

// Validate assures plugin is configured correctly.
func (p *Config) Validate() error {
	if len(p.Token) == 0 {
		if len(p.UserName) == 0 {
			return errors.New("UserName not provided")
		}

		// not required for some test registries
		if len(p.Password) == 0 {
			log.Warn("Password not provided")
		}
	}

	if len(p.Registry) == 0 {
		log.Infof("Registry not provided, using default registry %s", DefaultRegistry)
	}

	if len(p.Email) == 0 {
		log.Warn("Email not provied")
	}

	if p.SkipPing {
		log.Warn("Pre-publish auth check with registry will be skipped")
	}

	// tags cannot have semantic versioning
	// https://docs.npmjs.com/cli/dist-tag#caveats
	if len(p.Tag) != 0 {
		_, err := semver.NewVersion(p.Tag)
		if err == nil {
			return errors.New("tags should not have semantic versioning")
		}
	}

	switch strings.ToLower(p.AuditLevel) {
	case "l", "low", "all":
		p.AuditLevel = Low
	case "m", "mod", "moderate":
		p.AuditLevel = Moderate
	case "h", "high":
		p.AuditLevel = High
	case "c", "crit", "critical":
		p.AuditLevel = Critical
	case "n", "no", "none":
		p.AuditLevel = None
	default:
		log.Warn("audit_level is not recognized, the npm default (low)")

		p.AuditLevel = Low
	}

	log.WithFields(log.Fields{
		"audit-level": p.AuditLevel,
	}).Debug("audit level set")

	// access should be 'restricted' or 'public"
	// https://docs.npmjs.com/cli/v8/commands/npm-publish#access
	if len(p.Access) != 0 {
		switch p.Access {
		case "public", "restricted":
			break
		default:
			return errors.New("access is not recognized, use 'public' or 'restricted'")
		}
	}

	// workspaces should either be all or one
	if len(p.Workspace) > 0 && p.Workspaces {
		return errors.New("you must either specify a workspace or all workspaces, but not both")
	}

	return nil
}
