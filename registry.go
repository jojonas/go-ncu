package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
)

var packageCacheFilename = ".packages"
var packageCacheExpiration = 3 * time.Hour
var packageCacheCleanupInterval = 5 * time.Minute
var packageCache *cache.Cache

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

type Repository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Version struct {
	Name    string `json:"name"`
	Version string `json:"version"`

	Dependencies    Dependencies `json:"dependencies"`
	DevDependencies Dependencies `json:"devDependencies"`

	// Homepage        string            `json:"homepage"`
	// Repository      Repository        `json:"repository"`
	// Scripts         map[string]string `json:"scripts"`
	// Author          User              `json:"author"`
	// License         string            `json:"license"`
	// Readme          string            `json:"readme"`
	// ReadmeFilename  string            `json:"readmeFilename"`
	// Id              string            `json:"_id"`
	// Description     string            `json:"description"`
	// Dist            struct {
	//	SHASum  string `json:"shasum"`
	//	Tarball string `json:"tarball"`
	// } `json:"dist"`
	// NPMVersion  string `json:"_npmVersion"`
	// NPMUser     User   `json:"_npmUser"`
	//Maintainers []User `json:"maintainers"`
}

type Package struct {
	Name     string             `json:"name"`
	Versions map[string]Version `json:"versions"`
	Time     map[string]string  `json:"time"`

	// Id          string             `json:"_id"`
	// Rev         string             `json:"_rev"`
	// Description string             `json:"description"`
	// DistTags    map[string]string  `json:"dist-tags"`
	// Author      User               `json:"author"`
	// Repository  Repository         `json:"repository"`
	// Readme      string             `json:"readme"`
}

func ParsePackageResponse(data []byte) (*Package, error) {
	var pkg Package
	err := json.Unmarshal(data, &pkg)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling package response: %w", err)
	}
	return &pkg, nil
}

func initializePackageCache() {
	if packageCache != nil {
		return
	}

	items, err := ItemsFromFile(packageCacheFilename)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Warnf("Cannot load cache: %v", err)
		}

		packageCache = cache.New(packageCacheExpiration, packageCacheCleanupInterval)
	} else {
		packageCache = cache.NewFrom(packageCacheExpiration, packageCacheCleanupInterval, items)
	}
}

func persistPackageCache() {
	items := packageCache.Items()
	err := ItemsToFile(packageCacheFilename, items)
	if err != nil {
		log.Warnf("Cannot persist cache: %v", err)
	}
}

func packageJSONFromCache(name string) []byte {
	initializePackageCache()

	data, found := packageCache.Get(name)
	if !found {
		return nil
	}

	return data.([]byte)
}

func packageJSONFromAPI(name string) ([]byte, error) {
	url := fmt.Sprintf("https://registry.npmjs.com/%s", url.QueryEscape(name))
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body of %s: %w", url, err)
	}

	return body, nil
}

func GetPackageJSON(name string) ([]byte, error) {
	data := packageJSONFromCache(name)
	if data != nil {
		log.Debugf("Got package %q from cache.", name)
		return data, nil
	}

	log.Debugf("Loading package %q from API...", name)

	data, err := packageJSONFromAPI(name)
	if err != nil {
		return nil, fmt.Errorf("getting package %q from API: %w", name, err)
	}

	packageCache.Set(name, data, cache.DefaultExpiration)
	persistPackageCache()

	return data, nil
}

func GetPackage(name string) (*Package, error) {
	data, err := GetPackageJSON(name)
	if err != nil {
		return nil, fmt.Errorf("retrieving package %q: %w", name, err)
	}

	pkg, err := ParsePackageResponse(data)
	if err != nil {
		return nil, fmt.Errorf("parsing package %q: %w", name, err)
	}

	if pkg.Name != name {
		return nil, fmt.Errorf("expected package %q, got %q", name, pkg.Name)
	}

	return pkg, nil
}
