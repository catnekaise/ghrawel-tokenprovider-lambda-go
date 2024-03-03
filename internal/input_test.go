package internal

import (
	"github.com/google/go-github/v60/github"
	"reflect"
	"testing"
)

func Test_readRepositories(t *testing.T) {
	type args struct {
		endpointType            string
		repositorySelectionMode string
		repo                    *string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "owner endpoint, ALLOW_OWNER repo selection, no repos",
			args: args{
				endpointType:            "DYNAMIC_OWNER",
				repositorySelectionMode: "ALLOW_OWNER",
				repo:                    nil,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "owner endpoint, ALLOW_OWNER repo selection, single repo",
			args: args{
				endpointType:            "DYNAMIC_OWNER",
				repositorySelectionMode: "ALLOW_OWNER",
				repo:                    github.String("repo-1"),
			},
			want:    []string{"repo-1"},
			wantErr: false,
		},
		{
			name: "owner endpoint, ALLOW_OWNER repo selection, multi repos",
			args: args{
				endpointType:            "DYNAMIC_OWNER",
				repositorySelectionMode: "ALLOW_OWNER",
				repo:                    github.String("repo-1,repo-2"),
			},
			want:    []string{"repo-1", "repo-2"},
			wantErr: false,
		},
		{
			name: "default endpoint, AT_LEAST_ONE repo selection, one repo",
			args: args{
				endpointType:            "DEFAULT",
				repositorySelectionMode: "AT_LEAST_ONE",
				repo:                    github.String("repo-1"),
			},
			want:    []string{"repo-1"},
			wantErr: false,
		},
		{
			name: "default endpoint, AT_LEAST_ONE repo selection, no repo",
			args: args{
				endpointType:            "DEFAULT",
				repositorySelectionMode: "AT_LEAST_ONE",
				repo:                    nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "default endpoint, AT_LEAST_ONE repo selection, invalid repo",
			args: args{
				endpointType:            "DEFAULT",
				repositorySelectionMode: "AT_LEAST_ONE",
				repo:                    github.String("repo#"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "default endpoint, AT_LEAST_ONE repo selection, multiple repos",
			args: args{
				endpointType:            "DEFAULT",
				repositorySelectionMode: "AT_LEAST_ONE",
				repo:                    github.String("repo-1,repo-2"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "owner endpoint, ALLOW_OWNER repo selection, no repo",
			args: args{
				endpointType:            "DYNAMIC_OWNER",
				repositorySelectionMode: "ALLOW_OWNER",
				repo:                    nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readRepo(tt.args.endpointType, tt.args.repositorySelectionMode, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("readRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readRepo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
