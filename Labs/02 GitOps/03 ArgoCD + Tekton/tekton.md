## <font color='red'> 3.1 ArgoCD + Tekton POC</font>
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

#### <font color='red'> 3.1.1 Install k3s Rancher</font>
k3d is a lightweight wrapper to run k3s (Rancher Lab’s minimal Kubernetes distribution) in docker.
k3d makes it very easy to create single- and multi-node k3s clusters in docker, e.g. for local development on Kubernetes.

This step is optional. If you already have a cluster, perfect, but if not, you can create a local one based on k3d.  
Ensure you're in the correct directory..

download k3d:
```
curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash
```
create k3d cluster:
```
./create-local-cluster.sh
```

---

#### <font color='red'> 3.1.2 Install Tekton + Argo CD</font>
The POC script:
* Installs Tekton + Argo CD, including secrets to access to Git repo
* Creates the volume and claim necessary to execute pipelines
* Deploys Tekton dashboard
* Deploys Sonarqube
* Deploys Nexus and configure an standard instance
* Creates the configmap associated to Maven settings.xml, ready to publish artifacts in Nexus (with user and password)
* Installs Tekton tasks and pipelines
* Git-clone (from Tekton Hub)
* Maven (from Tekton Hub)
* Buildah (from Tekton Hub) - builds OCI images
* Prepare Image (custom task: poc/conf/tekton/tasks/prepare-image-task.yaml)
* Push to GitOps repo (custom task: poc/conf/tekton/tasks/push-to-gitops-repo.yaml)
* Installs Argo CD application, configured to check changes in GitOps repository (resources/gitops_repo)

```
./setup-poc.sh
```
** Be patient. The process takes some minutes. Ignore the Nexus error. It will continue..  :)

---

#### <font color='red'> 3.1.3 Access Tekton + Argo CD + Tests</font>

**Tekton**  

the tests run:
* SonarQube® is an automatic code review tool to detect bugs, vulnerabilities, and code smells in your code. It can integrate with your existing workflow to enable continuous code inspection across your project branches and pull requests.


to access Tekton dashboard:
```
kubectl proxy --port=8080
```
access the pipline:

  > in browser: http://localhost:8080/api/v1/namespaces/tekton-pipelines/services/tekton-dashboard:http/proxy/#/namespaces/cicd/pipelineruns

watch the video..

to view the pods executing the pipeline:
```
kubectl get pods -n cicd -l "tekton.dev/pipelineRun=products-ci-pipelinerun"
```

**Sonarqube**

to access Sonarqube to check quality issues:

  > in browser: http://localhost:9000/projects

user: admin
password: admin123  

**Nexus**

access Nexus to check how the artifact has been published:

Nexus is a repository manager. It allows you to proxy, collect, and manage your dependencies so that you are not constantly juggling a collection of JARs. 

the last stage in CI part consist on performing a push action to GitOps repository. In this stage, content from GitOps repo is cloned, commit information is updated in cloned files (Kubernentes descriptors) and a push is done. 

watch the video..!

**ArgoCD**  

to access the ArgoD dashboard:
```
kubectl port-forward svc/argocd-server -n argocd 9090:8080
```

  > in browser: https://localhost:9090

user: admin
password: 
```
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d && echo
```

In this dashboard you should be the "product service" application that manages synchronization between Kubernetes cluster and GitOps repository.

The application is "healthy" but as the objects associated with Product Service (Pods, Services, Deployment,...etc) aren't still deployed to the Kubernetes cluster sync status is "unknown".

Once the "pipelinerun" ends and changes are pushed to GitOps repository, Argo CD compares content deployed in the Kubernetes cluster (associated to Products Service) with content pushed to the GitOps repository and synchronizes Kubernetes cluster against the repository.



clean up:
```
delete-local-cluster.sh
```

---
