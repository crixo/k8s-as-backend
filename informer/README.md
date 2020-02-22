- init go module
use ```replace``` to point to the local folder instead of the remote one

- write biz logic code

- vendoring adding vendor code-related
```
go mod vendor
```

- build the docker image
```
docker build -t crixo/k8s-as-backend-informer:v.0.0.0 .
```

- laod image in kind
```
kind load docker-image crixo/k8s-as-backend-informer:v.0.0.0 --name standard
```

- apply k8s resources
```
k apply -f artifacts/crd.yaml,artifacts/app.yaml
sleep 10
k apply -f artifacts/todo.yaml  
```

- get resource by rest api call
```
k proxy
# list
http://127.0.0.1:8001/apis/k8sasbackend.com/v1/todos/
# list
curl http://127.0.0.1:8001/apis/k8sasbackend.com/v1/namespaces/default/todos
# single item
curl http://127.0.0.1:8001/apis/k8sasbackend.com/v1/namespaces/default/todos/buy-book
# filter list using filed selector
http://localhost:8001/apis/k8sasbackend.com/v1/namespaces/default/todos?fieldSelector=metadata.name%3Dbuy-book
# filter by label
http://localhost:8001/apis/k8sasbackend.com/v1/namespaces/default/todos?labelSelector=pippo%3Dpluto
```

## Enable field to be used w/ fieldSeletor
Not possible yet: [53459 Enable arbitrary CRD field selectors by supporting a whitelist of fields in CRD spec](https://github.com/kubernetes/kubernetes/issues/53459)

### How to identify fields supported by fieldSelector
https://kubernetes.slack.com/archives/C0X2V69D0/p1582371668432800
- [events](https://github.com/kubernetes/kubernetes/blob/a51d57459604577fb835a505c71a51fa45250d72/pkg/registry/core/event/strategy.go#L100)
- [pods](https://github.com/kubernetes/kubernetes/blob/a51d57459604577fb835a505c71a51fa45250d72/pkg/registry/core/pod/strategy.go#L213)
