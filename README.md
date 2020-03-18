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

## Collaboration
My current skills on *golang* are very limited. Any contribution to speed up the implementation, suggestions aimed to improve the quality of the code and the repo organization are more then welcome.

I'm currently grabbing existing code samples from [multiple sources](notes.md) to quickly match the goal described above so prove the technical feasibility.

## Kubernetes Operator
I'm also working on a [k8s operator](operator/README.md) to deploy and monitor the described solution.