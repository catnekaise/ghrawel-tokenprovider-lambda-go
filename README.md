# ghrawel Token Provider Lambda Go
This application is used by [catnekaise/ghrawel](https://github.com/catnekaise/ghrawel) and that repository also contains all documentation. For specifics about this application, read [here](https://github.com/catnekaise/ghrawel/blob/main/docs/token-provider/application.md) and [here](https://github.com/catnekaise/ghrawel/blob/main/docs/token-provider/infrastructure.md).

## Environment Variables

| Var             | Examples                           |
|-----------------|------------------------------------|
| SECRETS_STORAGE | PARAMETER_STORE or SECRETS_MANAGER |
| SECRETS_PREFIX  | /catnekaise/github-apps            |
| DEBUG_LOGGING   | true                               |