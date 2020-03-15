

```
kubectl apply -f deploy/crds/k8s-as-backend.example.com_k8sasbackends_crd.yaml

kubectl apply -f deploy/crds/k8s-as-backend.example.com_v1alpha1_k8sasbackend_cr.yaml

operator-sdk run --local --namespace=default
```

http://localhost/default/example-k8sasbackend/todo-app/swagger-ui

	//openssl req -in ~/Coding/golang/k8s-as-backend/operator/server-test.csr -noout -text

[Meta Accessor](https://github.com/kubernetes/apimachinery/blob/master/pkg/api/meta/meta.go)