openssl req -in ~/Coding/golang/k8s-as-backend/operator/server-test.csr -noout -text

export PEM_FOLDER="/Users/cristiano/Coding/golang/k8s-as-backend/operator/certs"
operator-sdk run --local --namespace=default


[Meta Accessor](https://github.com/kubernetes/apimachinery/blob/master/pkg/api/meta/meta.go)