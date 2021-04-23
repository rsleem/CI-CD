## <font color='red'> 2.1 Skaffold </font>
GitOps is a way to do Kubernetes cluster management and application delivery.  
It works by using Git as a single source of truth for declarative infrastructure and applications. With GitOps, the use of software agents can alert on any divergence between Git with what's running in a cluster, and if there's a difference, Kubernetes reconcilers automatically update or rollback the cluster depending on the case. 
With Git at the center of your delivery pipelines, developers use familiar tools to make pull requests to accelerate and simplify both application deployments and operations tasks to Kubernetes.

In this lab we're going to:
* check kustomize
* check skaffold

* run various projects to illustrate the features of skaffold

list of projects:
* getting-started - deploys app with kubectl
* getting-started-kustomize - deploys app with kustomize to 3 environments (dev, staging, prod) 


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

**Pre-requistes:**
verify kustomize:
```
kustomize version
```
helm version:
```
helm version
```
verify skaffold:
```
skaffold version
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
minikube tunnel:
```
minikube tunnel
```

---

#### <font color='red'> 1.1.2 Skaffold </font>
Skaffold is a command line tool that facilitates continuous development for Kubernetes-native applications. Skaffold handles the workflow for building, pushing, and deploying your application, and provides building blocks for creating CI/CD pipelines. 

This enables you to focus on iterating on your application locally while Skaffold continuously deploys to your local or remote Kubernetes cluster.

skaffold commands:

  > open in browser: https://skaffold.dev/docs/references/cli/


**getting-started**  
Getting started with a simple go app.
This is a simple example based on:
* building a single Go file app and with a multistage Dockerfile using local docker to build
* tagging using the default tagPolicy (gitCommit)
* deploying a single container pod using kubectl

switch to getting-started directory:  

build a skaffold.yaml:
```
skaffold init
```
Note: Save the skaffold.yaml file.  There's a skaffold.yaml.bak just in case..!

Open the skaffold.yaml file
Notice its detected that the app will be deployed using kubectl.

deploy app:
```
skaffold dev --no-prune=false --cache-artifacts=false
```
Change the 'Hello world!' in the main.go

Ctrl+C will stop app. Note using the above flags will prune and delete the images.

clean up:
```
skaffold delete
```

**getting-started-kustomize**
This is a simple example based on:
* building a single Go file app and with a multistage Dockerfile using local docker to build
* tagging using the default tagPolicy (gitCommit)
* deploying a single container pod using kustomize

switch to getting-started-kustomize directory:  

build a skaffold.yaml:
```
skaffold init
```
Note: Save the skaffold.yaml file.  There's a skaffold.yaml.bak just in case..!

Open the skaffold.yaml file
Notice its detected that the app will be deployed using kustomize.

just deploy the app once:
```
skaffold run
```

clean up
```
skaffold delete
```

---