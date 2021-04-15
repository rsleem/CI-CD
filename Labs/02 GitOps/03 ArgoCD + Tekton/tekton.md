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
**poc:** this is the main directory. contains 3 scripts:
* create-local-cluster.sh: this script creates a local Kubernetes cluster based on K3D.
* delete-local-cluster.sh: this script removes the local cluster
* setup-poc.sh: this script installs and configure everything neccessary in the cluster (Tekton, Argo CD, Nexus, Sonar, etc...)
  
**resources:** directory used to manage the two repositories (code and gitops):
* sources-repo: source code of the app 
* gitops-repo: repository where Kubernetes files associated to the service to be deployed are




---

## <font color='red'> 3.1.1 Install Rancher</font>
k3d is a lightweight wrapper to run k3s (Rancher Lab’s minimal Kubernetes distribution) in docker.
k3d makes it very easy to create single- and multi-node k3s clusters in docker, e.g. for local development on Kubernetes.

