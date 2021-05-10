## <font color='red'> 2.3 Argo CD </font>
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

the next step is to create Kubernetes cluster: 
* install minikube

start minikube:
```
minikube start
```
start tunnel:
```
minikube tunnel
```

---

#### <font color='red'> 2.3.1 Install ArgoCD </font>
The next step is to: 
* install ArgoCD

create a namespace:
```
kubectl create namespace argocd
```
install argocd:
```
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```
watch the Pods:
```
watch kubectl get pods -n argocd
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

#### <font color='red'> 2.3.2 Deploy a Guestbook App </font>
There are a couple of apps you can deploy at: http://github.com/jporeilly/ArgoCD-demos.git
* guestbook - kubectl
* guestbook - kustomize
* guestbook - helm

In this lab were going to deploy the guestbook - kubectl app.

**watch the video to see how to add an app** 

Notice: the STATUS: OutOfSync and HEALTH: Missing. That’s because ArgoCD creates applications with manual triggers by default.  

“Sync” is the terminology ArgoCD uses to describe the application on your target cluster as being up to date with the sources ArgoCD is pulling from. 
You have set up ArgoCD to monitor the GitHub repository with the configuration files. Once the initial sync is completed, a change will cause the status in ArgoCD to change to OutOfSync.

get the guestbook ClusterIP:
```
kg all -n guestbook -o wide
```
guestbook can then be accessed at:

 > in browser: http://[svc-ClusterIP]  

</br>
 
 > more examples can be found: https://github.com/argoproj/argocd-example-apps


clean up:

to stop  minikube:
```
minikube stop
```
to delete  minikube:
```
minikube delete
```

---