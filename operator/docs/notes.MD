# Notes

## Cert tools
```
openssl req -in ~/Coding/golang/k8s-as-backend/operator/server-test.csr -noout -text
```

open cert
```
openssl x509 -in ~/Desktop/k8s/kab01_operator-in-cluster_cert.pem -noout -text
```


### Verifying that a Certificate is issued by a CA
How to use OpenSSL on the command line to verify that a certificate was issued by a specific CA, given that CA's certificate
```
openssl verify -verbose -CAfile cacert.pem  server.crt
```
answer should be
```
server.crt: OK
```
If you get any other message, the certificate was not issued by that CA.

```
CA_BUNDLE=$(kubectl config view --raw -o json | jq -r '.clusters[] | select(.name == "'$(kubectl config current-context)'") | .cluster."certificate-authority-data"')
echo $CA_BUNDLE | base64 -D >> cacert.pem
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

## go mod pinning and retrieving a pinned version

from replace entry
```
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
```
[Pseudo_versions](https://golang.org/cmd/go/#hdr-Pseudo_versions)
```k8s.io/client-go v0.0.0-{commit time in UTC}-{prefix of the commit hash}```

get the full hash

- download the whole repo
```
git clone https://github.com/kubernetes/client-go.git
```

- get the specific commit by prefix
```
git show {prefix of the commit hash}
```

- download the whole repo at that hash
```
https://github.com/{username}/{projectname}/archive/{sha}.zip
https://github.com/kubernetes/client-go/archive/bec269661e48cb1e5fbb0d037f356ffe9e9978a0.zip
```

- Open point

  - How to find that specific hash, so all the others set by operator-SDK, is related to the *kubernetes-1.16.2* release? IOW how to set manually all the replace entries pointing to specific commit related to to given kubernetes release/tag? Does it exist a script to get ```k8s.io/XXXX v0.0.0-{commit time in UTC}-{prefix of the commit hash}``` for each k8s sub packages?

  - dig more on ```k8s.io/client-go v12.0.0+incompatible``` vs ```k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48``` why *require* uses a specific version *v12.0.0* and not just *v0.0.0*?
    [+incompatible](https://github.com/golang/go/wiki/Modules#can-a-module-consume-a-v2-package-that-has-not-opted-into-modules-what-does-incompatible-mean)

  - operator-SDK pinning vs manual pinning done for *admission-webhook* copying from the go.mod *replace* section from specific k8s release.

  - Why dependencies in go.mod are changing time by time? Is it happening only when you install new dependencies?