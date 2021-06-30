package main

import (
	"sort"

	"github.com/Masterminds/semver/v3"
)

func ExtractVersions(versionMap map[string]Version) semver.Collection {
	var versions semver.Collection
	for _, versionObj := range versionMap {
		version, err := semver.NewVersion(versionObj.Version)

		if err != nil {
			log.Warnf("Cannot parse version %q: %v", versionObj.Version, err)
		}
		versions = append(versions, version)
	}
	return versions
}

func sortMostRecentFirst(c *semver.Collection) {
	sort.Sort(sort.Reverse(*c))
}

func NewestMatching(constraint *semver.Constraints, candidates semver.Collection) *semver.Version {
	sortMostRecentFirst(&candidates)

	for _, candidate := range candidates {
		if constraint.Check(candidate) {
			return candidate
		}
	}
	return nil
}

func Newest(candidates semver.Collection) *semver.Version {
	sortMostRecentFirst(&candidates)

	for _, candidate := range candidates {
		// exclude prerelease versions
		if candidate.Prerelease() == "" {
			return candidate
		}
	}

	return nil
}

const (
	MajorDowngrade      = -4
	MinorDowngrade      = -3
	PatchDowngrade      = -2
	PrereleaseDowngrade = -1
	Equal               = 0
	PrereleaseUpgrade   = 1
	PatchUpgrade        = 2
	MinorUpgrade        = 3
	MajorUpgrade        = 4
)

func Compare(from, to *semver.Version) int {
	if from == nil {
		return MajorUpgrade
	}

	if from.Major() < to.Major() {
		return MajorUpgrade
	}
	if from.Major() > to.Major() {
		return MajorDowngrade
	}
	if from.Minor() < to.Minor() {
		return MinorUpgrade
	}
	if from.Minor() > to.Minor() {
		return MinorDowngrade
	}
	if from.Patch() < to.Patch() {
		return PatchUpgrade
	}
	if from.Patch() > to.Patch() {
		return PatchDowngrade
	}

	return to.Compare(from)
}
