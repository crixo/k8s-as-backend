

## How to apply additional/server-side-coding validataion to a CRD before it will be persisted in etcd

- https://banzaicloud.com/blog/k8s-admission-webhooks/
- https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#admission-webhooks


## Dependency Management
The project uses `go mod` but it is requited by `code-generator` for the project
to be in the `GOPATH`. You should export `GO111MODULE=on` to be sure that the
project uses `go mod` even if it is inside the `GOPATH`.

## Generate code
Checkout the project in your GOPATH because `code-generator` still uses
`GOPATH`. And this is the command I used to generate the code inside
`pkg/client`.

First get `code-generator` from a different shell w/o GO111MODULE exported
```
go get -u k8s.io/code-generator
```
then move to ~/go/src/k8s.io/code-generator and run to get all package dependencies
```
go get ./...
```

then run from the same shell the following command

```
~/go/src/k8s.io/code-generator/generate-groups.sh all \
    github.com/crixo/k8s-as-backend/pkg/client \
    github.com/crixo/k8s-as-backend/pkg/apis k8sasbackend:v1
```