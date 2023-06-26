// Fortio CLI/Main utilities.
//
// (c) 2023 Fortio Authors
// See LICENSE

// Package scli extends [cli] for server main()s
// [ServerMain] allows the setup of a confimap/directory watch for flags
// and a config endpoint (uses [fortio.org/dflag]).
// Configure using the [cli] package variables (at minimum [MinArgs] unless your
// binary only accepts flags), setup additional [flag] before calling
// [ServerMain].
// It also includes [NumFD] utility function to cross platform get the number
// of open file descriptors (handles on windows) held by your go process.
package scli // import "fortio.org/scli"

import (
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"fortio.org/cli"
	"fortio.org/dflag"
	"fortio.org/dflag/configmap"
	"fortio.org/dflag/dynloglevel"
	"fortio.org/dflag/endpoint"
	"fortio.org/log"
	"fortio.org/version"
)

// NormalizePort parses port and returns host:port if port is in the form
// of host:port already or :port if port is only a port (doesn't contain :).
// Copied from fortio.org/fnet.NormalizePort to avoid dependency loop.
func NormalizePort(port string) string {
	if strings.ContainsAny(port, ":") {
		return port
	}
	return ":" + port
}

// ServerMain extends [cli.Main] and returns true if a config port server has been started
// caller needs to select {} after their own code is ready.
// [cli.ExitFunction] will have been called (ie program will have exited exited)
// if there are usage errors (wrong number of arguments, bad flags etc...).
// It sets up (optional) config-dir to watch and listen on config-port for dynamic flag
// changes and UI/api.
func ServerMain() bool {
	configDir := flag.String("config-dir", "", "Config `directory` to watch for dynamic flag changes")
	configPort := flag.String("config-port", "", "Config `port` to open for dynamic flag UI/api")
	dynloglevel.LoggerFlagSetup("loglevel")
	dflag.DynBool(flag.CommandLine, "logger-json", true,
		"Log in JSON format, use -logger-json=false to disable").WithSyncNotifier(func(_ bool, newValue bool) {
		log.Debugf("Changing log format to JSON %v", newValue)
		log.Config.JSON = newValue
	})
	dflag.DynBool(flag.CommandLine, "logger-timestamp", true,
		"Timestamps emitted in JSON logs, use -logger-timestamp=false to disable").WithSyncNotifier(func(_ bool, newValue bool) {
		log.Debugf("Changing log format to JSON timestamp %v", newValue)
		log.Config.NoTimestamp = !newValue
	})
	cli.ServerMode = true
	cli.Main() // will call ExitFunction() if there are usage errors
	if *configDir != "" {
		if _, err := configmap.Setup(flag.CommandLine, *configDir); err != nil {
			log.Critf("Unable to watch config/flag changes in %v: %v", *configDir, err)
		}
	}
	shortScliV, _, _ := version.FromBuildInfoPath("fortio.org/scli")

	hasStartedServer := false
	if *configPort != "" {
		// Sort of inlining fortio.org/fhttp.HTTPServer here to avoid
		// a dependency loop.
		port := NormalizePort(*configPort)
		m := http.NewServeMux()
		s := &http.Server{
			Addr:        port,
			Handler:     m,
			ReadTimeout: 3 * time.Second,
		}
		setURL := "/set"
		ep := endpoint.NewFlagsEndpoint(flag.CommandLine, setURL)
		m.HandleFunc("/", ep.ListFlags)
		m.HandleFunc(setURL, ep.SetFlag)
		ln, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("Unable to serve config on %s: %v", s.Addr, err)
		}
		log.Infof("Fortio scli %v dflag config server listening on %s", shortScliV, ln.Addr())
		go func() {
			err := s.Serve(ln)
			if err != nil {
				log.Fatalf("Unable to serve config on %s: %v", s.Addr, err)
			}
		}()
		hasStartedServer = true
	}
	log.Infof("Starting %s %s", cli.ProgramName, cli.LongVersion)
	return hasStartedServer
}

// UntilInterrupted runs forever or until interrupted (ctrl-c or shutdown signal (kill -INT)).
func UntilInterrupted() {
	// listen for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// Block until a signal is received.
	<-c
	log.Warnf("Interrupt received.")
}
