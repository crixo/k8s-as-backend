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
#proposed default cluster name works just fine
sh create-cluster.sh 
# wait until done ~1 minute
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
kubectl apply -f deploy/crds/k8s-as-backend.example.com_v1alpha1_k8sasbackend_cr.yaml
```
Go back to the terminal/tab where is running the operator app. You should see a lot of logs describing the tasks accomplished and the deployed resources. 

- Verify the CR has been updated by the deployment workflow adding the PEM just created or founded.
```
kubectl get k8sasbackends.k8s-as-backend.example.com  example-k8sasbackend -o yaml
```

- Open the browser and test the app

Use the [TodoApp](http://localhost/default/example-k8sasbackend/todo-app/swagger-ui/index.html) that expose your business app. Create some todo and browse it through the swagger ui. *code* property value has to be unique within your app scope otherwise the request will fail.

- Check containers log to ensure the full workflow
```
k get po
# check admission controller logs
k logs k8s-as-backend-webhook-server-USE_YOUR_DEPLOYMENT_UNIQUE_IDENTIFIER
# check infomer logs
k logs todo-app-USE_YOUR_DEPLOYMENT_UNIQUE_IDENTIFIER -c informer
```

- Clean up everything
```
cd operator
kubectl delete -f deploy/crds/k8s-as-backend.example.com_v1alpha1_k8sasbackend_cr.yaml
```
At the present cluster-wide resource are not removed if no longer needed

## Notes
The images on docker hub may not be up to date. Use the specific folder at root level to build the latest version of the images.
