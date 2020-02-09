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