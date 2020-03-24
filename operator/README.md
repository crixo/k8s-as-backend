# K8sAsBackend Operator

## Dependencies
you have to install the following dependencies

- docker engine (eg. [docker for desktop](https://docs.docker.com/install/)  for windows and mac users)
- [kind](https://kind.sigs.k8s.io/) The *right way* to have k8s locally - v0.6.0 w/ k8s image v1.16.3
- [golang 1.13`*`](https://golang.org/doc/install) *THE* k8s programming language
- [operator-SDK`*`](https://github.com/operator-framework/operator-sdk) - v0.15.2 The framework selected to build k8s operator

`*` required only to simplify development and debugging activities. The final solution, as per [operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/), runs within a docker container inside the k8s cluster. 

## Demo steps
Start from the repo root, not form this operator folder. Each step starts from repo root

- Create cluster

A 3 node cluster (1 master, 2 workers) will be created with a nginx-controller configured and tested w/ dummy apps (foo and bar).  
The nginx-controller uses ports 80 and 443: make sure you do not have other services currently running on those port otherwise the installation will fail.
```
cd kind
# proposed default cluster name works just fine
sh create-cluster.sh 
# wait until done at least the cluster creation itself (~1 minute)
# nginx-ingress-controller installation&checks may take a bit longer
cd ..
```

- Deploy operator **CustomResourceDefinition**

```
cd operator
kubectl apply -f deploy/crds/k8s-as-backend.example.com_k8sasbackends_crd.yaml
cd ..
```

- Start the operator locally

Use a dedicated bash terminal. 
All the docker images currently part of the workload are stored on the public docker registry so kind cluster will be able to get it but you need an internet connection. To speed up the operator demo itself, I suggest you to use the pre-load images script.
```
#proposed default cluster name works just fine of course it has to match w/ the cluster creation step.
sh kind/preload-image.sh
```

Before running the app, make sure no PEM files are still present at ~/operator/certs.  You need at least a set of new PEM anytime you create a new k8s cluster. Certs folder clean up will be addressed in future releases.

The app will be in interactive mode sending logs to stdout
```
cd operator
operator-sdk run --local --namespace=default
```
You'll get some initial log reflecting the current operator configuration, then app waits for works to do through CR, the next step.

- Deploy your first operator **CustomResource**
Open a dedicated terminal tab pointing to repo root.
```
cd operator
kubectl apply -f deploy/crds/kab01.yaml
```
Go back to the terminal/tab where is running the operator app. You should see a lot of logs describing the tasks accomplished and the deployed resources. 

- Verify the CR has been updated by the deployment workflow adding the PEM just created or founded.
```
kubectl get k8sasbackends.k8s-as-backend.example.com  kab01 -o yaml
```

- Verify main secondary resources has been successfully deployed
```
kubectl get all
```

- Open the browser and test the app

Use the [TodoApp](http://localhost/default/kab01/todo-app/swagger-ui/index.html) that expose your business app. Create some todo and browse it through the swagger ui. *code* property value has to be unique within your app scope otherwise the request will fail.

- Check containers log to ensure the full workflow
```
kubectl get po
# check admission controller logs
kubectl logs kab01-todos-webhook-server-USE_YOUR_DEPLOYMENT_UNIQUE_IDENTIFIER
# check infomer logs
kubectl logs kab01-todo-app-USE_YOUR_DEPLOYMENT_UNIQUE_IDENTIFIER -c informer
```

## Browse the operator instance via api-server
Start api-server proxy into a dedicated terminal tab
```
kubectl proxy
```
then you can get [all the endpoints/paths exposed by the api-server](http://127.0.0.1:8001).  
The paths list includes k8s built-in paths and custom paths created through CRDs.

You can browse your CRDs definitions and CR instances through api server:
- [operator CRD](http://127.0.0.1:8001/apis/k8s-as-backend.example.com/v1alpha1)
- [operator CRs/instances](http://127.0.0.1:8001/apis/k8s-as-backend.example.com/v1alpha1/k8sasbackends) 
- [operator CRs/instances by namespace](http://127.0.0.1:8001/apis/k8s-as-backend.example.com/v1alpha1/namespaces/default/k8sasbackends)
- [operator CR/instance by namespace and name](http://127.0.0.1:8001/apis/k8s-as-backend.example.com/v1alpha1/namespaces/default/k8sasbackends/kab01)

- Clean up everything
```
cd operator
kubectl delete -f deploy/crds/kab01.yaml
```
At the present cluster-wide resource are not removed if no longer needed

## Notes
The images on docker hub may not be up to date. Use the specific folder at root level to build the latest version of the images.
