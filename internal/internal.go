package internal

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/rs/zerolog/log"
)

const (
	ReleaseInfoSource string = "https://releases.hashicorp.com/index.json"
)

type Releases struct {
	Terraform struct {
		Versions map[string]struct{}
	}
}

func GetModule(dir string) (*tfconfig.Module, error) {
	path, err := filepath.Abs(dir)
	if err != nil {
		log.Error().Err(err).Msg("Failed to resolve module path")
		return nil, err
	}
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Error().Err(err).Msg("Module path does not exist")
		return nil, err
	}
	if !info.IsDir() {
		err := errors.New("path is not a directory")
		log.Error().Err(err).Msg("Module path does not exist")
		return nil, err
	}
	if !tfconfig.IsModuleDir(path) {
		err := errors.New("the directory has no terraform configuration files")
		log.Error().Err(err).Str("dir", path).Msg("Module not found")
		return nil, err
	}
	module, _ := tfconfig.LoadModule(dir)
	return module, nil
}

func GetConstraints(module *tfconfig.Module) (version.Constraints, error) {
	if len(module.RequiredCore) == 0 {
		err := errors.New(`terraform "required_version" attribute is required`)
		log.Error().Err(err).Msg("Required version not found")
		return nil, err
	}
	v := strings.Join(module.RequiredCore, ", ")
	constraints, err := version.NewConstraint(v)
	if err != nil {
		log.Error().Err(err).Str("constraints", v).Msg("Failed to parse required version")
		return nil, err
	}
	return constraints, nil
}

func GetVersions() ([]*version.Version, error) {
	response, err := http.Get(ReleaseInfoSource)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
	}
	response.Body.Close()
	if response.StatusCode > 299 {
		log.Error().
			Str("source", ReleaseInfoSource).Int("code", response.StatusCode).Bytes("body", body).
			Msg("Failed to request the release information source")
	}
	var releases Releases
	err = json.Unmarshal(body, &releases)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal response")
		return nil, err
	}
	versions := make([]*version.Version, 0, len(releases.Terraform.Versions))
	for k := range releases.Terraform.Versions {
		v, err := version.NewVersion(k)
		if err != nil {
			log.Error().Err(err).Str("version", k).Msg("Failed to parse the given terraform version")
			continue
		}
		versions = append(versions, v)
	}
	sort.Sort(sort.Reverse(version.Collection(versions)))
	log.Debug().Any("versions", versions).Msg("Available terraform")
	return versions, nil
}

func GetLatestRequired(constraints version.Constraints, versions []*version.Version) (*version.Version, error) {
	for _, v := range versions {
		if constraints.Check(v) {
			log.Info().Str("version", v.String()).Msg("Required version found")
			return v, nil
		}
	}
	err := errors.New("unsupported terraform core version")
	log.Error().Err(err).Str("constraints", constraints.String()).Msg("Failed to find required version")
	return nil, err
}
