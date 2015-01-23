package main

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshsys "github.com/cloudfoundry/bosh-agent/system"

	bhimporter "github.com/cppforlife/bosh-hub/release/importer"
	bhwatcher "github.com/cppforlife/bosh-hub/release/watcher"
	bhstemsimp "github.com/cppforlife/bosh-hub/stemcell/importer"
)

type Config struct {
	Repos ReposOptions

	APIKey string

	// Does not start web server; just does background work
	ActAsWorker bool

	Watcher  bhwatcher.FactoryOptions
	Importer bhimporter.FactoryOptions

	StemcellImporter bhstemsimp.FactoryOptions
}

func NewConfigFromPath(path string, fs boshsys.FileSystem) (Config, error) {
	var config Config

	bytes, err := fs.ReadFile(path)
	if err != nil {
		return config, bosherr.WrapError(err, "Reading config %s", path)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, bosherr.WrapError(err, "Unmarshalling config")
	}

	return config, nil
}
