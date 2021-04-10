## <font color='red'> 2.1 GitOps </font>
GitOps is a way to do Kubernetes cluster management and application delivery.  
It works by using Git as a single source of truth for declarative infrastructure and applications. With GitOps, the use of software agents can alert on any divergence between Git with what's running in a cluster, and if there's a difference, Kubernetes reconcilers automatically update or rollback the cluster depending on the case. 
With Git at the center of your delivery pipelines, developers use familiar tools to make pull requests to accelerate and simplify both application deployments and operations tasks to Kubernetes.

In this lab we're going to:
* install skaffold
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


#### <font color='red'> 1.1.1 K8s Cluster </font>
The next step is to create Kubernetes cluster: 
* install minikube
* check kustomize

start minikube:
```
minikube start
```
add ingress:
```
minikube addons enable ingress
```
minikube tunnel:
```
minikube tunnel
```
verify kustomize:
```
kustomize version
```

---

#### <font color='red'> 1.1.2 Skaffold </font>
Skaffold is a command line tool that facilitates continuous development for Kubernetes-native applications. Skaffold handles the workflow for building, pushing, and deploying your application, and provides building blocks for creating CI/CD pipelines. 

This enables you to focus on iterating on your application locally while Skaffold continuously deploys to your local or remote Kubernetes cluster.

In this lab we're going to:
* install Skaffold
* download a sample go app

* Use skaffold dev to build and deploy your app every time your code changes,
* Use skaffold run to build and deploy your app once, similar to a CI/CD pipeline

install skaffold:
```
sudo snap install skaffold
```
check installed:
```
skaffold version
```
skaffold commands:

  > open in browser: https://skaffold.dev/docs/references/cli/










--- 

#### <font color='red'> 1.1.3 Kustomize Base + Overlays + Variants </font>
An overlay is a kustomization that depends on another kustomization.  
The kustomizations an overlay refers to (via file path, URI or other method) are called bases.  
An overlay is unusable without its bases.  
An overlay may act as a base to another overlay.  

Overlays make the most sense when there is more than one, because they create different variants of a common base - e.g. development, QA, staging and production environment variants.  

These variants use the same overall resources, and vary in relatively simple ways, e.g. the number of replicas in a deployment, the CPU to a particular pod, the data source used in a ConfigMap, etc.  

in this lab we're going to:
* configure overlays - staging & production
* configure variants
* patching

**staging overlay**

switch to staging overlay directory and tree:
```
tree staging
```
view the map.yaml # changing the configmap values

switch to helloworld directory:
```
kustomize build overlays/staging | kubectl apply -f -
```
verify deployment:
```
kg all
```
Note: make a note of the External IP of the service.

 > open in browser: http://Service-External-IP:8666

Note: Version 1: Good Morning!  These values are being pulled from the configmap. 


**production overlay**

switch to production overlay directory and tree:
```
tree production
```
* view the deployment.yaml #changing the replica count
* view the kustomization. yaml #matched on labels and the patched & mergred with base/deployment.yaml

switch to 01 kustomize directory:
```
kustomize build helloworld/overlays/production | kubectl apply -f -
```
verify deployment:
```
kg all
```
Note: make a note of the External IP of the service.

 > open in browser: http://Service-External-IP:8666

Note: Version 1: Good Morning!  These values are being pulled from the configmap. 

compare the output directly to see how staging and production differ:
```
diff \
  <(kustomize build helloworld/overlays/staging) \
  <(kustomize build helloworld/overlays/production) |\
  more
```

  > For further examples: https://github.com/kubernetes-sigs/kustomize/tree/master/examples


---