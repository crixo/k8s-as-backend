

```
kubectl apply -f deploy/crds/k8s-as-backend.example.com_k8sasbackends_crd.yaml

kubectl apply -f deploy/crds/k8s-as-backend.example.com_v1alpha1_k8sasbackend_cr.yaml

operator-sdk run --local --namespace=default
```

[Meta Accessor](https://github.com/kubernetes/apimachinery/blob/master/pkg/api/meta/meta.go)