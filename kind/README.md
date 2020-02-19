# Nginx ingress controller for Kind

- create the cluster with the specific configuration
```
kind create cluster --config 3nodes-ingress-controller.yaml --name nginx-contr-3nodes
```

- [mandatory ingress-nginx components](https://kubernetes.github.io/ingress-nginx/deploy/#prerequisite-generic-deployment-command)
https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.27.0/deploy/static/mandatory.yaml
```
kubectl apply -f mandatory.yaml
```


- [expose the nginx service using NodePort](https://kubernetes.github.io/ingress-nginx/deploy/#bare-metal)
https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.27.0/deploy/static/provider/baremetal/service-nodeport.yaml
```
kubectl apply -f service-nodeport.yaml
```

- patch nginx ingress controller deployment
```
kubectl patch deployments -n ingress-nginx nginx-ingress-controller --patch "$(cat nginx-ingress-controller-deployment-patch.yaml)"
```

- deploy the k8s resources, including ingress, to test the configuration
https://kind.sigs.k8s.io/examples/ingress/usage.yaml
```
kubectl apply -f usage.yaml
```
- test from host machine
```
# should output "foo"
curl localhost/foo
# should output "bar"
curl localhost/bar
```