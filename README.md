# Kubernetes as backend

## Project Scope
The goals of this project/repo are the followings:
- Get hands-on with k8s CRD investigating what is needed to use CRD in a real production scenario.
- Understanding the k8s API model with a special attention to its extensibility
- Assuming you decided to use k8s as the hosting platform for your application workloads, be challenged on the idea of using k8s framework to store and share entities managed by your business application running in k8s.

According to the so called microservices approach, independent applications running within an applications ecosystem need to store their (owned) entities and share it with other applications part of the same ecosystem.  
Usually to accomplish this goal each application exposes a set of (REST) API to receive entities data and persist it into a dedicated/owned storage.  
Before persisting these entities, some business rules are applied to accept it and after the entities are persisted, the related data are broadcast within the ecosystem the application belongs to in case other applications are interested to these data/entities as well.  A message broker part of the ecosystem is usually adopted to share these resources across the applications.

Goal of this project is to accomplish the scenario described above leveraging on the built-in k8s features in order to implement the architecture described in the following diagram
![](images/k8s-as-backend.png?raw=true)

The next goal is collecting the potential issues of this approach, although technically feasible, highlighting the disadvantages and the risks related to it.

## API Strategy
Sync one-time execution vs async loop reconciliation.
![](images/kab-API-strategy.png?raw=true)

Business logic applied by the k8s controller should be delegated to a queue and executed form there. K8s provides its own [queueing mechanism](https://godoc.org/k8s.io/client-go/util/workqueue).

I'm currently grabbing existing code samples from [multiple sources](notes.md) to quickly match the goal described above so prove the technical feasibility.

## Manual deploy
with plain yaml
```
cd kind

# create cluster
sh create-cluster.sh

# deploy workloads
# TODO: verify the 'y' parameter due to docker build context path within nested scripts
sh deploy.sh k8s-as-backend y
```

with the operator using the e2e test framework with not clean-ip
```
cd operator

sh run-all.sh
```


## Kubernetes Operator
I'm also working on a [k8s operator](operator/README.md) to deploy and monitor the described solution.

Project resources dependencies vs operator scope
![](images/kab-resource-deps.png?raw=true)

## Reference

### CRD & Informer
- [Extending Kubernetes](https://get.oreilly.com/ind_extending-kubernetes.html)
- [gianarb/todo-crd](https://github.com/gianarb/todo-crd)
- [Code Generator](https://github.com/kubernetes/code-generator)

### Admission Controller

- [Dynamic Admission Control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers)
- [Admission Webhook server sample](https://github.com/kubernetes/kubernetes/tree/v1.16.11/test/images/agnhost#webhook-kubernetes-external-admission-webhook)
- [Banzai Tutorial](https://banzaicloud.com/blog/k8s-admission-webhooks/)

### Multi-container POD patterns

- [The Distributed System ToolKit: Patterns for Composite Containers](https://kubernetes.io/blog/2015/06/the-distributed-system-toolkit-patterns/)
- [Designing Distributed Systems](https://learning.oreilly.com/library/view/designing-distributed-systems/9781491983638/)
- [Distributed Application Runtime](https://dapr.io/)

### Operator

- [Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [Operator-SDK](https://github.com/operator-framework/operator-sdk)
- [Kubernetes Operators](https://learning.oreilly.com/library/view/kubernetes-operators/9781492048039/)

## Collaboration
My current skills on *golang* are very limited. Any contribution to speed up the implementation, suggestions aimed to improve the quality of the code and the repo organization are more then welcome.