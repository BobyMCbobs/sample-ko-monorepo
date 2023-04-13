# sample-ko-monorepo

> A sample Go app for demonstrating Ko with

## Features

- build each application, where defined in [config/](./config)
- sign container images with [Cosign](https://docs.sigstore.dev/cosign/overview/)
- upload [resolved](https://ko.build/reference/ko_resolve/) manifest to release

## Usage

### Install dependencies

- [kubectl](https://kubectl.sigs.k8s.io/installation/kubectl/)
- [kind](https://kind.sigs.k8s.io)
- [kn](https://knative.dev/docs/client/install-kn/)
- [kn-quickstart](https://knative.dev/docs/getting-started/quickstart-install/)
- [cosign](https://docs.sigstore.dev/cosign/installation/)

### Set up

1. generate a cosign key pair
```shell
cosign generate-key-pair
```

2. set GitHub secrets with `COSIGN_PRIVATE_KEY` with the contents of `cosign.key` and `COSIGN_PASSWORD` with the password (if applicable)

3. commit the public key `cosign.pub` to the root of the repo

4. under Settings -> Code and automation -> Actions -> General, set _Allow GitHub Actions to create and approve pull requests_ to `true`

alternative: import existing keys if using them

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

## TODOs

- [ ] dependency security scanning
- [ ] automatic dependency updates
- [ ] Go version upgrade auto-PR
