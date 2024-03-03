package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/catnekaise/ghrawel-tokenprovider-lambda-go/pkg/api"
	"github.com/google/go-github/v60/github"
	"net/http"
)

var knownInstallations = map[int64]map[string]int64{}

type installationTokenOptions struct {
	Repositories []string         `json:"repositories,omitempty"`
	Permissions  *api.Permissions `json:"permissions"`
}

func createClient(privateKey *string, appId int64) (*github.Client, error) {

	itr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, appId, []byte(*privateKey))

	if err != nil {
		return nil, err
	}

	return github.NewClient(&http.Client{Transport: itr}), nil
}

func findInstallation(ctx context.Context, client *github.Client, appId int64, owner string) (*int64, error) {

	if installations, ok := knownInstallations[appId]; ok {
		if installationId, ok2 := installations[owner]; ok2 {
			return &installationId, nil
		}
	} else {
		knownInstallations[appId] = map[string]int64{}
	}

	var installationId *int64

	appInstallations, _, err := client.Apps.ListInstallations(ctx, &github.ListOptions{})

	if err != nil {
		return nil, err
	}

	for _, installation := range appInstallations {

		if *installation.Account.Login == owner {
			installationId = installation.ID
			break
		}
	}

	if installationId == nil {
		return nil, errors.New("could not find installation")
	}

	knownInstallations[appId][owner] = *installationId

	return installationId, nil
}

func getToken(ctx context.Context, client *github.Client, installationId *int64, permissions api.Permissions, repo []string) (*github.InstallationToken, error) {

	u := fmt.Sprintf("app/installations/%v/access_tokens", *installationId)

	body := installationTokenOptions{
		Repositories: repo,
		Permissions:  &permissions,
	}

	request, err := client.NewRequest("POST", u, body)

	if err != nil {
		return nil, err
	}

	token := new(github.InstallationToken)
	_, err = client.Do(ctx, request, token)

	if err != nil {
		return nil, err
	}

	return token, nil
}
