package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/catnekaise/ghrawel-tokenprovider-lambda-go/pkg/api"
	"log/slog"
	"os"
	"regexp"
)

var prefixRegex = regexp.MustCompile(`^/$|^/[a-zA-Z][a-zA-Z0-9/-]+[a-zA-Z]$`)

func Start() {

	logger := slog.New(logger{minLevel: level(), jsonHandler: handler()})
	slog.SetDefault(logger)

	lambda.Start(run)
}

func run(ctx context.Context, req api.Input) (*tokenResponse, error) {

	secretsPrefix := os.Getenv("SECRETS_PREFIX")
	secretsStorage := os.Getenv("SECRETS_STORAGE")

	if secretsStorage != api.SecretsStorageParameterStore && secretsStorage != api.SecretsStorageSecretsManager {
		slog.ErrorContext(ctx, fmt.Sprintf("Unknown SECRETS_STORAGE %q", secretsStorage))
		return nil, createErrorResponse("Error", 500)
	}

	if !prefixRegex.MatchString(secretsPrefix) {
		slog.ErrorContext(ctx, fmt.Sprintf("Invalid SECRETS_PREFIX %q", secretsPrefix))
		return nil, createErrorResponse("Error", 500)
	}

	ctx = contextWithLoggerFields(ctx, req)
	logInitialRequest(ctx, req)

	return handleInput(ctx, req, secretsStorage, secretsPrefix)
}

func handleInput(ctx context.Context, req api.Input, secretsStorage string, secretsPrefix string) (*tokenResponse, error) {

	if ok, err := isRepositorySelectionMode(req.TokenContext.TargetRule.RepositorySelectionMode); !ok {
		slog.ErrorContext(ctx, err.Error())
		return nil, createErrorResponse("Error", 500)
	}

	if ok, err := isEndpointType(req.TokenContext.Endpoint.Type); !ok {
		slog.ErrorContext(ctx, err.Error())
		return nil, createErrorResponse("Error", 500)
	}

	owner := req.TokenRequest.Owner

	if ok := readOwner(owner); ok == false {
		slog.InfoContext(ctx, fmt.Sprintf("InputError - Value of provided owner (%q) is invalid", owner))
		return nil, createErrorResponse("Value of provided owner is invalid", 400)
	}

	repos, err := readRepo(req.TokenContext.Endpoint.Type, req.TokenContext.TargetRule.RepositorySelectionMode, req.TokenRequest.Repo)

	if err != nil {
		slog.InfoContext(ctx, fmt.Sprintf("InputError - repositories under selection mode %s", req.TokenContext.TargetRule.RepositorySelectionMode))
		return nil, createErrorResponse("Invalid repository selection.", 400)
	}

	return handle(ctx, req, secretsStorage, secretsPrefix, owner, repos)
}

func handle(ctx context.Context, req api.Input, secretsStorage string, secretsPrefix string, owner string, repos []string) (*tokenResponse, error) {

	privateKey, err := getPrivateKey(ctx, secretsStorage, secretsPrefix, req.TokenContext.App.Name)

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("PrivateKeyError - %s", err.Error()))
		return nil, createErrorResponse("Error", 500)
	}

	client, err := createClient(privateKey, req.TokenContext.App.Id)

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("GitHubClientError - %s", err.Error()))
		return nil, createErrorResponse("Error", 500)
	}

	installationId, err := findInstallation(ctx, client, req.TokenContext.App.Id, owner)

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("InstallationNotFound - Could not find installation for %s", owner))
		return nil, createErrorResponse("Error", 500)
	}

	token, err := getToken(ctx, client, installationId, req.TokenContext.Permissions, repos)

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("TokenError - %s", err.Error()))
		return nil, createErrorResponse("Error", 500)
	}

	slog.InfoContext(ctx, "TokenCreated")

	return &tokenResponse{Token: *token.Token}, nil
}

func createErrorResponse(message string, statusCode int) error {

	e := errorResponse{SelectionPattern: fmt.Sprintf("CK_ERR_%v", statusCode), Message: message}

	str, err := json.Marshal(e)

	if err != nil {
		panic(err)
	}

	return errors.New(string(str))
}

type tokenResponse struct {
	Token string `json:"token"`
}

type errorResponse struct {
	SelectionPattern string `json:"selectionPattern"`
	Message          string `json:"message"`
}
