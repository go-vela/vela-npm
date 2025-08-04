// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/go-vela/vela-npm/internal/npm"
	"github.com/go-vela/vela-npm/version"
)

func main() { //nolint: funlen // length of main function is acceptable for CLI applications
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	// create new CLI application
	// Plugin Information
	cmd := cli.Command{
		Name:      "vela-npm",
		Usage:     "Vela npm plugin for publishing NodeJS packages",
		Copyright: "Copyright 2022 Target Brands, Inc. All rights reserved.",
		Authors: []any{
			&mail.Address{
				Name:    "Vela Admins",
				Address: "vela@target.com",
			},
		},
		Version: v.Semantic(),
		Action:  run,
	}
	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "token",
			Aliases:     []string{"t"},
			Usage:       "auth token",
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_TOKEN"),
				cli.EnvVar("PLUGIN_TOKEN"),
				cli.EnvVar("NPM_TOKEN"),
				cli.File("/vela/parameters/npm/token"),
				cli.File("/vela/secrets/npm/token"),
			),
		},
		&cli.StringFlag{
			Name:        "username",
			Aliases:     []string{"u"},
			Usage:       "name of user",
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_USERNAME"),
				cli.EnvVar("PLUGIN_USERNAME"),
				cli.EnvVar("NPM_USERNAME"),
				cli.File("/vela/parameters/npm/username"),
				cli.File("/vela/secrets/npm/username"),
				cli.File("/vela/secrets/managed-auth/username"),
			),
		},
		&cli.StringFlag{
			Name:        "password",
			Aliases:     []string{"p"},
			Usage:       "password for user",
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_PASSWORD"),
				cli.EnvVar("PLUGIN_PASSWORD"),
				cli.EnvVar("NPM_PASSWORD"),
				cli.File("/vela/parameters/npm/password"),
				cli.File("/vela/secrets/npm/password"),
				cli.File("/vela/secrets/managed-auth/password"),
			),
		},
		&cli.StringFlag{
			Name:        "registry",
			Aliases:     []string{"r"},
			Usage:       "npm registry",
			Value:       npm.DefaultRegistry,
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_REGISTRY"),
				cli.EnvVar("PLUGIN_REGISTRY"),
				cli.EnvVar("NPM_REGISTRY"),
				cli.File("/vela/parameters/npm/registry"),
				cli.File("/vela/secrets/npm/registry"),
			),
		},
		&cli.StringFlag{
			Name:        "email",
			Aliases:     []string{"e"},
			Usage:       "email for user",
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_EMAIL"),
				cli.EnvVar("PLUGIN_EMAIL"),
				cli.EnvVar("NPM_EMAIL"),
				cli.File("/vela/parameters/npm/email"),
				cli.File("/vela/secrets/npm/email"),
			),
		},
		&cli.BoolFlag{
			Name:        "strict-ssl",
			Usage:       "enables strict SSL",
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_STRICT_SSL"),
				cli.EnvVar("PLUGIN_STRICT_SSL"),
				cli.EnvVar("STRICT_SSL"),
				cli.File("/vela/parameters/npm/strict_ssl"),
				cli.File("/vela/secrets/npm/strict_ssl"),
			),
		},
		&cli.BoolFlag{
			Name:        "always-auth",
			Usage:       "enables always auth",
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_ALWAYS_AUTH"),
				cli.EnvVar("PLUGIN_ALWAYS_AUTH"),
				cli.EnvVar("ALWAYS_AUTH"),
				cli.File("/vela/parameters/npm/always_auth"),
				cli.File("/vela/secrets/npm/always_auth"),
			),
		},
		&cli.BoolFlag{
			Name:        "skip-ping",
			Usage:       "skips auth ping",
			Value:       false,
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_SKIP_PING"),
				cli.EnvVar("PLUGIN_SKIP_PING"),
				cli.EnvVar("SKIP_PING"),
				cli.File("/vela/parameters/npm/skip_ping"),
				cli.File("/vela/secrets/npm/skip_ping"),
			),
		},
		&cli.BoolFlag{
			Name:        "first-publish",
			Usage:       "(DEPRECATED): skips version lookup and verification for first time publishes",
			DefaultText: "N/A",
		},
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:       "info",
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_LOG"),
				cli.EnvVar("PARAMETER_LOG_LEVEL"),
				cli.EnvVar("PLUGIN_LOG"),
				cli.EnvVar("PLUGIN_LOG_LEVEL"),
				cli.EnvVar("LOG_LEVEL"),
				cli.EnvVar("LOG"),
				cli.File("/vela/parameters/npm/log_level"),
				cli.File("/vela/secrets/npm/log_level"),
			),
		},
		&cli.BoolFlag{
			Name:        "dry-run",
			Usage:       "publish command will only do dry run",
			Value:       false,
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_DRY_RUN"),
				cli.EnvVar("PLUGIN_DRY_RUN"),
				cli.EnvVar("DRY_RUN"),
				cli.File("/vela/parameters/npm/dry_run"),
				cli.File("/vela/secrets/npm/dry_run"),
			),
		},
		&cli.StringFlag{
			Name:        "tag",
			Usage:       "publish package with given tag",
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_TAG"),
				cli.EnvVar("PLUGIN_TAG"),
				cli.EnvVar("TAG"),
				cli.File("/vela/parameters/npm/tag"),
				cli.File("/vela/secrets/npm/tag"),
			),
		},
		&cli.StringFlag{
			Name:        "audit-level",
			Usage:       "The level at which an npm audit will fail - options: (none|low|moderate|high|critical)",
			Value:       "none",
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_AUDIT_LEVEL"),
				cli.EnvVar("PARAMETER_AUDIT"),
				cli.EnvVar("PLUGIN_AUDIT_LEVEL"),
				cli.EnvVar("PLUGIN_AUDIT"),
				cli.EnvVar("AUDIT_LEVEL"),
				cli.EnvVar("AUDIT"),
				cli.File("/vela/parameters/npm/audit_level"),
				cli.File("/vela/secrets/npm/audit_level"),
			),
		},
		&cli.StringFlag{
			Name:        "access",
			Usage:       "Tells the registry whether this package should be published as public or restricted. Only applies to scoped packages, which default to restricted",
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_ACCESS"),
				cli.EnvVar("PARAMETER_ACCESS"),
				cli.EnvVar("ACCESS"),
				cli.File("/vela/parameters/npm/access"),
				cli.File("/vela/secrets/npm/access"),
			),
		},
		&cli.StringFlag{
			Name:        "ci",
			Usage:       "set to CI environment",
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI"),
				cli.File("/vela/parameters/npm/ci"),
				cli.File("/vela/secrets/npm/ci"),
			),
		},
		&cli.BoolFlag{
			Name:        "workspaces",
			Usage:       "publish all workspaces",
			Value:       false,
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_WORKSPACES"),
				cli.EnvVar("PLUGIN_WORKSPACES"),
				cli.EnvVar("WORKSPACES"),
				cli.EnvVar("WS"),
				cli.File("/vela/parameters/npm/workspaces"),
				cli.File("/vela/secrets/npm/workspaces"),
			),
		},
		&cli.StringFlag{
			Name:        "workspace",
			Usage:       "publish a specific workspace",
			DefaultText: "N/A",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_WORKSPACE"),
				cli.EnvVar("PLUGIN_WORKSPACE"),
				cli.EnvVar("WORKSPACE"),
				cli.EnvVar("W"),
				cli.File("/vela/parameters/npm/workspace"),
				cli.File("/vela/secrets/npm/workspace"),
			),
		},
	}

	if err = cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(_ context.Context, c *cli.Command) error {
	// set the log level for the plugin
	switch strings.ToLower(c.String("log-level")) {
	case "t", "trace":
		log.SetLevel(log.TraceLevel)
	case "d", "debug":
		log.SetLevel(log.DebugLevel)
	case "w", "warn":
		log.SetLevel(log.WarnLevel)
	case "e", "error":
		log.SetLevel(log.ErrorLevel)
	case "f", "fatal":
		log.SetLevel(log.FatalLevel)
	case "p", "panic":
		log.SetLevel(log.PanicLevel)
	case "i", "info":
		fallthrough
	default:
		log.SetLevel(log.InfoLevel)
	}

	if c.IsSet("ci") {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	} else {
		log.SetFormatter(&log.TextFormatter{
			ForceColors:   true,
			FullTimestamp: false,
			PadLevelText:  true,
		})
	}

	// docs reference
	log.WithFields(log.Fields{
		"code":     "https://github.com/go-vela/vela-npm",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/pipeline/npm/",
		"registry": "https://hub.docker.com/r/target/vela-npm",
		"version":  "1.0.0",
	}).Info("Vela NPM Plugin")

	config := &npm.Config{
		Token:           c.String("token"),
		UserName:        c.String("username"),
		Password:        c.String("password"),
		Registry:        c.String("registry"),
		Email:           c.String("email"),
		StrictSSL:       c.Bool("strict-ssl"),
		IsStrictSSLSet:  c.IsSet("strict-ssl"),
		AlwaysAuth:      c.Bool("always-auth"),
		IsAlwaysAuthSet: c.IsSet("always-auth"),
		SkipPing:        c.Bool("skip-ping"),
		DryRun:          c.Bool("dry-run"),
		Tag:             c.String("tag"),
		AuditLevel:      c.String("audit-level"),
		Access:          c.String("access"),
		Workspaces:      c.Bool("workspaces"),
		Workspace:       c.String("workspace"),
	}

	p := npm.NewPlugin(config)

	// validate plugin inputs
	if err := p.Validate(); err != nil {
		return err
	}

	// run the plugin
	if err := p.Exec(); err != nil {
		return err
	}

	return nil
}
