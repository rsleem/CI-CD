## <font color='red'> 1.1 Kustomize </font>
Kustomize lets you lets you create an entire Kubernetes application out of individual pieces — without touching the YAML configuration filesfor the individual components.  For example, you can combine pieces from different sources, keep your customizations — or kustomizations, as the case may be — in source control, and create overlays for specific situations. 


And it is part of Kubernetes 1.14 or later. Kustomize enables you to do that by creating a file that ties everything together, or optionally includes “overrides” for individual parameters.

In this lab we're going to:
* install minikube
* check kustomize
* sync apps on k8s-argocd cluster
* metrics in Grafana & Prometheus

---

#### <font color='red'>IMPORTANT:</font> 
<strong>Please ensure you start with a clean environment. 
If you have previously run minikube, you will need to delete the existing instance.</strong>

to stop  minikube:
```
minikube stop
```
to delete  minikube:
```
minikube delete
```

---







* GitOps is a way to do Kubernetes cluster management and application delivery.  It works by using Git as a single source of truth for declarative infrastructure and applications. With GitOps, the use of software agents can alert on any divergence between Git with what's running in a cluster, and if there's a difference, Kubernetes reconcilers automatically update or rollback the cluster depending on the case. With Git at the center of your delivery pipelines, developers use familiar tools to make pull requests to accelerate and simplify both application deployments and operations tasks to Kubernetes.

obviously you will require a GitHub account if you want to try yourself.

* Docker Hub


---

#### <font color='red'> 1.1.1 K8s Cluster </font>
The next step is to create k8s-dev Kubernetes cluster: 
* install minikube
* check kustomize

start minikube:
```
minikube start
```
minikube tunnel:
```
minikube tunnel
```
verify kustomize:
```
k version
```

deploy nginx:
```
kubectl apply -f 01_nginx-deployment.yaml
```
show labels:
```
kubectl get all --show-lables
```
Note: the labels..

clean up:
```
kubectl delete -f 01_nginx-deployment.yaml
```

---

#### <font color='red'> 1.1.2 Kustomize </font>

In this lab we're going to:
* configure kustomization.yaml
* deploy helloworld app
* verify deployment

switch to helloworld directory.

tree the base directory:
```
tree base
```
to view the concatenated output:
```
kustomize build base
```
switch to 01 kustomize directory:
```
kustomize build helloworld/base | k apply -f -
```
verify deployment:
```
kg all
```
Note: make a note of the External IP of the service.

  > in browser: http://[External-IP]:8666

Note: Version 1: Good Morning!  These values are being pulled from the configmap. 

