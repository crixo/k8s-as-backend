

## How to apply additional/server-side-coding validataion to a CRD before it will be persisted in etcd

- https://banzaicloud.com/blog/k8s-admission-webhooks/
  - [banzai sample](https://github.com/banzaicloud/admission-webhook-example)
- [Dynamic Admission Control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers)
  - [admission webhook server sample](https://github.com/kubernetes/kubernetes/blob/v1.13.0/test/images/webhook/main.go)
  - [deploy webhook server](https://github.com/kubernetes/kubernetes/blob/v1.15.0/test/e2e/apimachinery/webhook.go#L301)


## Dependency Management
The project uses `go mod` but it is requited by `code-generator` for the project
to be in the `GOPATH`. You should export `GO111MODULE=on` to be sure that the
project uses `go mod` even if it is inside the `GOPATH`.


## Generate code
- install or upgrade go 1.13 using homebrew
- clone starting repo with pkg/apis/GROUP_NAME/VERSION/doc.go,register.go,types.go
- add the hack w/ all files folder from https://github.com/kubernetes/sample-controller/tree/master/hack
- adjust the update-codegen.sh with the specific repo reference
- init gomod
```
go mod init github.com/crixo/k8s-as-backend
```
- create autogenerate client
```
go get k8s.io/code-generator
go get k8s.io/client-go@kubernetes-1.16.3
go mod vendor
# optional - # sh ./hack/verify-codegen.sh ``` (optional)
sh ./hack/update-codegen.sh 
```

- write biz logic code

- re-vendoring adding vendor code-related
```
go mod vendor
```

- build the docker image
```
docker build -t crixo/k8s-as-backend .
```

- laod image in kind
```
docker tag crixo/k8s-as-backend crixo/k8s-as-backend:kind
kind load docker-image crixo/k8s-as-backend:kind --name standard
```

### reference
- https://github.com/kubernetes/code-generator

## checks
- get list of enabled admission controller
https://stackoverflow.com/questions/51489955/how-to-obtain-the-enable-admission-controller-list-in-kubernetes
```in console on the control-plane
cat /etc/kubernetes/manifests/kube-apiserver.yaml
```
kubectl -n kube-system describe po kube-apiserver-YOUR-CLUSTER_NAME_REF  
```

## Generate code old
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