# unity-management-console

This console is the core management portal for the Unity environment. It configures and customizes each environment with
various specifics for running Unity software. It will also install the software requested by the user.

## Environment Variables

Local testing and development

```AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
AWS_SESSION_TOKEN
GITHUB_TOKEN```

## Running
Grab a management console zip file from the [releases page](https://github.com/unity-sds/unity-management-console/releases)
Unzip to your target destination.

```shell

unzip managementconsole.zip
cd management-console
./main
```

## Development

```shell
npm run dev
```
Run the development environment. This is a Svelte only environment it provides a fake backend and http responses.

```shell
npm run build-all
```
Build the frontend and the backend. This also runs various linters and other code quality checks.

```shell
npm run serve
```
Launch the frontend and the backend locally for testing.
