package api

import "github.com/aws/aws-lambda-go/events"

const (
	RepositorySelectionModeAllowOwner = "ALLOW_OWNER"
	RepositorySelectionModeAtLeastOne = "AT_LEAST_ONE"
	EndpointTypeDefault               = "DEFAULT"
	EndpointTypeStaticOwner           = "STATIC_OWNER"
	EndpointTypeDynamicOwner          = "DYNAMIC_OWNER"
	SecretsStorageParameterStore      = "PARAMETER_STORE"
	SecretsStorageSecretsManager      = "SECRETS_MANAGER"
)

type Permissions struct {
	Actions                                 *string `json:"actions,omitempty"`
	Administration                          *string `json:"administration,omitempty"`
	Checks                                  *string `json:"checks,omitempty"`
	Codespaces                              *string `json:"codespaces,omitempty"`
	Contents                                *string `json:"contents,omitempty"`
	DependabotSecrets                       *string `json:"dependabot_secrets,omitempty"`
	Deployments                             *string `json:"deployments,omitempty"`
	Environments                            *string `json:"environments,omitempty"`
	Issues                                  *string `json:"issues,omitempty"`
	Metadata                                *string `json:"metadata,omitempty"`
	Packages                                *string `json:"packages,omitempty"`
	Pages                                   *string `json:"pages,omitempty"`
	PullRequests                            *string `json:"pull_requests,omitempty"`
	RepositoryCustomProperties              *string `json:"repository_custom_properties,omitempty"`
	RepositoryHooks                         *string `json:"repository_hooks,omitempty"`
	RepositoryProjects                      *string `json:"repository_projects,omitempty"`
	SecretScanningAlerts                    *string `json:"secret_scanning_alerts,omitempty"`
	Secrets                                 *string `json:"secrets,omitempty"`
	SecurityEvents                          *string `json:"security_events,omitempty"`
	SingleFile                              *string `json:"single_file,omitempty"`
	Statuses                                *string `json:"statuses,omitempty"`
	VulnerabilityAlerts                     *string `json:"vulnerability_alerts,omitempty"`
	Workflows                               *string `json:"workflows,omitempty"`
	Members                                 *string `json:"members,omitempty"`
	OrganizationAdministration              *string `json:"organization_administration,omitempty"`
	OrganizationCustomRoles                 *string `json:"organization_custom_roles,omitempty"`
	OrganizationCustomOrgRoles              *string `json:"organization_custom_org_roles,omitempty"`
	OrganizationCustomProperties            *string `json:"organization_custom_properties,omitempty"`
	OrganizationCopilotSeatManagement       *string `json:"organization_copilot_seat_management,omitempty"`
	OrganizationAnnouncementBanners         *string `json:"organization_announcement_banners,omitempty"`
	OrganizationEvents                      *string `json:"organization_events,omitempty"`
	OrganizationHooks                       *string `json:"organization_hooks,omitempty"`
	OrganizationPersonalAccessTokens        *string `json:"organization_personal_access_tokens,omitempty"`
	OrganizationPersonalAccessTokenRequests *string `json:"organization_personal_access_token_requests,omitempty"`
	OrganizationPlan                        *string `json:"organization_plan,omitempty"`
	OrganizationProjects                    *string `json:"organization_projects,omitempty"`
	OrganizationPackages                    *string `json:"organization_packages,omitempty"`
	OrganizationSecrets                     *string `json:"organization_secrets,omitempty"`
	OrganizationSelfHostedRunners           *string `json:"organization_self_hosted_runners,omitempty"`
	OrganizationUserBlocking                *string `json:"organization_user_blocking,omitempty"`
	TeamDiscussions                         *string `json:"team_discussions,omitempty"`
	EmailAddresses                          *string `json:"email_addresses,omitempty"`
	Followers                               *string `json:"followers,omitempty"`
	GitSshKeys                              *string `json:"git_ssh_keys,omitempty"`
	GpgKeys                                 *string `json:"gpg_keys,omitempty"`
	InteractionLimits                       *string `json:"interaction_limits,omitempty"`
	Profile                                 *string `json:"profile,omitempty"`
	Starring                                *string `json:"starring,omitempty"`
}

type TokenRequest struct {
	Owner string  `json:"owner"`
	Repo  *string `json:"repo"`
}

type TokenContext struct {
	ProviderName string      `json:"providerName"`
	Permissions  Permissions `json:"permissions"`
	App          App         `json:"app"`
	Endpoint     Endpoint    `json:"endpoint"`
	TargetRule   TargetRule  `json:"targetRule"`
}

type App struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Endpoint struct {
	Type string `json:"type"`
}

type TargetRule struct {
	RepositorySelectionMode string `json:"repositorySelectionMode"`
}

type Input struct {
	events.APIGatewayProxyRequest
	TokenRequest TokenRequest `json:"tokenRequest"`
	TokenContext TokenContext `json:"tokenContext"`
}

func IsOwnerEndpoint(value string) bool {
	return value == EndpointTypeDynamicOwner || value == EndpointTypeStaticOwner
}
