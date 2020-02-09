## Generate code
- install or upgrade go 1.13 using homebrew
- clone starting repo with pkg/apis/GROUP_NAME/VERSION/doc.go,register.go,types.go
- add the hack w/ all files folder from https://github.com/kubernetes/sample-controller/tree/master/hack
- adjust the update-codegen.sh with the specific repo reference
- init go mod
```
go mod init github.com/crixo/k8s-as-backend/library
```
- create autogenerate client
```
go get k8s.io/code-generator
go get k8s.io/client-go@kubernetes-1.16.3
go mod vendor
# optional - # sh ./hack/verify-codegen.sh ``` (optional)
sh ./hack/update-codegen.sh 
```