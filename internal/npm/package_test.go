// SPDX-License-Identifier: Apache-2.0
package npm

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestPackage_Validate_NoRegistry(t *testing.T) {
	p := &packageJSON{
		Name:    "test-package",
		Version: "1.0.0",
		PublishConfig: publishConfig{
			Registry: "",
		},
	}
	err := p.Validate("")

	logrus.Warn(err)

	if err != nil {
		t.Fail()
	}
}

func TestPackage_Validate_NoName(t *testing.T) {
	p := &packageJSON{
		Version: "1.0.0",
		PublishConfig: publishConfig{
			Registry: "",
		},
	}
	err := p.Validate("")

	logrus.Warn(err)

	if err == nil {
		t.Fail()
	}
}

func TestPackage_Validate_NoVersion(t *testing.T) {
	p := &packageJSON{
		Name: "test-package",
		PublishConfig: publishConfig{
			Registry: "",
		},
	}
	err := p.Validate("")

	logrus.Warn(err)

	if err == nil {
		t.Fail()
	}
}

func TestPackage_Validate_BadVersion(t *testing.T) {
	p := &packageJSON{
		Name:    "test-package",
		Version: "beta",
		PublishConfig: publishConfig{
			Registry: "",
		},
	}
	err := p.Validate("")

	logrus.Warn(err)

	if err == nil {
		t.Fail()
	}
}

func TestPackage_Validate_RegistryMismatch(t *testing.T) {
	p := &packageJSON{
		Name:    "test-package",
		Version: "1.0.0",
		PublishConfig: publishConfig{
			Registry: "someRegistry",
		},
	}
	err := p.Validate("otherRegistry")

	logrus.Warn(err)

	if err == nil {
		t.Fail()
	}
}
