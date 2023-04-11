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

## Install products

launch a local kind cluster, pre-installed with Knative
```shell
kn quickstart kind
```

apply the pre-built release
```shell
kubectl apply -f https://github.com/BobyMCbobs/sample-ko-monorepo/releases/download/v0.2.0/release.yaml
```

## TODOs

- [ ] dependency security scanning
- [ ] automatic dependency updates
- [ ] Go version upgrade auto-PR
