## <font color='red'>2.3 Argo CD + Tekton</font>
Argo CD watches cluster objects stored in a Git repository and manages the create, update, and delete (CRUD) processes for objects within the repository. Tekton is a CI/CD tool that handles all parts of the development lifecycle, from building images to deploying cluster objects.

In this Lab you will:
* install k3s Rancher
* install Tekton + ArgoCD
* install Nexus
* install SonarQube

* create a Task
* run the Task - Taskrun

---

#### <font color='red'>2.3.1 Install k3s Rancher</font>
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

#### <font color='red'>2.3.2 Install Tekton + Argo CD</font>
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

---

#### <font color='red'>2.3.3 Access Argo CD + Tekton</font>
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
kubectl patch service tekton-dashboard -n cicd -p '{"spec": {"type": "LoadBalancer"}}'
```
verify external-ip:
```
kubectl get svc -n cicd
```
access the pipline:

  > in browser: http://[external-ip]:9097

---

#### <font color='red'>2.3.4 Tekton Tasks</font>
List of Tekton Tasks:
* helloworld
* add a parameter
* multiple steps

> lots more examples: https://github.com/tektoncd/pipeline/tree/main/examples/v1beta1/taskruns

* ensure you're in the Tasks directory..

create a namespace to run tasks:
```
k create namespace tasks 
```

---

**hello world**
Simple Hello World example to show you how to:
* create a Task
* use a TaskRun to instantiate and execute a Task outside of a Pipeline

A Task defines a series of steps that run in a desired order and complete a set amount of build work. Every Task runs as a Pod on your Kubernetes cluster with each step as its own container. 


* view helloworld-task.yaml
to register the task:
```
kubectl apply -f helloworld-task.yaml -n tasks
```
details about your created Task:
```
tkn task describe echo-hello-world -n tasks
```
* view helloworld-taskrun.yaml
to run this task:
```
kubectl apply -f helloworld-taskrun.yaml -n tasks
```
check status:
```
tkn taskrun describe echo-hello-world-task-run -n tasks
```
* view in Tekton dashboard

---

**add a parameter**

Tasks can also take parameters. This way, you can pass various flags to be used in this Task. These parameters can be instrumental in making your Tasks more generic and reusable across Pipelines.

In this next example, you will create a task that will ask for a person's name and then say Hello to that person.

Starting with the previous example, you can add a params property to your task's spec. A param takes a name and a type. You can also add a description and a default value for this Task.

For this parameter, the name is person, the description is Name of person to greet, the default value is World, and the type of parameter is a string. If you don't provide a parameter to this Task, the greeting will be "Hello World".

You can then access those params by using variable substitution. In this case, change the word "World" in the args line to $(params.person).
```
kubectl apply -f 02_add-param/param.yaml -n tasks
tkn task start --showlog hello
tkn task start --showlog -p person=James hello
```

---

**multiple tasks**
Your tasks can have more than one step. In this next example, you will change this Task to use two steps. The first one will write to a file, and the second one will output the content of that file. The steps will run in the order in which they are defined in the steps array.

First, start by adding a new step called write-hello. In here, you will use the same UBI base image. Instead of using a single command, you can also write a script. You can do this with a script parameter, followed by a | and the actual script to run. In this script, start by echoing "Preparing greeting", then echo the "Hello $(params.person)" that you had in the previous example into the ~/hello.txt file. Finally, add a little pause with the sleep command and echo "Done".

For the second step, you can create a new step called say-hello. This second step will run in its container but share the /tekton folder from the previous step. In the first step, you created a file in the "~" folder, which maps to "/tekton/home". For this second step, you can use an image node:14, and the file you created in the first step will be accessible. You can also run a NodeJS script as long as you specify the executable in the #! line of your script. In this case, you can write a script that will output the content of the ~/hello.txt file.
```
kubectl apply -f 03_multi-steps/step.yaml -n tasks
tkn task start --showlog hello
```

---


#### <font color='red'>2.3.4 Tekton Pipelines</font>
Tasks are useful, but you will usually want to run more than one Task. In fact, tasks should do one single thing so you can reuse them across pipelines or even within a single pipeline. For this next example, you will start by writing a generic task that will echo whatever it receives in the parameters.

List of Tekton Pipelines:
* hello
* add a parameter
* multiple steps

> lots more examples: https://github.com/tektoncd/pipeline/tree/main/examples/v1beta1/taskruns

* ensure you're in the Pipelines directory..

create a namespace to run pipelines:
```
k create namespace pipelines 
```

---

**hello**
A pipeline is a series of tasks that can run either in parallel or sequentially. In this Pipeline, you will use the say-something tasks twice with different outputs.




---

## <font color='red'>ArgoCD + Tekton POC</font>
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


## <font color='red'>Pre-requsites</font>
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
* install tekton CLI
```
curl -LO https://github.com/tektoncd/cli/releases/download/v0.17.2/tkn_0.17.2_Linux_x86_64.tar.gz
sudo tar xvzf tkn_0.17.2_Linux_x86_64.tar.gz -C /usr/local/bin/ tkn
```

* Creates the configmap associated to Maven settings.xml, ready to publish artifacts in Nexus (with user and password)
* Installs Tekton tasks and pipelines - added later in POC
* Git-clone (from Tekton Hub)
* Maven (from Tekton Hub)
* Buildah (from Tekton Hub) - builds OCI images
* Prepare Image (custom task: poc/conf/tekton/tasks/prepare-image-task.yaml)
* Push to GitOps repo (custom task: poc/conf/tekton/tasks/push-to-gitops-repo.yaml)
* Installs Argo CD application, configured to check changes in GitOps repository (resources/gitops_repo)



to view the pods executing the pipeline:
```
kubectl get pods -n cicd -l "tekton.dev/pipelineRun=products-ci-pipelinerun"
```





The application is "healthy" but as the objects associated with Product Service (Pods, Services, Deployment,...etc) aren't still deployed to the Kubernetes cluster sync status is "unknown".

Once the "pipelinerun" ends and changes are pushed to GitOps repository, Argo CD compares content deployed in the Kubernetes cluster (associated to Products Service) with content pushed to the GitOps repository and synchronizes Kubernetes cluster against the repository.

In this dashboard you should be the "product service" application that manages synchronization between Kubernetes cluster and GitOps repository.


#### <font color='red'>Access Tekton + Argo CD + Tests</font>


**Sonarqube**
the tests run:
* SonarQube® is an automatic code review tool to detect bugs, vulnerabilities, and code smells in your code. It can integrate with your existing workflow to enable continuous code inspection across your project branches and pull requests.
to access Sonarqube to check quality issues:

  > in browser: http://localhost:9000/projects

user: admin
password: admin123  

**Nexus**
Nexus is a repository manager. It allows you to proxy, collect, and manage your dependencies so that you are not constantly juggling a collection of JARs.

access Nexus to check how the artifact has been published:

  > in browser: http://localhost:9001

user: admin
password: admin123 

the last stage in CI part consist on performing a push action to GitOps repository. In this stage, content from GitOps repo is cloned, commit information is updated in cloned files (Kubernentes descriptors) and a push is done. 

watch the video..!

---





clean up:
```
delete-local-cluster.sh
```

---
