## <font color='red'> 1.2 ArgoCD & Tekton </font>
Lab demonstrates a possible GitOps workflow using Argo CD and Tekton. We are using Argo CD to setup our Kubernetes clusters dev and prod (in the following we will only use the dev cluster) and Tekton to build and update our example application.

In this lab we're going to:
* install k8s-argocd cluster
* configure to pull image from Docker Hub
* configure to pull app from GitHub
* sync app on k8s-argocd cluster

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

Pre-requisties:
* Kustomize lets you lets you create an entire Kubernetes application out of individual pieces — without touching the YAML configuration filesfor the individual components.  For example, you can combine pieces from different sources, keep your customizations — or kustomizations, as the case may be — in source control, and create overlays for specific situations. And it is part of Kubernetes 1.14 or later. Kustomize enables you to do that by creating a file that ties everything together, or optionally includes “overrides” for individual parameters.

ensure Snap is installed:
```
sudo yum install epel-release
sudo yum install snapd
sudo systemctl enable --now snapd.socket
sudo ln -s /var/lib/snapd/snap /snap
```
install Kustomize:
```
sudo snap install kustomize
```

* GitOps is a way to do Kubernetes cluster management and application delivery.  It works by using Git as a single source of truth for declarative infrastructure and applications. With GitOps, the use of software agents can alert on any divergence between Git with what's running in a cluster, and if there's a difference, Kubernetes reconcilers automatically update or rollback the cluster depending on the case. With Git at the center of your delivery pipelines, developers use familiar tools to make pull requests to accelerate and simplify both application deployments and operations tasks to Kubernetes.

you will require a GitHub account.

* Docker Hub

---

#### <font color='red'> 1.1.1 Dev K8s Cluster </font>
The next step is to create ArgoCD Kubernetes cluster: 
* k8s-dev - install ArgoCD

start k8s-dev cluster:
```
minikube start -p k8s-dev
```
enable ingress:
```
minikube addons enable ingress -p k8s-dev
```
confirm that your k8s-dev context is set correctly:
```
kubectl config use-context k8s-dev
```

---


#### <font color='red'> 1.1.2 Install ArgoCD - Dev K8s Cluster </font>
install ArgoCD:
```
kustomize build clusters/argocd/dev | k apply -f -
```
verify that ArgoCD:
```
kgpo -n argocd
```
deploy our manifests to the cluster using the app of apps pattern. 
create a new Application, which manages all other applications (including ArgoCD):
```
k apply -f clusters/apps/dev.yaml
```
add our Ingresses to the /etc/hosts file:
```
sudo echo "`minikube ip --profile=k8s-dev` argocd-dev.fake grafana-dev.fake prometheus-dev.fake tekton-dev.fake server-dev.fake" | sudo tee -a /etc/hosts
```

  > open in browser: http://argocd-dev.fake