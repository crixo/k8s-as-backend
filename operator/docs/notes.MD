# Notes

## Cert tools
```
openssl req -in ~/Coding/golang/k8s-as-backend/operator/server-test.csr -noout -text
```

## Extra settings via env variable
```
export PEM_FOLDER="/Users/cristiano/Coding/golang/k8s-as-backend/operator/certs"
operator-sdk run --local --namespace=default
```
TODO: is it possible to provide flag through the SDK launcher?


## Update scaffolding
From [operator-SDK user guides](https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md)

After modifying the *_types.go file always run the following command to update the generated code for that resource type:
```
operator-sdk generate k8s
```

Now that GO types have fields and possibly annotations, the CRD corresponding to the API's group and kind must be updated. To do so, run the following command:
```
operator-sdk generate crds
```

## Testing
run all tests in all packages
```
go test ./...
```

run a specific test
```
go test -v ./pkg/controller/k8sasbackend/webhookserver/
```

## References
- [Meta Accessor](https://github.com/kubernetes/apimachinery/blob/master/pkg/api/meta/meta.go)