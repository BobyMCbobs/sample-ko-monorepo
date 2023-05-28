<!-- generate TOC using `go run sigs.k8s.io/mdtoc@latest --inplace README.md` -->
<!-- toc -->
- [sample-ko-monorepo](#sample-ko-monorepo)
  - [Features](#features)
  - [Automations](#automations)
  - [Usage](#usage)
    - [Install dependencies](#install-dependencies)
    - [Set up](#set-up)
  - [Install products](#install-products)
  - [Locally run binaries](#locally-run-binaries)
  - [Locally build](#locally-build)
  - [Signatures and attestations](#signatures-and-attestations)
  - [Verifying](#verifying)
  - [Troubleshooting](#troubleshooting)
    - [images fail to push](#images-fail-to-push)
  - [TODOs](#todos)
  - [Related](#related)
<!-- /toc -->

# sample-ko-monorepo

> A sample Go app for demonstrating Ko with

## Features

- build each application, where Go package main entrypoints are
- sign container images with [Cosign](https://docs.sigstore.dev/cosign/overview/)

## Automations

| Name              | Description                                                                     | Link                                             |
|-------------------|---------------------------------------------------------------------------------|--------------------------------------------------|
| Build             | Builds and signs Go based container images (ko, cosign)                         | [link](.github/workflows/build.yml)              |
| Go test           | Runs `go test` against the repo                                                 | [link](.github/workflows/go-test.yml)            |
| Lint              | Lints for code quality (golangci)                                               | [link](.github/workflows/golangci-lint.yml)      |
| Image promotion   | Tags images using image digests                                                 | [link](.github/workflows/image-promotion.yml)    |
| Conform           | Ensures that commits in PRs are standardised                                    | [link](.github/workflows/policy-conformance.yml) |
| Update Go version | Ensures that the Go version which the applications use, is on the latest stable | [link](.github/workflows/update-go-version.yaml) |

all of the actions are implementing reusable workflows.

## Usage

### Install dependencies

- [kubectl](https://kubectl.sigs.k8s.io/installation/kubectl/)
- [kind](https://kind.sigs.k8s.io)
- [kn](https://knative.dev/docs/client/install-kn/)
- [kn-quickstart](https://knative.dev/docs/getting-started/quickstart-install/)
- [cosign](https://docs.sigstore.dev/cosign/installation/)

### Set up

1. under Settings -> Code and automation -> Actions -> General, set _Allow GitHub Actions to create and approve pull requests_ to `true`

2. add a branch protection rule under Settings -> Code and automation -> Add rule
entering

```yaml
Branch name pattern: main
Require a pull request before merging: true
Require status checks to pass before merging: true
  Require branches to be up to date before merging: true
  Status checks:
    - golangci / lint
    - conform / conform
Require signed commits
```

## Install products

launch a local kind cluster, pre-installed with Knative
```shell
kn quickstart kind
```

apply the pre-built release
```shell
kubectl apply -f https://github.com/BobyMCbobs/sample-ko-monorepo/raw/main/deploy/release.yaml
```

## Locally run binaries

```shell
go run cmd/webthingy/main.go
```

```shell
go run cmd/mission-critical-service/main.go
```

## Locally build

```shell
export KO_DOCKER_REPO=ghcr.io/bobymcbobs/sample-ko-monorepo
ko resolve --bare -f config/
```

## Signatures and attestations

```shell
cosign tree IMAGE_REF
```


## Verifying

container images are able to be verified with the following command

```shell
cosign verify ghcr.io/bobymcbobs/sample-ko-monorepo/mission-critical-service@sha256:405b54637c79a0b0934d0d7f01464f358fe1fd118fefb1d9b77c8a351e9471b6 --certificate-identity https://github.com/BobyMCbobs/sample-ko-monorepo/.github/workflows/reusable-build.yml@refs/heads/main --certificate-oidc-issuer https://token.actions.githubusercontent.com
```

SBOMs attestations are able to be verified with the following command

```shell
cosign verify-attestation ghcr.io/bobymcbobs/sample-ko-monorepo/mission-critical-service@sha256:405b54637c79a0b0934d0d7f01464f358fe1fd118fefb1d9b77c8a351e9471b6 --certificate-identity https://github.com/BobyMCbobs/sample-ko-monorepo/.github/workflows/reusable-build.yml@refs/heads/main --certificate-oidc-issuer https://token.actions.githubusercontent.com  | jq -r .payload | base64 -d | jq -r .predicate.Data | bom document outline -
```

## Troubleshooting

### images fail to push

adjust the actions package access settings in
1. go to github.com/{{org/user}}
2. go to the packages tab
3. click on the package failing
4. ensure that the Actions repository access is set up to point to the source repo
5. set _manage Actions access_ role field to `write`

## TODOs

- [ ] dependency security scanning
- [ ] automatic dependency updates
- [x] Go version upgrade auto-PR
- [x] add build dependency cache

## Related

- [sample-docker-monorepo](https://github.com/BobyMCbobs/sample-docker-monorepo)
