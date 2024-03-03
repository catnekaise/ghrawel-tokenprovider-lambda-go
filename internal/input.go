package internal

import (
	"errors"
	"fmt"
	"github.com/catnekaise/ghrawel-tokenprovider-lambda-go/pkg/api"
	"log/slog"
	"regexp"
	"strings"
)

var ownerRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]+$`)
var repoRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.-]+$`)

func isRepositorySelectionMode(value string) (bool, error) {

	switch value {
	case api.RepositorySelectionModeAtLeastOne:
		return true, nil
	case api.RepositorySelectionModeAllowOwner:
		return true, nil
	}

	return false, errors.New("invalid repository selection mode")

}

func isEndpointType(value string) (bool, error) {

	switch value {
	case api.EndpointTypeDefault:
		return true, nil
	case api.EndpointTypeStaticOwner:
		return true, nil
	case api.EndpointTypeDynamicOwner:
		return true, nil
	}

	return false, errors.New("invalid endpoint type")
}

func readOwner(owner string) bool {

	return ownerRegex.MatchString(owner)
}

func readRepo(endpointType string, repositorySelectionMode string, repo *string) ([]string, error) {

	var repositories []string

	if repo == nil {
		repositories = []string{}
	} else {
		repositories = strings.Split(*repo, ",")
	}

	switch repositorySelectionMode {
	case api.RepositorySelectionModeAtLeastOne:
		if len(repositories) == 0 {
			return nil, errors.New("at least one repository has to be specified")
		} else if len(repositories) != 1 {
			return nil, errors.New(fmt.Sprintf("Exactly one repository must be specified"))
		}

		break
	case api.RepositorySelectionModeAllowOwner:
		if api.IsOwnerEndpoint(endpointType) && len(repositories) == 0 {
			return nil, nil
		}
		break
	}

	invalid := false

	for _, repository := range repositories {

		valid := repoRegex.MatchString(repository)

		if !valid {
			invalid = true
			slog.Debug(fmt.Sprintf("Invalid name (%q) provided for repository", repository))
		}
	}

	if invalid == true {
		return nil, errors.New(fmt.Sprintf("One or more repositories specified with an invalid name"))
	}

	return repositories, nil

}
