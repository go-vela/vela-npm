// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.
package npm

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/go-vela/vela-npm/test"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/afero"
)

const npmrcDefaults = "json=true\ncolor=false\nloglevel=silent\nupdate-notifier=false\n"
const auth = "dGVzdHVzZXI6dGVzdHBhc3M="

func createTestPlugin(t *testing.T, c *Config) (*plugin, *test.MockOSContext, afero.Fs) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := test.NewMockOSContext(ctrl)
	a := &afero.Afero{Fs: afero.NewMemMapFs()}

	return &plugin{config: c, cli: m, os: a}, m, a.Fs
}

func TestMain_ValidateNPMCommand_Success(t *testing.T) {
	p, mock, _ := createTestPlugin(t, &Config{})

	res := `{"npm": "6.13.2"}`
	// Asserts that the first and only call to RunCommandBytes() is passed ("npm", "versions").
	// Anything else will fail.
	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"version"})).
		Return([]byte(res), nil)

	// should complete successfully
	if err := p.verifyNpm(); err != nil {
		t.Fail()
	}
}

func TestMain_ValidateNPMCommand_Fail(t *testing.T) {
	p, mock, _ := createTestPlugin(t, &Config{})

	// Asserts that the first and only call to Bar() is passed 99.
	// Anything else will fail.
	mock.
		EXPECT().
		RunCommandBytes("npm", "version").
		Return(nil, errors.New("command failed (return status code 1): npm is not a recognized command"))

	if err := p.verifyNpm(); err == nil {
		t.Fail()
	}
}

func TestPlugin_verifyPackage_Valid(t *testing.T) {
	c := &Config{
		Tag:      "1.0.0",
		UserName: "testuser",
	}
	p, _, fs := createTestPlugin(t, c)
	// write test package.json to file system
	testPackage := &packageJSON{
		Name:    "vela-npm",
		Version: "1.0.0",
	}

	jsonContents, err := json.Marshal(testPackage)
	if err != nil {
		t.Fail()
	}

	err = afero.WriteFile(fs, "package.json", jsonContents, 0644)
	if err != nil {
		t.Fail()
	}

	_, err = p.verifyPackage(".")
	if err != nil {
		t.Error(err)
	}
}

func TestPlugin_verifyPackage_Invalid(t *testing.T) {
	c := &Config{
		Tag:      "1.0.0",
		UserName: "testuser",
	}
	p, _, _ := createTestPlugin(t, c)

	_, err := p.verifyPackage(".")
	if err == nil {
		t.Fail()
	}
}

func TestPlugin_createNpmrc_CreatesFile(t *testing.T) {
	c := &Config{
		Registry: "http://registry.test.com",
		UserName: "testuser",
		Password: "testpass",
	}
	p, mock, fs := createTestPlugin(t, c)
	home := path.Join("usr", "mctestface")
	fs.MkdirAll(home, 0755) // nolint: errcheck // testing
	mock.
		EXPECT().
		GetHomeDir().
		Return(home, nil)
	mock.EXPECT().RunCommand("npm", "config", "list")

	err := p.createNpmrc()
	if err != nil {
		t.Error(err)
	}

	f, err := afero.ReadFile(fs, path.Join(home, ".npmrc"))
	if err != nil {
		t.Error(err)
	}

	npmrc := string(f)
	testNpmrc := fmt.Sprintf("%s//registry.test.com/:_auth=%s\nregistry=http://registry.test.com\n", npmrcDefaults, auth)

	if npmrc != testNpmrc {
		t.Errorf("%s != %s", npmrc, testNpmrc)
	}
}

func TestPlugin_createNpmrc_AuthToken(t *testing.T) {
	c := &Config{
		Token:    "test-token",
		Registry: "http://registry.test.com",
	}
	p, mock, fs := createTestPlugin(t, c)
	home := path.Join("usr", "mctestface")
	fs.MkdirAll(home, 0755) // nolint: errcheck // testing
	mock.
		EXPECT().
		GetHomeDir().
		Return(home, nil)
	mock.EXPECT().RunCommand("npm", "config", "list")

	err := p.createNpmrc()
	if err != nil {
		t.Error(err)
	}

	f, err := afero.ReadFile(fs, path.Join(home, ".npmrc"))
	if err != nil {
		t.Error(err)
	}

	npmrc := string(f)
	testNpmrc := fmt.Sprintf("%s//registry.test.com/:_authToken=\"%s\"\nregistry=%s\n", npmrcDefaults, c.Token, c.Registry)

	if npmrc != testNpmrc {
		t.Errorf("%s != %s", npmrc, testNpmrc)
	}
}

func TestPlugin_createNpmrc_Registry(t *testing.T) {
	c := &Config{
		UserName: "testuser",
		Password: "testpass",
		Registry: "http://registry.test.com",
	}
	p, mock, fs := createTestPlugin(t, c)
	home := path.Join("usr", "mctestface")
	fs.MkdirAll(home, 0755) // nolint: errcheck // testing
	mock.
		EXPECT().
		GetHomeDir().
		Return(home, nil)
	mock.EXPECT().RunCommand("npm", "config", "list")

	err := p.createNpmrc()
	if err != nil {
		t.Error(err)
	}

	f, err := afero.ReadFile(fs, path.Join(home, ".npmrc"))
	if err != nil {
		t.Error(err)
	}

	npmrc := string(f)
	testNpmrc := fmt.Sprintf("%s//registry.test.com/:_auth=%s\nregistry=%s\n", npmrcDefaults, auth, c.Registry)

	if npmrc != testNpmrc {
		t.Errorf("%s != %s", npmrc, testNpmrc)
	}
}

func TestPlugin_createNpmrc_Email(t *testing.T) {
	c := &Config{
		UserName: "testuser",
		Password: "testpass",
		Email:    "testuser@test.com",
		Registry: "http://registry.test.com",
	}
	p, mock, fs := createTestPlugin(t, c)
	home := path.Join("usr", "mctestface")
	fs.MkdirAll(home, 0755) // nolint: errcheck // testing
	mock.
		EXPECT().
		GetHomeDir().
		Return(home, nil)
	mock.EXPECT().RunCommand("npm", "config", "list")

	err := p.createNpmrc()
	if err != nil {
		t.Error(err)
	}

	f, err := afero.ReadFile(fs, path.Join(home, ".npmrc"))
	if err != nil {
		t.Error(err)
	}

	npmrc := string(f)
	auth := b64.StdEncoding.EncodeToString([]byte("testuser:testpass"))
	testNpmrc := fmt.Sprintf("%s//registry.test.com/:_auth=%s\nregistry=http://registry.test.com\nemail=%s\n", npmrcDefaults, auth, c.Email)

	if npmrc != testNpmrc {
		t.Errorf("%s != %s", npmrc, testNpmrc)
	}
}

func TestPlugin_createNpmrc_StrictSSLSet(t *testing.T) {
	c := &Config{
		UserName:       "testuser",
		Password:       "testpass",
		IsStrictSSLSet: true,
		StrictSSL:      true,
	}
	p, mock, fs := createTestPlugin(t, c)
	home := path.Join("usr", "mctestface")
	fs.MkdirAll(home, 0755) // nolint: errcheck // testing
	mock.
		EXPECT().
		GetHomeDir().
		Return(home, nil)
	mock.EXPECT().RunCommand("npm", "config", "list")

	err := p.createNpmrc()
	if err != nil {
		t.Error(err)
	}

	f, err := afero.ReadFile(fs, path.Join(home, ".npmrc"))
	if err != nil {
		t.Error(err)
	}

	npmrc := string(f)
	auth := b64.StdEncoding.EncodeToString([]byte("testuser:testpass"))
	testNpmrc := fmt.Sprintf("%s_auth=%s\nstrict-ssl=true\n", npmrcDefaults, auth)

	if npmrc != testNpmrc {
		t.Errorf("%s != %s", npmrc, testNpmrc)
	}
}

func TestPlugin_createNpmrc_AlwaysAuthSet(t *testing.T) {
	c := &Config{
		UserName:        "testuser",
		Password:        "testpass",
		IsAlwaysAuthSet: true,
		AlwaysAuth:      true,
	}
	p, mock, fs := createTestPlugin(t, c)
	home := path.Join("usr", "mctestface")
	fs.MkdirAll(home, 0755) // nolint: errcheck // testing
	mock.
		EXPECT().
		GetHomeDir().
		Return(home, nil)
	mock.EXPECT().RunCommand("npm", "config", "list")

	err := p.createNpmrc()
	if err != nil {
		t.Error(err)
	}

	f, err := afero.ReadFile(fs, path.Join(home, ".npmrc"))
	if err != nil {
		t.Error(err)
	}

	npmrc := string(f)
	auth := b64.StdEncoding.EncodeToString([]byte("testuser:testpass"))
	testNpmrc := fmt.Sprintf("%s_auth=%s\nalways-auth=true\n", npmrcDefaults, auth)

	if npmrc != testNpmrc {
		t.Errorf("%s != %s", npmrc, testNpmrc)
	}
}

func TestPlugin_createNpmrc_All(t *testing.T) {
	c := &Config{
		UserName:        "testuser",
		Password:        "testpass",
		Registry:        "http://registry.test.com",
		Email:           "testuser@test.com",
		IsStrictSSLSet:  true,
		StrictSSL:       true,
		IsAlwaysAuthSet: true,
		AlwaysAuth:      true,
	}
	p, mock, fs := createTestPlugin(t, c)
	home := path.Join("usr", "mctestface")
	fs.MkdirAll(home, 0755) // nolint: errcheck // testing
	mock.
		EXPECT().
		GetHomeDir().
		Return(home, nil)
	mock.EXPECT().RunCommand("npm", "config", "list")

	err := p.createNpmrc()
	if err != nil {
		t.Error(err)
	}

	f, err := afero.ReadFile(fs, path.Join(home, ".npmrc"))
	if err != nil {
		t.Error(err)
	}

	npmrc := string(f)
	auth := b64.StdEncoding.EncodeToString([]byte("testuser:testpass"))
	testNpmrc := fmt.Sprintf("%s//registry.test.com/:_auth=%s\nregistry=%s\nemail=%s\nstrict-ssl=true\nalways-auth=true\n",
		npmrcDefaults,
		auth,
		c.Registry,
		c.Email)

	if npmrc != testNpmrc {
		t.Errorf("%s != %s", npmrc, testNpmrc)
	}
}

func TestPlugin_authenticate(t *testing.T) {
	p, mock, _ := createTestPlugin(t, &Config{
		Registry: "http://registry.test.com",
	})
	mock.
		EXPECT().
		RunCommandString(gomock.Eq("npm"), gomock.Eq([]string{"whoami", "--registry", "http://registry.test.com"})).
		Times(1).
		Return("testuser", nil)
	mock.
		EXPECT().
		RunCommand(gomock.Eq("npm"), gomock.Eq([]string{"ping", "--registry", "http://registry.test.com"})).
		Times(1).
		Return(bytes.Buffer{}, nil)

	err := p.authenticate()
	if err != nil {
		t.Error(err)
	}
}

func TestPlugin_authenticate_SkipPing(t *testing.T) {
	c := &Config{
		SkipPing: true,
		Registry: "http://registry.test.com",
	}
	p, mock, _ := createTestPlugin(t, c)
	mock.
		EXPECT().
		RunCommandString(gomock.Eq("npm"), gomock.Eq([]string{"whoami", "--registry", "http://registry.test.com"})).
		Times(1).
		Return("testuser", nil)
	mock.
		EXPECT().
		RunCommand(gomock.Eq("npm"), gomock.Eq([]string{"ping", "--registry", "http://registry.test.com"})).
		Times(0).
		Return(bytes.Buffer{}, nil)

	err := p.authenticate()
	if err != nil {
		t.Error(err)
	}
}

func TestPlugin_validatePackageVersion(t *testing.T) {
	p, mock, _ := createTestPlugin(t, &Config{
		Registry: "http://registry.test.com",
	})
	testPackage := packageJSON{
		Name:    "vela-npm",
		Version: "2.0.0",
	}
	res := `["1.0.0"]`
	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"view", "vela-npm", "versions", "--registry", "http://registry.test.com"})).
		Return([]byte(res), nil)

	err := p.validatePackageVersion(testPackage)
	if err != nil {
		t.Error(err)
	}
}

func TestPlugin_validatePackageVersion_RegistryNotFound(t *testing.T) {
	p, mock, _ := createTestPlugin(t, &Config{
		Registry: "http://registry.test.com",
	})
	testPackage := packageJSON{
		Name:    "vela-npm",
		Version: "2.0.0",
	}
	res := `{
		"error": {
			"code": "ENOTFOUND",
			"summary": "request to https://not-a-registry/vela-npm failed, reason: getaddrinfo ENOTFOUND not-a-registry",
			"detail": "This is a problem related to network connectivity.\nIn most cases you are behind a proxy or have bad network settings.\n\nIf you are behind a proxy, please make sure that the\n'proxy' config is set properly.  See: 'npm help config'"
		}
	}`
	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"view", "vela-npm", "versions", "--registry", "http://registry.test.com"})).
		Return([]byte(res), errors.New("Process exited with status code 1"))

	err := p.validatePackageVersion(testPackage)
	if err == nil {
		t.Fail()
	}
}

func TestPlugin_validatePackageVersion_PackageNotFound(t *testing.T) {
	p, mock, _ := createTestPlugin(t, &Config{
		Registry: "http://registry.test.com",
	})
	testPackage := packageJSON{
		Name:    "vela-npm",
		Version: "2.0.0",
	}
	res := `{
		"error": {
			"code": "E404",
			"summary": "Not Found - GET https://registry.npmjs.org/vela-npm - not_found",
			"detail": "\n 'vela-npm@latest' is not in the npm registry.\nYou should bug the author to publish it (or use the name yourself!)\n\nNote that you can also install from a\ntarball, folder, http url, or git url."
		}
	}`
	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"view", "vela-npm", "versions", "--registry", "http://registry.test.com"})).
		Return([]byte(res), errors.New("Process exited with status code 1"))

	err := p.validatePackageVersion(testPackage)
	if err != nil {
		t.Fail()
	}
}

func TestPlugin_validatePackageVersion_Conflict(t *testing.T) {
	p, mock, _ := createTestPlugin(t, &Config{
		Registry: "http://registry.test.com",
	})
	testPackage := packageJSON{
		Name:    "vela-npm",
		Version: "1.0.0",
	}
	res := `["1.0.0"]`
	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"view", "vela-npm", "versions", "--registry", "http://registry.test.com"})).
		Return([]byte(res), nil)

	err := p.validatePackageVersion(testPackage)
	if err == nil {
		t.Fail()
	}
}

func TestPlugin_audit_Skipped(t *testing.T) {
	c := &Config{
		AuditLevel: None,
	}
	p, mock, _ := createTestPlugin(t, c)
	mock.
		EXPECT().
		RunCommandBytes("npm", "audit", "--production", "--audit-level=none").
		Times(0)

	err := p.audit()
	if err != nil {
		t.Error(err)
	}
}

func TestPlugin_audit_ErrorContainLevel(t *testing.T) {
	c := &Config{
		AuditLevel: Low,
	}
	p, mock, _ := createTestPlugin(t, c)
	res := `{
		"metadata": {
			"vulnerabilities": {
				"info": 0,
				"low": 1,
				"moderate": 0,
				"high": 0,
				"critical": 0
			}
		}
	}`

	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"audit", "--production", "--audit-level=low"})).
		Return([]byte(res), fmt.Errorf("Command failed with exit code 1"))

	err := p.audit()

	if err == nil || !strings.Contains(err.Error(), "npm audit --production --audit-level=low") {
		t.Error(fmt.Errorf("audit: the audit command should give feedback to user on how to diagnose error instead got: %w", err))
	}
}

func TestPlugin_audit_FailIfNoPackageLock(t *testing.T) {
	c := &Config{
		AuditLevel: Critical,
	}
	p, mock, _ := createTestPlugin(t, c)
	res := `{
		"error": {
			"code": "ENOLOCK",
			"summary": "This command requires an existing lockfile.",
			"detail": "Try creating one first with: npm i --package-lock-only\nOriginal error: loadVirtual requires existing shrinkwrap file"
		}
	}`

	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"audit", "--production", "--audit-level=critical"})).
		Return([]byte(res), fmt.Errorf("Command failed with exit code 1"))

	err := p.audit()
	if err == nil {
		t.Error("audit: critical should error when there is no package lock file")
	}
}

func TestPlugin_audit_CriticalTriggerCritical(t *testing.T) {
	c := &Config{
		AuditLevel: Critical,
	}
	p, mock, _ := createTestPlugin(t, c)
	res := `{
		"metadata": {
			"vulnerabilities": {
				"info": 0,
				"low": 0,
				"moderate": 0,
				"high": 0,
				"critical": 1
			}
		}
	}`

	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"audit", "--production", "--audit-level=critical"})).
		Return([]byte(res), fmt.Errorf("Command failed with exit code 1"))

	err := p.audit()
	if err == nil {
		t.Error("audit: critical should error when critical is found")
	}
}

func TestPlugin_audit_CriticalTriggerHigh(t *testing.T) {
	c := &Config{
		AuditLevel: High,
	}
	p, mock, _ := createTestPlugin(t, c)
	res := `{
		"metadata": {
			"vulnerabilities": {
				"info": 0,
				"low": 0,
				"moderate": 0,
				"high": 0,
				"critical": 1
			}
		}
	}`

	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"audit", "--production", "--audit-level=high"})).
		Return([]byte(res), fmt.Errorf("Command failed with exit code 1"))

	err := p.audit()
	if err == nil {
		t.Error("audit: high should error when critical is found")
	}
}

func TestPlugin_Publish(t *testing.T) {
	p, mock, _ := createTestPlugin(t, &Config{
		Registry: "http://registry.test.com",
	})
	res := ` {
		"name": "@go-vela/vela-npm",
		"version": "1.0.0",
		"files": [
			{
				"path": "README.md",
				"size": 67,
				"mode": 420
			},
			{
				"path": "index.js",
				"size": 80,
				"mode": 420
			},
			{
				"path": "package.json",
				"size": 392,
				"mode": 420
			}
		]
	}`

	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"publish", "--quiet", "--registry", "http://registry.test.com"})).
		Return([]byte(res), nil)

	err := p.publish()
	if err != nil {
		t.Error(err)
	}
}

func TestPlugin_Publish_Access(t *testing.T) {
	p, mock, _ := createTestPlugin(t, &Config{
		Access:   "public",
		Registry: "http://registry.test.com",
	})

	res := ` {
		"name": "@go-vela/vela-npm",
		"version": "1.0.0",
		"files": [
			{
				"path": "README.md",
				"size": 67,
				"mode": 420
			},
			{
				"path": "index.js",
				"size": 80,
				"mode": 420
			},
			{
				"path": "package.json",
				"size": 392,
				"mode": 420
			}
		]
	}`

	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"publish", "--quiet", "--access", "public", "--registry", "http://registry.test.com"})).
		Return([]byte(res), nil)

	err := p.publish()
	if err != nil {
		t.Error(err)
	}
}

func TestPlugin_Publish_DryRun(t *testing.T) {
	c := &Config{
		DryRun:   true,
		Registry: "http://registry.test.com",
	}
	p, mock, _ := createTestPlugin(t, c)
	res := ` {
		"name": "@go-vela/vela-npm",
		"version": "1.0.0",
		"files": [
			{
				"path": "README.md",
				"size": 67,
				"mode": 420
			},
			{
				"path": "index.js",
				"size": 80,
				"mode": 420
			},
			{
				"path": "package.json",
				"size": 392,
				"mode": 420
			}
		]
	}`

	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"publish", "--quiet", "--dry-run", "--registry", "http://registry.test.com"})).
		Return([]byte(res), nil)

	err := p.publish()
	if err != nil {
		t.Error(err)
	}
}

func TestPlugin_Publish_Tag(t *testing.T) {
	c := &Config{
		Tag:      "beta",
		Registry: "http://registry.test.com",
	}
	p, mock, _ := createTestPlugin(t, c)
	res := ` {
		"name": "@go-vela/vela-npm",
		"version": "1.0.0",
		"files": [
			{
				"path": "README.md",
				"size": 67,
				"mode": 420
			},
			{
				"path": "index.js",
				"size": 80,
				"mode": 420
			},
			{
				"path": "package.json",
				"size": 392,
				"mode": 420
			}
		]
	}`

	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"publish", "--quiet", "--tag", "beta", "--registry", "http://registry.test.com"})).
		Return([]byte(res), nil)

	err := p.publish()
	if err != nil {
		t.Error(err)
	}
}

func TestPlugin_Publish_Workspaces(t *testing.T) {
	c := &Config{
		Workspaces: true,
		Registry:   "http://registry.test.com",
	}
	p, mock, _ := createTestPlugin(t, c)
	res := `{
		"@vela-npm/1": {
			"id": "@vela-npm/1@1.0.0",
			"name": "@vela-npm/1",
			"version": "1.0.0",
			"files": [
				{
					"path": "README.md",
					"size": 68,
					"mode": 420
				},
				{
					"path": "index.js",
					"size": 80,
					"mode": 420
				},
				{
					"path": "package.json",
					"size": 386,
					"mode": 420
				}
			]
		},
		"@vela-npm/2": {
			"id": "@vela-npm/2@1.0.0",
			"name": "@vela-npm/2",
			"version": "1.0.0",
			"files": [
				{
					"path": "README.md",
					"size": 68,
					"mode": 420
				},
				{
					"path": "index.js",
					"size": 80,
					"mode": 420
				},
				{
					"path": "package.json",
					"size": 386,
					"mode": 420
				}
			]
		}
	}`

	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"publish", "--quiet", "--workspaces", "--registry", "http://registry.test.com"})).
		Return([]byte(res), nil)

	err := p.publish()
	if err != nil {
		t.Error(err)
	}
}

func TestPlugin_Publish_Workspace(t *testing.T) {
	c := &Config{
		Workspace: "example",
		Registry:  "http://registry.test.com",
	}
	p, mock, _ := createTestPlugin(t, c)
	res := `{
		"@vela-npm/1": {
			"id": "@vela-npm/1@1.0.0",
			"name": "@vela-npm/1",
			"version": "1.0.0",
			"files": [
				{
					"path": "README.md",
					"size": 68,
					"mode": 420
				},
				{
					"path": "index.js",
					"size": 80,
					"mode": 420
				},
				{
					"path": "package.json",
					"size": 386,
					"mode": 420
				}
			]
		}
	}`

	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"publish", "--quiet", "--workspace", "example", "--registry", "http://registry.test.com"})).
		Return([]byte(res), nil)

	err := p.publish()
	if err != nil {
		t.Error(err)
	}
}

func TestPlugin_Publish_All(t *testing.T) {
	c := &Config{
		DryRun:     true,
		Tag:        "beta",
		Workspaces: true,
		Registry:   "http://registry.test.com",
	}
	p, mock, _ := createTestPlugin(t, c)
	res := `{
		"@vela-npm/1": {
			"id": "@vela-npm/1@1.0.0",
			"name": "@vela-npm/1",
			"version": "1.0.0",
			"files": [
				{
					"path": "README.md",
					"size": 68,
					"mode": 420
				},
				{
					"path": "index.js",
					"size": 80,
					"mode": 420
				},
				{
					"path": "package.json",
					"size": 386,
					"mode": 420
				}
			]
		},
		"@vela-npm/2": {
			"id": "@vela-npm/2@1.0.0",
			"name": "@vela-npm/2",
			"version": "1.0.0",
			"files": [
				{
					"path": "README.md",
					"size": 68,
					"mode": 420
				},
				{
					"path": "index.js",
					"size": 80,
					"mode": 420
				},
				{
					"path": "package.json",
					"size": 386,
					"mode": 420
				}
			]
		}
	}`

	mock.
		EXPECT().
		RunCommandBytes(gomock.Eq("npm"), gomock.Eq([]string{"publish", "--quiet", "--dry-run", "--tag", "beta", "--workspaces", "--registry", "http://registry.test.com"})).
		Return([]byte(res), nil)

	err := p.publish()
	if err != nil {
		t.Error(err)
	}
}
