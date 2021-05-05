## <font color='red'>3.1 PoC - Argo CD + Tekton</font>
Argo CD watches cluster objects stored in a Git repository and manages the create, update, and delete (CRUD) processes for objects within the repository. Tekton is a CI/CD tool that handles all parts of the development lifecycle, from building images to deploying cluster objects.

In this Lab you will:
* install k3s Rancher
* install Tekton + ArgoCD
* install Nexus
* install SonarQube


* create a CI Pipeline

For those of you that want to have some fun, SonarQube and Nexus have also been installed and configured.

**Sonarqube**
the tests run:
* SonarQube® is an automatic code review tool to detect bugs, vulnerabilities, and code smells in your code. It can integrate with your existing workflow to enable continuous code inspection across your project branches and pull requests.
to access Sonarqube to check quality issues:

  > in browser: http://localhost:9000/projects

user: admin
password: admin123  

**Nexus**
Nexus is a repository manager. It allows you to proxy, collect, and manage your dependencies so that you are not constantly juggling a collection of JARs.

access Nexus to check the artifact has been published:

  > in browser: http://localhost:9001

user: admin
password: admin123 


---

#### <font color='red'>Pre-requisites Install k3s Rancher</font>
k3d is a lightweight wrapper to run k3s (Rancher Lab’s minimal Kubernetes distribution) in docker.
k3d makes it very easy to create single- and multi-node k3s clusters in docker, e.g. for local development on Kubernetes.

This step is optional. If you already have a cluster, perfect, but if not, you can create a local one based on k3d.  
Ensure you're in the correct directory..

download k3d (this step has been completed):
```
curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash
```
create k3d cluster:
```
./create-local-cluster.sh
```

---

#### <font color='red'>Pre-requisites Install Tekton + Argo CD</font>
Theres a script:
* Installs Tekton + Argo CD, including secrets to access to Git repo
* Creates the volume and claim necessary to execute pipelines
* Deploys Tekton dashboard
* Deploys Sonarqube
* Deploys Nexus and configure an standard instance

run the script:
```
./setup-poc.sh
```
** Be patient. The process takes some minutes. Ignore the Nexus error. It will continue..  :)

* Creates the configmap associated to Maven settings.xml, ready to publish artifacts in Nexus (with user and password)
* Installs Tekton tasks and pipelines 
* Git-clone (from Tekton Hub)
* Maven (from Tekton Hub)
* Buildah (from Tekton Hub) - builds OCI images
* Prepare Image (custom task: poc/conf/tekton/tasks/prepare-image-task.yaml)
* Push to GitOps repo (custom task: poc/conf/tekton/tasks/push-to-gitops-repo.yaml)
* Installs Argo CD application, configured to check changes in GitOps repository (resources/gitops_repo)


---

#### <font color='red'>Access Argo CD + Tekton</font>
to access the ArgoD dashboard:
```
kubectl port-forward svc/argocd-server -n argocd 9070:443
```

  > in browser: https://localhost:9070

user: admin
password: 
```
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d && echo
```


to access Tekton dashboard:
```
kubectl patch service tekton-dashboard -n tekton-pipelines -p '{"spec": {"type": "LoadBalancer"}}'
```
verify external-ip:
```
kubectl get svc -n cicd
```
access the pipline:

  > in browser: http://[external-ip]:9097


---


### <font color='red'>CI-Pipeline</font>
This POC illustrates GitOps CI/CD pipelines. 

CI stages implemented by Tekton:
* Checkout: in this stage, source code repository is cloned
* Build & Test: in this stage, we use Maven to build and execute test
* Code Analisys: code is evaluated by Sonarqube
* Publish: if everything is ok, artifact is published to Nexus
* Build image: in this stage, we build the image and publish to local registry
* Push to GitOps repo: this is the final CI stage, in which Kubernetes descriptors are cloned from the GitOps repository, they are modified in order to insert commit info and then, a push action is performed to upload changes to GitOps repository.

CD stages implemented by ArgoCD:
* Argo CD detects that the repository has changed and perform the sync action against the Kubernetes cluster.

directory structure:  

**poc:**   
this is the main directory. contains 3 scripts:
* create-local-cluster.sh: this script creates a local Kubernetes cluster based on K3D.
* delete-local-cluster.sh: this script removes the local cluster
* setup-poc.sh: this script installs and configure everything neccessary in the cluster (Tekton, Argo CD, Nexus, SonarQube, etc...)
  
**resources:**   
directory used to manage the two repositories (code and gitops):
* sources-repo: source code of the app 
* gitops-repo: repository used for Kubernetes deployment YAML files.

---

### <font color='red'>Troubleshoot</font>
* ensure centos is the owner
```
cd tekton-argocd-poc
sudo chown -R centos:centos tekton-argocd-poc
```
* ensure the following files are +x
```
cd tekton-argocd-poc/poc
sudo chmod +x create-local-cluster.sh
sudo chmod +x delete-local-cluster.sh
sudo chmod +x setup-poc.sh
```

---

clean up:
```
delete-local-cluster.sh
```

---
