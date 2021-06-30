package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stderr)
	log.SetLevel(logrus.InfoLevel)
	//log.SetLevel(logrus.DebugLevel)
}

func main() {
	filename := "testdata/package-large.json"
	pkg, err := ReadPackageJSON(filename)
	if err != nil {
		log.Fatalf("reading package.json: %v", err)
	}

	log.Infof("Dependencies")
	runSuggestions(pkg.Dependencies)

	if len(pkg.DevDependencies) > 0 {
		log.Infof("Development Dependencies")
		runSuggestions(pkg.DevDependencies)
	}

	if len(pkg.PeerDependencies) > 0 {
		log.Infof("Peer Dependencies")
		runSuggestions(pkg.PeerDependencies)
	}

	if len(pkg.OptionalDependencies) > 0 {
		log.Infof("Optional Dependencies")
		runSuggestions(pkg.OptionalDependencies)
	}
}

func runSuggestions(deps Dependencies) {
	var names []string
	for name := range deps {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		log.Debugf("Processing dependency %q...", name)

		constraintStr := deps[name]

		if strings.HasPrefix(constraintStr, "npm:") {
			// resolve aliases
			rest := strings.TrimPrefix(constraintStr, "npm:")
			parts := strings.Split(rest, "@")
			if len(parts) != 2 {
				log.Warnf("Could not resolve alias %q.", constraintStr)
				continue
			}

			name = parts[0]
			constraintStr = parts[1]
		}

		constraint, err := semver.NewConstraint(constraintStr)
		if err != nil {
			log.Warnf("Cannot parse version constraint %q: %v", constraintStr, err)
			continue
		}

		pkg, err := GetPackage(name)
		if err != nil {
			log.Errorf("Could not get information about package %q: %v", name, err)
			continue
		}

		candidates := ExtractVersions(pkg.Versions)
		current := NewestMatching(constraint, candidates)
		newest := Newest(candidates)

		if current != nil && !newest.GreaterThan(current) {
			continue
		}

		fmt.Printf("%-40s %16s  →  %16s\n", name, constraint, newest)
	}
}
