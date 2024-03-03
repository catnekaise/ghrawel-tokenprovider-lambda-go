package internal

import (
	"context"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/catnekaise/ghrawel-tokenprovider-lambda-go/pkg/api"
	"log/slog"
	"os"
)

type logger struct {
	minLevel    slog.Level
	jsonHandler *slog.JSONHandler
}

type extraFields struct {
	FunctionArn       string
	GithubAppId       int64
	RequestId         string
	TokenProviderName string
	TokenRequestOwner string
	TokenRequestRepo  *string
	User              string
	UserArn           string
}

type ctxKey struct{}

func level() slog.Level {

	if os.Getenv("DEBUG_LOGGING") == "true" {
		return slog.LevelDebug
	}

	return slog.LevelInfo
}

func handler() *slog.JSONHandler {

	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level(),
	})
}

func (m logger) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= m.minLevel
}

func (m logger) Handle(ctx context.Context, record slog.Record) error {

	xfValue := ctx.Value(ctxKey{})

	if xfValue == nil {
		return m.jsonHandler.Handle(ctx, record)
	}

	fields := xfValue.(extraFields)

	attrs := []slog.Attr{
		{
			Key:   "awsRequestId",
			Value: slog.StringValue(fields.RequestId),
		},
		{
			Key:   "functionArn",
			Value: slog.StringValue(fields.FunctionArn),
		},
		{
			Key:   "tokenProviderName",
			Value: slog.StringValue(fields.TokenProviderName),
		},
		{
			Key:   "userArn",
			Value: slog.StringValue(fields.UserArn),
		},
		{
			Key:   "user",
			Value: slog.StringValue(fields.User),
		},
		{
			Key:   "tokenRequestOwner",
			Value: slog.StringValue(fields.TokenRequestOwner),
		},
		{
			Key:   "githubAppId",
			Value: slog.Int64Value(fields.GithubAppId),
		},
	}

	if fields.TokenRequestRepo != nil {
		attrs = append(attrs, slog.Attr{Key: "tokenRequestRepo", Value: slog.StringValue(*fields.TokenRequestRepo)})

	}

	return m.jsonHandler.WithAttrs(attrs).Handle(ctx, record)
}

func (m logger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return m.jsonHandler.WithAttrs(attrs)
}

func (m logger) WithGroup(name string) slog.Handler {
	return m.jsonHandler.WithGroup(name)
}

func contextWithLoggerFields(ctx context.Context, req api.Input) context.Context {

	fields := extraFields{
		RequestId:         req.RequestContext.RequestID,
		TokenProviderName: req.TokenContext.ProviderName,
		UserArn:           req.RequestContext.Identity.UserArn,
		User:              req.RequestContext.Identity.User,
		TokenRequestOwner: req.TokenRequest.Owner,
		TokenRequestRepo:  req.TokenRequest.Repo,
		GithubAppId:       req.TokenContext.App.Id,
	}

	if lmbCtx, ok := lambdacontext.FromContext(ctx); ok {
		fields.FunctionArn = lmbCtx.InvokedFunctionArn
	}

	return context.WithValue(ctx, ctxKey{}, fields)
}

func logInitialRequest(ctx context.Context, req api.Input) {

	attrs := []slog.Attr{
		{
			Key:   "path",
			Value: slog.StringValue(req.RequestContext.Path),
		},
		{
			Key:   "endpointType",
			Value: slog.StringValue(req.TokenContext.Endpoint.Type),
		},
		{
			Key:   "repositorySelectionMode",
			Value: slog.StringValue(req.TokenContext.TargetRule.RepositorySelectionMode),
		},
		{
			Key:   "githubAppName",
			Value: slog.StringValue(req.TokenContext.App.Name),
		},
		{
			Key:   "userAgent",
			Value: slog.StringValue(req.RequestContext.Identity.UserAgent),
		},
		{
			Key:   "permissions",
			Value: slog.AnyValue(req.TokenContext.Permissions),
		},
	}

	if req.RequestContext.Identity.CognitoIdentityPoolID != "" {
		attrs = append(attrs, slog.Attr{Key: "cognitoIdentityPoolID", Value: slog.StringValue(req.RequestContext.Identity.CognitoIdentityPoolID)})
		attrs = append(attrs, slog.Attr{Key: "cognitoIdentityID", Value: slog.StringValue(req.RequestContext.Identity.CognitoIdentityID)})
		attrs = append(attrs, slog.Attr{Key: "cognitoAuthenticationProvider", Value: slog.StringValue(req.RequestContext.Identity.CognitoAuthenticationProvider)})
		attrs = append(attrs, slog.Attr{Key: "cognitoAuthenticationType", Value: slog.StringValue(req.RequestContext.Identity.CognitoAuthenticationType)})
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "Init", attrs...)
}
