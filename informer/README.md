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
```