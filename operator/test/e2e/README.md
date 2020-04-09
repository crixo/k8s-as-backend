# End-to-End Tests

https://github.com/operator-framework/operator-sdk/blob/master/doc/test-framework/writing-e2e-tests.md


```
kubectl create namespace operator-test
SKIP_CLEAN_UP=1 operator-sdk test local ./test/e2e --namespace operator-test --up-local

operator-sdk test local ./test/e2e --namespace operator-test --up-local --go-test-flags "-skipcleanup=true" --global-manifest "./deploy/crds/k8s-as-backend.example.com_k8sasbackends_crd.yaml"
```

http://localhost/operator-test/example-k8sasbackend/todo-app/swagger-ui/index.html

