// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.
package npm

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-vela/vela-npm/internal/shell"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Plugin interface.
type Plugin interface {
	Validate() error
	Exec() error
}

// implementation for Plugin.
type plugin struct {
	config *Config
	cli    shell.OSContext
	os     *afero.Afero
}

type version struct {
	NPM  string `json:"npm"`
	Node string `json:"node"`
}

type publishResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type workspacesPublishResponse map[string]publishResponse

// NewPlugin creates a new plugin struct given Config.
func NewPlugin(c *Config) Plugin {
	return &plugin{
		config: c,
		cli:    shell.NewOSContext(),
		os:     &afero.Afero{Fs: afero.NewOsFs()},
	}
}

// Validate assures plugin is configured correctly.
func (p *plugin) Validate() error {
	if p.os == nil {
		return errors.New("no file system handler provided")
	}

	if p.cli == nil {
		return errors.New("no shell handler provided")
	}

	if err := p.config.Validate(); err != nil {
		return err
	}

	return nil
}

// Exec runs the plugin.
func (p *plugin) Exec() error {
	// check for an .npmrc file in the root, if it exists, rename it as it could interfere with our configuration
	log.Trace("Checking for local .npmrc")
	exists, osErr := p.os.Exists(".npmrc")
	if osErr != nil {
		return fmt.Errorf("failed when looking for local .npmrc: %w", osErr)
	}
	if exists {
		log.Trace("Renaming local .npmrc")
		osErr := p.os.Fs.Rename(".npmrc", ".tmp-npmrc")
		if osErr != nil {
			return fmt.Errorf("failed to create rename local .npmrc: %w", osErr)
		}
	}

	pluginErr := p.runSteps()

	// restore .npmrc
	if exists {
		log.Trace("Restoring local .npmrc")
		osErr := p.os.Fs.Rename(".tmp-npmrc", ".npmrc")
		if osErr != nil {
			return fmt.Errorf("failed to restore rename local .npmrc: %w", osErr)
		}
	}

	if pluginErr != nil {
		return pluginErr
	}
	return nil
}

func (p *plugin) runSteps() error {
	// run through plugin steps
	if err := p.createNpmrc(); err != nil {
		return err
	}

	if err := p.verifyNpm(); err != nil {
		return err
	}

	if err := p.authenticate(); err != nil {
		return err
	}
	// check for workspaces in root package.json
	workspaces, err := p.checkForWorkspaces()
	if err != nil {
		log.Debug("Failed to get workspaces %w", err)
	}
	// if not working with workspaces, use root
	if !p.config.Workspaces && len(p.config.Workspace) == 0 {
		// using workspaces but none specified
		if len(workspaces) > 0 {
			return errors.New("using workspaces but none are specified")
		}

		np, err := p.verifyPackage(".")
		if err != nil {
			return fmt.Errorf("failed to verify package.json: %w", err)
		}

		if err := p.validatePackageVersion(np); err != nil {
			return err
		}
	}

	if len(workspaces) > 0 {
		// if specific workspace is given, filter only that one
		if len(p.config.Workspace) > 0 {
			workspaces = []string{p.config.Workspace}
		}

		for _, w := range workspaces {
			np, err := p.verifyPackage(w)
			if err != nil {
				return fmt.Errorf("failed to verify package.json: %w", err)
			}

			if err := p.validatePackageVersion(np); err != nil {
				return err
			}
		}
	}

	if err := p.audit(); err != nil {
		return err
	}

	if err := p.publish(); err != nil {
		return err
	}

	return nil
}

// VerifyNpm makes sure npm command exists.
func (p *plugin) verifyNpm() error {
	// verify npm exists and can by run
	// https://docs.npmjs.com/cli/version
	o, err := p.cli.RunCommandBytes("npm", "version")

	if err != nil {
		return fmt.Errorf("NPM version command failed: %w", err)
	}

	// convert to JSON to display npm version
	var versions version

	err = json.Unmarshal(o, &versions)
	if err != nil {
		return fmt.Errorf("failed to convert npm version response to JSON: %w", err)
	}

	log.WithFields(log.Fields{
		"npm":  versions.NPM,
		"node": versions.Node,
	}).Info("Verifying npm command")

	return nil
}

func (p *plugin) checkForWorkspaces() ([]string, error) {
	log.Trace("Checking for workspaces...")

	nodePackage := packageJSON{}

	f, err := p.os.ReadFile("package.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read package.json %w", err)
	}

	if err := json.Unmarshal(f, &nodePackage); err != nil {
		return nil, fmt.Errorf("failed to marshall package.json: %w", err)
	}

	if len(nodePackage.Workspaces) > 0 {
		log.Trace(nodePackage.Workspaces)

		return nodePackage.Workspaces, nil
	}

	return nil, errors.New("no workspaces found")
}

// verifyPackage makes sure the current package version is not already in the registry.
func (p *plugin) verifyPackage(prefix string) (packageJSON, error) {
	log.Trace("Verifying node package...")

	nodePackage := packageJSON{}

	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	f, err := p.os.ReadFile(prefix + "package.json")
	if err != nil {
		return nodePackage, fmt.Errorf("failed to read package.json %w", err)
	}

	if err := json.Unmarshal(f, &nodePackage); err != nil {
		return nodePackage, fmt.Errorf("failed to marshall package.json: %w", err)
	}

	if err := nodePackage.Validate(p.config.Registry); err != nil {
		return nodePackage, err
	}

	log.Trace("... node package verified")

	return nodePackage, nil
}

// createNpmrc creates .npmrc file to be used by npm commands.
func (p *plugin) createNpmrc() error {
	// create file to write to, written in multiple steps
	log.Trace("Creating .npmrc...")

	// set default home directory for root user
	home := "/root"

	// capture current user running commands
	hd, err := p.cli.GetHomeDir()
	if err == nil {
		home = hd
	}

	// create full path for .npmrc file
	fp := filepath.Join(home, ".npmrc")

	log.WithFields(log.Fields{
		"path": fp,
	}).Info("Creating .npmrc configuration file")

	// send Filesystem call to create directory path for .npmrc file
	if err = p.os.Fs.MkdirAll(filepath.Dir(fp), 0777); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := p.os.Create(fp)
	if err != nil {
		return fmt.Errorf("failed to create .npmrc file: %w", err)
	}

	defer f.Close()

	// write defaults
	// use JSON responses
	if _, err = f.WriteString("json=true\n"); err != nil {
		return fmt.Errorf("failed to write json: %w", err)
	}

	log.Debug("json successfully written")
	// no color
	if _, err = f.WriteString("color=false\n"); err != nil {
		return fmt.Errorf("failed to write color: %w", err)
	}

	log.Debug("color successfully written")
	// log level silent
	if _, err = f.WriteString("loglevel=silent\n"); err != nil {
		return fmt.Errorf("failed to write loglevel: %w", err)
	}

	log.Debug("loglevel successfully written")
	// disable update notifier
	if _, err = f.WriteString("update-notifier=false\n"); err != nil {
		return fmt.Errorf("failed to write update-notifier: %w", err)
	}

	log.Debug("update-notifier successfully written")

	// write auth config
	if len(p.config.Token) != 0 {
		// use token
		registry, _ := url.Parse(p.config.Registry)
		registry.Scheme = "" // Reset the scheme to empty. This makes it so we will get a protocol relative URL.
		registryString := registry.String()

		if !strings.HasSuffix(registryString, "/") {
			registryString = registryString + "/"
		}

		log.WithFields(log.Fields{
			"registry": registryString,
		}).Trace("_authToken registry string")

		auth := fmt.Sprintf("%s:_authToken=\"%s\"", registryString, p.config.Token)

		if _, err = f.WriteString(auth + "\n"); err != nil {
			return fmt.Errorf("failed to write _authToken: %w", err)
		}

		log.Debug("_authToken successfully written")
	} else {
		// user username/password
		auth := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", p.config.UserName, p.config.Password)))
		if _, err = f.WriteString("_auth=" + auth + "\n"); err != nil {
			return fmt.Errorf("failed to write _auth: %w", err)
		}

		log.Debug("_auth successfully written")
	}

	// write registry config if it exists
	if len(p.config.Registry) != 0 {
		if _, err = f.WriteString("registry=" + p.config.Registry + "\n"); err != nil {
			return fmt.Errorf("failed to write registry: %w", err)
		}

		log.WithFields(log.Fields{
			"registry": p.config.Registry,
		}).Debug("Registry successfully written")
	}

	// write email config if it exists
	if len(p.config.Email) != 0 {
		if _, err = f.WriteString("email=" + p.config.Email + "\n"); err != nil {
			return fmt.Errorf("failed to write email: %w", err)
		}

		log.WithFields(log.Fields{
			"email": p.config.Email,
		}).Debug("Email successfully written")
	}

	// if strict-ssl is given, write what is given to config, else will rely on npm to default (true)
	// https://docs.npmjs.com/misc/config#strict-ssl
	if p.config.IsStrictSSLSet {
		if _, err = f.WriteString("strict-ssl=" + strconv.FormatBool(p.config.StrictSSL) + "\n"); err != nil {
			return fmt.Errorf("failed to write strict-ssl: %w", err)
		}

		log.WithFields(log.Fields{
			"strict-ssl": strconv.FormatBool(p.config.StrictSSL),
		}).Debug("StrictSSL successfully written")
	}

	// if always-auth is given, write what is given to config, else rely on npm to default (false)
	// https://docs.npmjs.com/misc/config#always-auth
	if p.config.IsAlwaysAuthSet {
		if _, err = f.WriteString("always-auth=" + strconv.FormatBool(p.config.AlwaysAuth) + "\n"); err != nil {
			return fmt.Errorf("failed to write always-auth: %w", err)
		}

		log.WithFields(log.Fields{
			"always-auth": strconv.FormatBool(p.config.AlwaysAuth),
		}).Debug("AlwaysAuth successfully written")
	}

	// Trace will output this command's output. Useful for debugging.
	p.cli.RunCommand("npm", "config", "list")

	log.Trace("... .npmrc successfully written")

	return nil
}

// authenticate attempts to communicate with npm.
func (p *plugin) authenticate() error {
	log.Info("Checking connection and authentication")
	// make sure auth config was written successfully
	// https://docs.npmjs.com/cli/whoami.html
	_, err := p.cli.RunCommandString("npm", "whoami")

	if err != nil {
		return fmt.Errorf("npm authentication failed")
	}

	// running npm ping will verify authentication
	// https://docs.npmjs.com/cli/ping.html
	// this can be skipped because not all registries support this
	if p.config.SkipPing {
		log.Warn("Skipping auth ping")
	} else {
		log.Debug("Attempting ping")

		_, err = p.cli.RunCommand("npm", "ping")
		if err != nil {
			return errors.New("ping failed, authentication unsuccessful")
		}
	}

	log.WithFields(log.Fields{
		"username": p.config.UserName,
	}).Trace("... Authentication completed")

	return nil
}

// validatePackageVersion checks package version against the registry, errors if current version is already there.
func (p *plugin) validatePackageVersion(nodePackage packageJSON) error {
	// we cannot publish a version if it already exists in the registry
	// https://docs.npmjs.com/cli-commands/view.html
	log.WithFields(log.Fields{
		"name":    nodePackage.Name,
		"version": nodePackage.Version,
	}).Info("Checking registry for the current version")

	out, cmdErr := p.cli.RunCommandBytes("npm", "view", nodePackage.Name, "versions")
	// There was an error getting versions but doesn't mean we can't run
	if cmdErr != nil {
		log.Trace(fmt.Errorf("versions command failed: %w", cmdErr))

		var errResp shell.NPMErrorResponse
		if err := json.Unmarshal(out, &errResp); err != nil {
			return fmt.Errorf("failed to convert npm error response: %w", err)
		}

		if errResp.ErrorBlock.Code == "ENOTFOUND" { // ENOTFOUND -> not a valid registry
			return fmt.Errorf(errResp.ErrorBlock.Summary)
		} else if errResp.ErrorBlock.Code == "E404" { // E404 -> valid registry but package doesn't exist yet... so it's ours to take!
			// Notify that we are publishing with a novel package name
			log.Info("Package does not already exist in the registry, publish will claim `" + nodePackage.Name + "`")

			return nil
		}
		// Unknown error response code
		return fmt.Errorf(errResp.ErrorBlock.Summary)
	}

	var versions []string
	if err := json.Unmarshal(out, &versions); err != nil {
		versionString := strings.ReplaceAll(string(out), "\"", "")
		versionString = strings.TrimSuffix(versionString, "\n")

		log.WithFields(log.Fields{
			"version": versionString,
		}).Debug("Possibly only one version in registry")
		// if only one version it will be a string instead of array
		versions = append(versions, versionString)
	}

	log.Debug("Versions found:")
	log.Debug(versions)

	for _, v := range versions {
		if v == nodePackage.Version {
			return errors.New("Package of version " + nodePackage.Version + " already exists")
		}
	}

	log.Trace("Version does not already exists in registry")

	return nil
}

func (p *plugin) audit() error {
	if p.config.AuditLevel == None {
		log.Warn("Audit level set to NONE, skipping audit check")

		return nil
	}

	// Running audit will error if given audit-level or higher is found
	// https://docs.npmjs.com/cli/v6/commands/npm-audit
	log.Info("Running audit check")

	out, cmdErr := p.cli.RunCommandBytes("npm", "audit", "--production", "--audit-level="+p.config.AuditLevel)
	if cmdErr != nil {
		log.Trace(fmt.Errorf("audit command failed: %w", cmdErr))

		var errResp shell.NPMErrorResponse
		if err := json.Unmarshal(out, &errResp); (err == nil && errResp != shell.NPMErrorResponse{}) {
			if errResp.ErrorBlock.Code == "ENOLOCK" { // ENOLOCK -> requires lockfile
				return fmt.Errorf(errResp.ErrorBlock.Summary + " " + errResp.ErrorBlock.Detail)
			} else if errResp.ErrorBlock.Code == "ENOAUDIT" { // ENOAUDIT -> valid registry but it doesn't support audits
				log.Warn(errResp.ErrorBlock.Summary + " Try adding a .npmrc to your project directory or set `audit-level: none`.")
			} else { // Unknown error response code
				return fmt.Errorf(errResp.ErrorBlock.Summary + " " + errResp.ErrorBlock.Detail)
			}
		}
	}

	if cmdErr != nil {
		return fmt.Errorf("audit failed for audit-level=%[1]s, run `npm audit --production --audit-level=%[1]s` to view vulnerabilities that need fixed", p.config.AuditLevel)
	}

	return nil
}

// publish runs the npm publish command.
// https://docs.npmjs.com/cli/publish
func (p *plugin) publish() error {
	log.Info("Building publish command")

	var args = []string{"publish", "--quiet"}

	// to see if publish would be successful but not actually publish we can do a dry run
	if p.config.DryRun {
		log.Info("Doing a dry run")

		args = append(args, "--dry-run")
	}

	if len(p.config.Tag) != 0 {
		log.WithFields(log.Fields{"tag": p.config.Tag}).Info("Tagging package")

		args = append(args, "--tag", p.config.Tag)
	}

	if len(p.config.Access) != 0 {
		log.WithFields(log.Fields{"access": p.config.Access}).Info("Setting package access")

		args = append(args, "--access", p.config.Access)
	}

	if p.config.Workspaces {
		log.Info("Publishing all workspaces")

		args = append(args, "--workspaces")
	}

	if len(p.config.Workspace) > 0 {
		log.Info("Publishing workspace " + p.config.Workspace)

		args = append(args, "--workspace", p.config.Workspace)
	}

	out, err := p.cli.RunCommandBytes("npm", args...)

	if err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	logFields := make(log.Fields)

	if p.config.Workspaces || len(p.config.Workspace) > 0 {
		var res workspacesPublishResponse
		if err := json.Unmarshal(out, &res); err != nil {
			log.Trace("Failed to convert npm publish response")
		} else {
			for w := range res {
				logFields[res[w].Name] = res[w].Version
			}
		}
	} else {
		var res publishResponse
		if err := json.Unmarshal(out, &res); err != nil {
			log.Trace("Failed to convert npm publish response")
		} else {
			logFields[res.Name] = res.Version
		}
	}

	log.WithFields(logFields).Info("Successfully published node package!")

	return nil
}
