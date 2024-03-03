package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/catnekaise/ghrawel-tokenprovider-lambda-go/pkg/api"
)

var paramStoreClient *ssm.Client
var secretsManagerClient *secretsmanager.Client

func getPrivateKey(ctx context.Context, storage string, prefix string, name string) (*string, error) {

	if storage == api.SecretsStorageParameterStore {
		return getPrivateKeyParameterStore(ctx, prefix, name)
	} else if storage == api.SecretsStorageSecretsManager {
		return getPrivateKeySecretsManager(ctx, prefix, name)
	}

	return nil, errors.New(fmt.Sprintf("Unknown storage type %q", storage))
}

func getPrivateKeySecretsManager(ctx context.Context, prefix string, name string) (*string, error) {

	if secretsManagerClient == nil {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, err
		}
		secretsManagerClient = secretsmanager.NewFromConfig(cfg)
	}

	secret, err := secretsManagerClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(fmt.Sprintf("%s/%s", prefix, name)),
	})

	if err != nil {
		return nil, err
	}

	return aws.String(*secret.SecretString), nil
}

func getPrivateKeyParameterStore(ctx context.Context, prefix string, name string) (*string, error) {

	if paramStoreClient == nil {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, err
		}
		paramStoreClient = ssm.NewFromConfig(cfg)
	}

	parameter, err := paramStoreClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("%s/%s", prefix, name)),
		WithDecryption: aws.Bool(true),
	})

	if err != nil {
		return nil, err
	}

	value := *parameter.Parameter.Value

	return &value, nil
}
