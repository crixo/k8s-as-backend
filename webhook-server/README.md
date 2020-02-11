https://github.com/banzaicloud/admission-webhook-example

https://github.com/kubernetes/kubernetes/tree/release-1.16/test/images/agnhost/webhook

https://github.com/microsoft/vscode-go/issues/2679


- vendoring adding vendor code-related
```
go mod vendor
```

- build the docker image
```
docker build -t crixo/k8s-as-backend-webhook-server:v.0.0.0 .
```

- laod image in kind
```
kind load docker-image crixo/k8s-as-backend-webhook-server:v.0.0.0 --name standard
```


## Error
```
~/Coding/golang/k8s-as-backend/informer$ k apply -f artifacts/todo.yaml 
Error from server (InternalError): error when applying patch:
{"metadata":{"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"k8sasbackend.com/v1\",\"kind\":\"Todo\",\"metadata\":{\"annotations\":{},\"name\":\"buy-book\",\"namespace\":\"default\"},\"spec\":{\"message\":\"Remember to buy a book about cloud on Amazon upd9.\",\"when\":\"2019-05-13T21:02:21Z\"}}\n"}},"spec":{"message":"Remember to buy a book about cloud on Amazon upd9."}}
to:
Resource: "k8sasbackend.com/v1, Resource=todos", GroupVersionKind: "k8sasbackend.com/v1, Kind=Todo"
Name: "buy-book", Namespace: "default"
Object: &{map["apiVersion":"k8sasbackend.com/v1" "kind":"Todo" "metadata":map["annotations":map["kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"k8sasbackend.com/v1\",\"kind\":\"Todo\",\"metadata\":{\"annotations\":{},\"name\":\"buy-book\",\"namespace\":\"default\"},\"spec\":{\"message\":\"Remember to buy a book about cloud on Amazon upd8.\",\"when\":\"2019-05-13T21:02:21Z\"}}\n"] "creationTimestamp":"2020-02-10T23:36:40Z" "generation":'\x01' "name":"buy-book" "namespace":"default" "resourceVersion":"640" "selfLink":"/apis/k8sasbackend.com/v1/namespaces/default/todos/buy-book" "uid":"f9076dd8-76fc-4a84-ad95-4783ed3d676d"] "spec":map["message":"Remember to buy a book about cloud on Amazon upd8." "when":"2019-05-13T21:02:21Z"]]}
for: "artifacts/todo.yaml": Internal error occurred: failed calling webhook "pod-policy.example.com": Post https://admission-webhook-example-svc.default.svc:443/crd?timeout=5s: dial tcp: lookup admission-webhook-example-svc.default.svc on 192.168.65.1:53: no such host
```