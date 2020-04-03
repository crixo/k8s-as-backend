# End-to-End Tests

https://github.com/operator-framework/operator-sdk/blob/master/doc/test-framework/writing-e2e-tests.md


```
kubectl create namespace operator-test
operator-sdk test local ./test/e2e --namespace operator-test --up-local
```