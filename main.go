package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/xlog"
	"golang.org/x/net/webdav"
)

func main() {
	var (
		homeDir, logLevel, login, password string
		port                               uint
	)

	flag.StringVar(&homeDir, "home_dir", ".", "home dir")
	flag.StringVar(&logLevel, "log_level", "info", "Log level (debug | info | warn | error | fatal)")
	flag.StringVar(&login, "login", "", "login for basic authorization")
	flag.StringVar(&password, "password", "", "password for basic authorization")
	flag.UintVar(&port, "port", 8080, "server port")

	flag.Parse()

	lvl, err := xlog.LevelFromString(logLevel)
	if err != nil {
		log.Fatal(err)
	}

	logger := xlog.New(xlog.Config{
		Level:  lvl,
		Output: xlog.NewConsoleOutput(),
	})
	log.SetOutput(logger)

	h := &webdav.Handler{
		FileSystem: webdav.Dir(homeDir),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			logger.Debugf("Method: \"%s\", Path: \"%s\"", r.Method, r.URL)
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if login != "" {
			l, p, b := r.BasicAuth()
			if l != login || p != password || !b {
				if b {
					logger.Warnf("No authorization, login: \"%s\", password: \"%s\"", l, p)
				}

				w.Header().Set("WWW-Authenticate", `Basic realm="WEBDAV"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				return
			}
		}

		h.ServeHTTP(w, r)

	})

	logger.Info("Start WEBDAV server")
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
