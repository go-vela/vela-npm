// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-vela/vela-npm/internal/npm"
	"github.com/go-vela/vela-npm/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	app := &cli.App{
		Name:      "vela-npm",
		Usage:     "Vela npm plugin for publishing NodeJS packages",
		Version:   v.Semantic(),
		HelpName:  "vela-npm",
		Copyright: "Copyright 2022 Target Brands, Inc. All rights reserved.",
		Authors: []*cli.Author{
			{
				Name:  "Vela Admins",
				Email: "vela@target.com",
			},
		},
		Action: run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "token",
				Aliases:     []string{"t"},
				Usage:       "auth token",
				EnvVars:     []string{"PARAMETER_TOKEN", "PLUGIN_TOKEN", "NPM_TOKEN"},
				FilePath:    string("/vela/parameters/npm/token,/vela/secrets/npm/token"),
				DefaultText: "N/A",
			},
			&cli.StringFlag{
				Name:        "username",
				Aliases:     []string{"u"},
				Usage:       "name of user",
				EnvVars:     []string{"PARAMETER_USERNAME", "PLUGIN_USERNAME", "NPM_USERNAME"},
				FilePath:    string("/vela/parameters/npm/username,/vela/secrets/npm/username"),
				DefaultText: "N/A",
			},
			&cli.StringFlag{
				Name:        "password",
				Aliases:     []string{"p"},
				Usage:       "password for user",
				EnvVars:     []string{"PARAMETER_PASSWORD", "PLUGIN_PASSWORD", "NPM_PASSWORD"},
				FilePath:    string("/vela/parameters/npm/password,/vela/secrets/npm/password"),
				DefaultText: "N/A",
			},
			&cli.StringFlag{
				Name:        "registry",
				Aliases:     []string{"r"},
				Usage:       "npm registry",
				EnvVars:     []string{"PARAMETER_REGISTRY", "PLUGIN_REGISTRY", "NPM_REGISTRY"},
				FilePath:    string("/vela/parameters/npm/registry,/vela/secrets/npm/registry"),
				Value:       npm.DefaultRegistry,
				DefaultText: "N/A",
			},
			&cli.StringFlag{
				Name:        "email",
				Aliases:     []string{"e"},
				Usage:       "email for user",
				EnvVars:     []string{"PARAMETER_EMAIL", "PLUGIN_EMAIL", "NPM_EMAIL"},
				FilePath:    string("/vela/parameters/npm/email,/vela/secrets/npm/email"),
				DefaultText: "N/A",
			},
			&cli.BoolFlag{
				Name:        "strict-ssl",
				Usage:       "enables strict SSL",
				EnvVars:     []string{"PARAMETER_STRICT_SSL", "PLUGIN_STRICT_SSL", "STRICT_SSL"},
				FilePath:    string("/vela/parameters/npm/strict_ssl,/vela/secrets/npm/strict_ssl"),
				DefaultText: "N/A",
			},
			&cli.BoolFlag{
				Name:        "always-auth",
				Usage:       "enables always auth",
				EnvVars:     []string{"PARAMETER_ALWAYS_AUTH", "PLUGIN_ALWAYS_AUTH", "ALWAYS_AUTH"},
				FilePath:    string("/vela/parameters/npm/always_auth,/vela/secrets/npm/always_auth"),
				DefaultText: "N/A",
			},
			&cli.BoolFlag{
				Name:        "skip-ping",
				Usage:       "skips auth ping",
				Value:       false,
				EnvVars:     []string{"PARAMETER_SKIP_PING", "PLUGIN_SKIP_PING", "SKIP_PING"},
				FilePath:    string("/vela/parameters/npm/skip_ping,/vela/secrets/npm/skip_ping"),
				DefaultText: "N/A",
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
				EnvVars:     []string{"PARAMETER_LOG", "PARAMETER_LOG_LEVEL", "PLUGIN_LOG", "PLUGIN_LOG_LEVEL", "LOG_LEVEL", "LOG"},
				FilePath:    string("/vela/parameters/npm/log_level,/vela/secrets/npm/log_level"),
				DefaultText: "N/A",
			},
			&cli.BoolFlag{
				Name:        "dry-run",
				Usage:       "publish command will only do dry run",
				Value:       false,
				EnvVars:     []string{"PARAMETER_DRY_RUN", "PLUGIN_DRY_RUN", "DRY_RUN"},
				FilePath:    string("/vela/parameters/npm/dry_run,/vela/secrets/npm/dry_run"),
				DefaultText: "N/A",
			},
			&cli.StringFlag{
				Name:        "tag",
				Usage:       "publish package with given tag",
				EnvVars:     []string{"PARAMETER_TAG", "PLUGIN_TAG", "TAG"},
				FilePath:    string("/vela/parameters/npm/tag,/vela/secrets/npm/tag"),
				DefaultText: "N/A",
			},
			&cli.StringFlag{
				Name:        "audit-level",
				Usage:       "The level at which an npm audit will fail - options: (none|low|moderate|high|critical)",
				Value:       "low",
				EnvVars:     []string{"PARAMETER_AUDIT_LEVEL", "PARAMETER_AUDIT", "PLUGIN_AUDIT_LEVEL", "PLUGIN_AUDIT", "AUDIT_LEVEL", "AUDIT"},
				FilePath:    string("/vela/parameters/npm/audit_level,/vela/secrets/npm/audit_level"),
				DefaultText: "N/A",
			},
			&cli.StringFlag{
				Name:        "access",
				Usage:       "Tells the registry whether this package should be published as public or restricted. Only applies to scoped packages, which default to restricted",
				EnvVars:     []string{"PARAMETER_ACCESS", "PARAMETER_ACCESS", "ACCESS"},
				FilePath:    string("/vela/parameters/npm/access,/vela/secrets/npm/access"),
				DefaultText: "N/A",
			},
			&cli.StringFlag{
				Name:        "ci",
				Usage:       "set to CI environment",
				EnvVars:     []string{"CI"},
				FilePath:    string("/vela/parameters/npm/ci,/vela/secrets/npm/ci"),
				DefaultText: "N/A",
			},
			&cli.BoolFlag{
				Name:        "workspaces",
				Usage:       "publish all workspaces",
				Value:       false,
				EnvVars:     []string{"PARAMETER_WORKSPACES", "PLUGIN_WORKSPACES", "WORKSPACES", "WS"},
				FilePath:    string("/vela/parameters/npm/workspaces,/vela/secrets/npm/workspaces"),
				DefaultText: "N/A",
			},
			&cli.StringFlag{
				Name:        "workspace",
				Usage:       "publish a specific workspace",
				EnvVars:     []string{"PARAMETER_WORKSPACE", "PLUGIN_WORKSPACE", "WORKSPACE", "W"},
				FilePath:    string("/vela/parameters/npm/workspace,/vela/secrets/npm/workspace"),
				DefaultText: "N/A",
			},
		},
	}

	if err = app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
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
