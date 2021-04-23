## <font color='red'> 2.2 ArgoCD </font>
Argo CD is a declarative, GitOps continuous delivery tool for Kubernetes.

In this lab we're going to:
* install Argo CD
* deploy a guestbook application

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

#### <font color='red'> 1.1.1 K8s Cluster </font>
The next step is to create Kubernetes cluster: 
* install minikube
* enable Ingress addon

start minikube:
```
minikube start
```
enable Ingress:
```
minikube addons enable ingress
```
verify ingress:
```
ksysgpo
```

---

#### <font color='red'> 1.1.1 Install ArgoCD </font>
The next step is to: 
* install ArgoCD

create a namespace:
```
kubectl create namespace argocd
```
install argocd:
```
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.0.0/manifests/install.yaml
```
verify deployed ArgoCD:
```
kgpo -n argocd
```
By default, the Argo CD API server is not exposed with an external IP.  

port-forward:
```
kubectl port-forward svc/argocd-server -n argocd 8080:443
```
the API server can then be accessed: 

  > in browser: http://localhost:8080

username: admin

to retieve the password:
```
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d && echo
```

---

#### <font color='red'> 1.1.2 Deploy a Guestbook App </font>

* guestbook - kubectl
* guestbook - kustomize
* guestbook - helm

verify status and configuration of your app

Notice: the STATUS: OutOfSync and HEALTH: Missing. That’s because ArgoCD creates applications with manual triggers by default.  

“Sync” is the terminology ArgoCD uses to describe the application on your target cluster as being up to date with the sources ArgoCD is pulling from. 
You have set up ArgoCD to monitor the GitHub repository with the configuration files. Once the initial sync is completed, a change will cause the status in ArgoCD to change to OutOfSync.

port-forward to expose app on localhost:9090:
```
kubectl port-forward svc/guestbook -n default 9090:8080
```


  > more examples can be found: https://github.com/argoproj/argocd-example-apps



clean up:

---

#### <font color='red'> 1.1.2 Deploy App of Apps pattern </font>
Declaratively specify one Argo CD app that consists only of other apps.