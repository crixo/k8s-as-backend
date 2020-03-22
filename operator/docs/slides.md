# Slides

- What is an operator?
You should picturing as a (virtual) Site Reliability Engineering(SRE):  
> SRE began at Google in response to the challenges of running huge systems with ever-increasing numbers of users and features. A key SRE objective is allowing services to grow without forcing the teams that run them to grow in direct proportion. To run systems at dramatic scale without a team of dramatic size, SREs write code to do deployment, operations, and maintenance tasks. SREs create software that runs other software, keeps it running, and manages it over time. [Kubernetes Operators](http://shop.oreilly.com/product/0636920234357.do)

- Why Operator as (GO) app  

  - Operator as Helm or Ansible

  - Ensure variable references:

      - Webhook-Server service name within certificate request

      - TodoApp service name within Webhook-Server deployment env var

- Monitoring many primitives k8s resources vs monitoring a single CR
  
  - use kubectl against new CRs

  - build your own dashboard against the api-server custom endpoints

- Operator-SDK vs Shared informer
  
  - Dependencies pinnings and cluster-based dependencies

  - scaffolding and code generation

- golang vs other languages

  - fully-functional client

  - same codebase used by k8s itself

  - frameworks available

  - same cluster lifecycle

- CRD for serverless approach  
deploying a CRD through api(server) you spawn a pod/container just for the requested job
