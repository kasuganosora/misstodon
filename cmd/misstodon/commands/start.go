package commands

import (
	_ "embed"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api"
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/acme/autocert"
)

//go:embed banner.txt
var banner string

var Start = &cli.Command{
	Name:  "start",
	Usage: "Start the server",
	Before: func(c *cli.Context) error {
		appVersion := global.AppVersion
		if !c.Bool("no-color") {
			appVersion = "\033[1;31;40m" + appVersion + "\033[0m"
		}
		fmt.Printf("\n%s  %s\n\n", banner, appVersion)
		return nil
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "bind",
			Aliases: []string{"b"},
			Usage:   "bind address",
		},
		&cli.StringFlag{
			Name: "fallback-server",
			Usage: "if proxy-server is not found in the request, the fallback server address will be used, " +
				`e.g. "misskey.io"`,
		},
	},
	Action: func(c *cli.Context) error {
		conf := global.Config
		if c.IsSet("fallbackServer") {
			conf.Proxy.FallbackServer = c.String("fallbackServer")
		}
		bindAddress, _ := utils.StrEvaluation(c.String("bind"), conf.Server.BindAddress)

		gin.SetMode(gin.ReleaseMode)
		r := gin.New()

		api.Router(r)

		logStart := log.Info().Str("address", bindAddress)
		switch {
		case conf.Server.AutoTLS && conf.Server.Domain != "":
			cacheDir, _ := filepath.Abs("./cert/.cache")
			tlsManager := &autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(conf.Server.Domain),
				Cache:      autocert.DirCache(cacheDir),
			}
			logStart.Msg("Starting server with AutoTLS")
			server := &http.Server{
				Addr:      ":https",
				TLSConfig: tlsManager.TLSConfig(),
				Handler:   r,
			}
			go func() {
				_ = http.ListenAndServe(":http", tlsManager.HTTPHandler(nil))
			}()
			return server.ListenAndServeTLS("", "")
		case conf.Server.TlsCertFile != "" && conf.Server.TlsKeyFile != "":
			logStart.Msg("Starting server with TLS")
			return r.RunTLS(bindAddress, conf.Server.TlsCertFile, conf.Server.TlsKeyFile)
		default:
			logStart.Msg("Starting server")
			return r.Run(bindAddress)
		}
	},
}
