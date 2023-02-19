// Fortio CLI/Main utilities.
//
// (c) 2023 Fortio Authors
// See LICENSE

// Package scli extends [cli] for server main()s
// [ServerMain] allows the setup of a confimap/directory watch for flags
// and a config endpoint (uses [fortio.org/dflag]).
// Configure using the [cli] package variables (at minimum [MinArgs] unless your
// binary only accepts flags), setup additional [flag]s before calling
// [ServerMain].
package scli // import "fortio.org/scli"

import (
	"flag"

	"fortio.org/cli"
	"fortio.org/dflag/configmap"
	"fortio.org/dflag/dynloglevel"
	"fortio.org/dflag/endpoint"
	"fortio.org/fortio/fhttp"
	"fortio.org/log"
)

// ServerMain extends [cli.Main] and returns true if a config port server has been started
// caller needs to select {} after its own code is ready.
// Will have called ExitFunction (ie exited) if there are usage errors
// (wrong number of arguments, bad flags etc...).
// It sets up (optional) config-dir to watch and listen on config-port for dynamic flag
// changes and UI/api.
func ServerMain() bool {
	configDir := flag.String("config-dir", "", "Config `directory` to watch for dynamic flag changes")
	configPort := flag.String("config-port", "", "Config `port` to open for dynamic flag UI/api")
	dynloglevel.LoggerFlagSetup("loglevel")
	cli.ServerMode = true
	cli.Main() // will call ExitFunction() if there are usage errors
	if *configDir != "" {
		if _, err := configmap.Setup(flag.CommandLine, *configDir); err != nil {
			log.Critf("Unable to watch config/flag changes in %v: %v", *configDir, err)
		}
	}
	hasStartedServer := false
	if *configPort != "" {
		mux, addr := fhttp.HTTPServer("config", *configPort) // err already logged
		if addr != nil {
			hasStartedServer = true
			setURL := "/set"
			ep := endpoint.NewFlagsEndpoint(flag.CommandLine, setURL)
			mux.HandleFunc("/", ep.ListFlags)
			mux.HandleFunc(setURL, ep.SetFlag)
		}
	}
	log.Printf("Starting %s %s", cli.ProgramName, cli.LongVersion)
	return hasStartedServer
}

// Plural adds an "s" to the noun if i is not 1.
func Plural(i int, noun string) string {
	return PluralExt(i, noun, "s")
}

// PluralExt returns the noun with an extension if i is not 1.
// Eg:
//
//	PluralExt(1, "address", "es") // -> "address"
//	PluralExt(3 /* or 0 */, "address", "es") // -> "addresses"
func PluralExt(i int, noun string, ext string) string {
	if i == 1 {
		return noun
	}
	return noun + ext
}
