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
	"runtime"
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

var ConfigMapUpdater *configmap.Updater // Expose the updater if advanced callers want to check warnings or stop it etc.

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
	dflag.DynBool(flag.CommandLine, "logger-file-line", true,
		"Filename and line numbers emitted in JSON logs, use -logger-file-line=false to disable").WithSyncNotifier(
		func(_ bool, newValue bool) {
			log.Debugf("Changing log format to JSON file and line %v", newValue)
			log.Config.LogFileAndLine = newValue
		})
	dflag.DynBool(flag.CommandLine, "logger-goroutine", true,
		"GoroutineID emitted in JSON/color logs, use -logger-goroutine=false to disable").WithSyncNotifier(
		func(_ bool, newValue bool) {
			log.Debugf("Changing log format to GoroutineID %v", newValue)
			log.Config.GoroutineID = newValue
		})
	cli.ServerMode = true
	cli.Main() // will call ExitFunction() if there are usage errors
	if *configDir != "" {
		var err error
		ConfigMapUpdater, err = configmap.Setup(flag.CommandLine, *configDir)
		if err != nil {
			log.Critf("Unable to watch config/flag changes in %v: %v", *configDir, err)
		} else if ConfigMapUpdater.Warnings() != 0 {
			log.S(log.Warning, "Unknown flags found", log.Int("count", ConfigMapUpdater.Warnings()), log.Str("dir", *configDir))
		}
	}

	// So http client library for instance ends up logging in JSON or color too and not break json parsing.
	log.InterceptStandardLogger(log.Warning)

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
		log.S(log.Info, "Fortio scli dflag config server listening", log.Str("version", shortScliV), log.Attr("addr", ln.Addr()))
		go func() {
			err := s.Serve(ln)
			if err != nil {
				log.Fatalf("Unable to serve config on %s: %v", s.Addr, err)
			}
		}()
		hasStartedServer = true
	}
	log.S(log.Info, "Starting", log.Str("command", cli.ProgramName), log.Str("version", cli.LongVersion),
		log.Int("go-max-procs", runtime.GOMAXPROCS(0)))
	return hasStartedServer
}

// UntilInterrupted runs forever or until interrupted (ctrl-c or shutdown signal (kill -INT or -TERM)).
// Kubernetes for instance sends a SIGTERM before killing a pod.
// You can place your clean shutdown code after this call in the main().
// UntilInterrupted forwards to [cli.UntilInterrupted], call that one directly in newer code.
func UntilInterrupted() {
	cli.UntilInterrupted()
}
