package internal

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/catnekaise/ghrawel-tokenprovider-lambda-go/pkg/api"
	"github.com/google/go-github/v60/github"
	"reflect"
	"regexp"
	"testing"
)

func Test_handleInput(t *testing.T) {
	type args struct {
		req api.Input
	}
	tests := []struct {
		name       string
		args       args
		want       *tokenResponse
		wantErr    bool
		wantErrInt int
	}{
		{
			name: "invalid owner",
			args: args{
				req: createTestInput("catnekaise#", github.String("example-repo"), nil, nil),
			},
			want:       nil,
			wantErr:    true,
			wantErrInt: 400,
		},
		{
			name: "invalid repo",
			args: args{
				req: createTestInput("catnekaise", github.String("example-repo#"), nil, nil),
			},
			want:       nil,
			wantErr:    true,
			wantErrInt: 400,
		},
		{
			name: "multi repos when DEFAULT endpoint type",
			args: args{
				req: createTestInput("catnekaise", github.String("example-repo,example-repo-2"), nil, nil),
			},
			want:       nil,
			wantErr:    true,
			wantErrInt: 400,
		},
		{
			name: "single repo when DYNAMIC_OWNER endpoint type, AT_LEAST_ONE repository selection mode",
			args: args{
				req: createTestInput("catnekaise", github.String("example-repo,example-repo-2"), github.String("DYNAMIC_OWNER"), nil),
			},
			want:       nil,
			wantErr:    true,
			wantErrInt: 400,
		},
		{
			name: "invalid repo when STATIC_OWNER endpoint type, ALLOW_OWNER repository selection mode",
			args: args{
				req: createTestInput("catnekaise", github.String("example-repo#"), github.String("STATIC_OWNER"), github.String("ALLOW_OWNER")),
			},
			want:       nil,
			wantErr:    true,
			wantErrInt: 400,
		},
		{
			name: "invalid endpoint type",
			args: args{
				req: createTestInput("catnekaise", github.String("example-repo,example-repo-2"), github.String("STANDARD"), nil),
			},
			want:       nil,
			wantErr:    true,
			wantErrInt: 500,
		},

		{
			name: "invalid repository selection mode",
			args: args{
				req: createTestInput("catnekaise", github.String("example-repo,example-repo-2"), nil, github.String("ALL")),
			},
			want:       nil,
			wantErr:    true,
			wantErrInt: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleInput(context.TODO(), tt.args.req, "PARAMETER_STORE", "/")
			if (err != nil) != tt.wantErr {
				t.Errorf("handleInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleInput() got = %v, want %v", got, tt.want)
			}

			errCode := fmt.Sprintf("CK_ERR_%v", tt.wantErrInt)
			if !regexp.MustCompile(errCode).MatchString(err.Error()) {
				t.Errorf("run() got = %v, want %v", err.Error(), errCode)
			}
		})
	}
}

func Test_run(t *testing.T) {

	t.Run("Secrets Storage", func(t *testing.T) {

		t.Setenv("SECRETS_PREFIX", "/")
		t.Setenv("SECRETS_STORAGE", "S3")

		_, err := run(context.TODO(), createTestInput("catnekaise", github.String("repo-1"), nil, nil))

		if err == nil {
			t.Error("run() did not return error as expected")
		}

		if !regexp.MustCompile("CK_ERR_500").MatchString(err.Error()) {
			t.Errorf("run() got = %v, want %v", err.Error(), "CK_ERR_500")
		}
	})

	t.Run("Secrets Prefix", func(t *testing.T) {
		t.Setenv("SECRETS_PREFIX", "#/")
		t.Setenv("SECRETS_STORAGE", "PARAMETER_STORE")

		_, err := run(context.TODO(), createTestInput("catnekaise", github.String("repo-1"), nil, nil))

		if err == nil {
			t.Error("run() did not return error as expected")
		}

		fmt.Println(err.Error())

		if !regexp.MustCompile("CK_ERR_500").MatchString(err.Error()) {
			t.Errorf("run() got = %v, want %v", err.Error(), "CK_ERR_500")
		}
	})

	t.Run("Bad Repo", func(t *testing.T) {
		t.Setenv("SECRETS_PREFIX", "/")
		t.Setenv("SECRETS_STORAGE", "SECRETS_MANAGER")

		_, err := run(context.TODO(), createTestInput("catnekaise", github.String("repo-##1"), nil, nil))

		if err == nil {
			t.Error("run() did not return error as expected")
		}

		if !regexp.MustCompile("CK_ERR_400").MatchString(err.Error()) {
			t.Errorf("run() got = %v, want %v", err.Error(), "CK_ERR_400")
		}
	})
}

func Test_prefix(t *testing.T) {
	tests := []struct {
		value string
		match bool
	}{
		{
			value: "/",
			match: true,
		},
		{
			value: "/catnekaise",
			match: true,
		},
		{
			value: "/catnekaise/github-apps",
			match: true,
		},
		{
			value: "/catnekaise/a/b/c/d/f/g/h/j/k/l",
			match: true,
		},
		{
			value: "catnekaise/github-apps",
			match: false,
		},
		{
			value: "/catnekaise/github_apps",
			match: false,
		},
		{
			value: "/catnekaise/github:apps",
			match: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if prefixRegex.MatchString(tt.value) != tt.match {
				t.Errorf("prefixRegex() got = %v, want %v", !tt.match, tt.match)
			}
		})
	}
}

func createTestInput(owner string, repo *string, endpointType *string, repositorySelectionMode *string) api.Input {

	t := "DEFAULT"
	r := "AT_LEAST_ONE"

	if endpointType != nil {
		t = *endpointType
	}

	if repositorySelectionMode != nil {
		r = *repositorySelectionMode
	}

	return api.Input{
		APIGatewayProxyRequest: events.APIGatewayProxyRequest{},
		TokenRequest: api.TokenRequest{
			Owner: owner,
			Repo:  repo,
		},
		TokenContext: api.TokenContext{
			ProviderName: "test",
			Permissions: api.Permissions{
				Contents: github.String("read"),
			},
			App: api.App{
				Name: "default",
				Id:   1234,
			},
			Endpoint: api.Endpoint{
				Type: t,
			},
			TargetRule: api.TargetRule{
				RepositorySelectionMode: r,
			},
		},
	}
}
