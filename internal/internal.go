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
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/rs/zerolog/log"
)

const (
	ReleaseListLimit    int    = 20
	ReleaseListEndpoint string = "https://api.releases.hashicorp.com/v1/releases/terraform"
)

type Release struct {
	Version          string
	TimestampCreated *time.Time `json:"timestamp_created,omitempty"`
}

type Error struct {
	Code    int
	Message string
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

func GetReleases(endpoint string) ([]*Release, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create request")
		return nil, nil
	}
	req.Header.Set("accept", "application/json")
	query := req.URL.Query()

	var releases []*Release
	var after *time.Time

	for {
		if after != nil {
			query.Set("after", after.Format(time.RFC3339Nano))
			req.URL.RawQuery = query.Encode()
		}
		log.Debug().Str("url", req.URL.String()).Msg("Got release list")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to make request")
			return nil, nil
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read response body")
		}
		defer resp.Body.Close()

		if contentType := resp.Header.Get("content-type"); contentType != "application/json" {
			log.Error().Str("content-type", contentType).Msg("Wrong release list response content type")
			return nil, nil
		}
		if resp.StatusCode > 200 {
			var error Error
			err = json.Unmarshal(body, &error)
			event := log.Error()
			if err != nil {
				event.Int("code", resp.StatusCode).Bytes("body", body)
			} else {
				event.Int("code", error.Code).Str("message", error.Message)
			}
			event.Msg("Failed to get release list response")
			return nil, nil
		}

		var releaseList []*Release
		err = json.Unmarshal(body, &releaseList)
		if err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal response")
			return nil, err
		}
		if len(releaseList) == 0 {
			break
		}

		releases = append(releases, releaseList...)
		if after = releaseList[len(releaseList)-1].TimestampCreated; after == nil {
			break
		}
	}

	return releases, nil
}

func GetVersions() ([]*version.Version, error) {
	releases, err := GetReleases(ReleaseListEndpoint)
	if err != nil {
		return nil, err
	}
	versions := make([]*version.Version, 0, len(releases))
	for _, r := range releases {
		v, err := version.NewVersion(r.Version)
		if err != nil {
			log.Error().Err(err).Str("version", r.Version).Msg("Failed to parse the given terraform version")
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
			log.Debug().Str("version", v.String()).Msg("Required version found")
			return v, nil
		}
	}
	err := errors.New("unsupported terraform core version")
	log.Error().Err(err).Str("constraints", constraints.String()).Msg("Failed to find required version")
	return nil, err
}
