// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.
package npm

import "testing"

func TestConfig_Validate_Valid(t *testing.T) {
	c := &Config{
		UserName: "testuser",
	}
	p, _, _ := createTestPlugin(t, c)

	err := p.Validate()
	if err != nil {
		t.Error(err)
	}
}

func TestConfig_Validate_NoOS(t *testing.T) {
	c := &Config{
		UserName: "testuser",
	}
	p, _, _ := createTestPlugin(t, c)

	err := p.Validate()
	if err != nil {
		t.Error(err)
	}
}

func TestConfig_Validate_NoUserName(t *testing.T) {
	c := &Config{}
	p, _, _ := createTestPlugin(t, c)

	err := p.Validate()
	if err == nil {
		t.Fail()
	}
}

func TestConfig_Validate_Token(t *testing.T) {
	c := &Config{
		Token: "token_abc",
	}
	p, _, _ := createTestPlugin(t, c)

	err := p.Validate()
	if err != nil {
		t.Fail()
	}
}

func TestConfig_Validate_BadTag(t *testing.T) {
	c := &Config{
		Tag:      "1.0.0",
		UserName: "testuser",
	}
	p, _, _ := createTestPlugin(t, c)
	err := p.Validate()

	if err == nil {
		t.Fail()
	}
}

func TestConfig_Validate_NormalizeAuditLevel_Info(t *testing.T) {
	c := &Config{
		UserName:   "testuser",
		AuditLevel: "l",
	}
	p, _, _ := createTestPlugin(t, c)

	err := p.Validate()
	if err != nil {
		t.Error(err)
	}

	if c.AuditLevel != Low {
		t.Error("AuditLevel not normalized")
	}
}

func TestConfig_Validate_NormalizeAuditLevel_Default(t *testing.T) {
	c := &Config{
		UserName:   "testuser",
		AuditLevel: "what",
	}
	p, _, _ := createTestPlugin(t, c)

	err := p.Validate()
	if err != nil {
		t.Error(err)
	}

	if c.AuditLevel != Low {
		t.Error("AuditLevel not defaulted")
	}
}

func TestConfig_Validate_Access_Public(t *testing.T) {
	c := &Config{
		UserName: "testuser",
		Access:   "public",
	}
	p, _, _ := createTestPlugin(t, c)

	err := p.Validate()
	if err != nil {
		t.Error(err)
	}
}

func TestConfig_Validate_Access_Restricted(t *testing.T) {
	c := &Config{
		UserName: "testuser",
		Access:   "restricted",
	}
	p, _, _ := createTestPlugin(t, c)

	err := p.Validate()
	if err != nil {
		t.Error(err)
	}
}
func TestConfig_Validate_Access_NotRecognized(t *testing.T) {
	c := &Config{
		UserName: "testuser",
		Access:   "protected",
	}
	p, _, _ := createTestPlugin(t, c)

	err := p.Validate()
	if err == nil {
		t.Error(err)
	}
}

func TestConfig_Validate_Workspace(t *testing.T) {
	c := &Config{
		UserName:   "testuser",
		Workspaces: true,
		Workspace:  "example",
	}
	p, _, _ := createTestPlugin(t, c)

	err := p.Validate()
	if err == nil {
		t.Error(err)
	}
}
