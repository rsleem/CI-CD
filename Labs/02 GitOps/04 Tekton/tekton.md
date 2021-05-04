## <font color='red'>4.1 Tekton</font>
Argo CD is a declarative, GitOps continuous delivery tool for Kubernetes.

In this lab we're going to:
* install Tekton
* install Tekton CLI
* install Tekton Dashboard on Kubernetes

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

#### <font color='red'>4.3.1 Install Tekton + CLI</font>


install tekton pipeline:
```
kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
```
verify installation:
```
kubectl get pods --namespace tekton-pipelines
```

**Tekton CLI**
install Tekton CLI:
```
curl -LO https://github.com/tektoncd/cli/releases/download/v0.17.2/tkn_0.17.2_Linux_x86_64.tar.gz
sudo tar xvzf tkn_0.17.2_Linux_x86_64.tar.gz -C /usr/local/bin/ tkn
```
To run a CI/CD workflow, you need to provide Tekton a Persistent Volume for storage purposes. Tekton requests a volume of 5Gi with the default storage class by default. 

check available persistent volumes and storage classes:
```
kubectl get pv
kubectl get storageclasses
```

**Tekton Dashboard**

install Tekton Dashboard on a Kubernetes cluster:
```
kubectl apply --filename https://storage.googleapis.com/tekton-releases/dashboard/latest/tekton-dashboard-release.yaml
```
access the Dashboard is using kubectl port-forward:
```
kubectl --namespace tekton-pipelines port-forward svc/tekton-dashboard 9097:9097
```

---

#### <font color='red'>4.3.2 Tekton Piepline</font>
In our first tekton pipeline a Go application simply prints the sum of two integers.
* run the application tests inside the cloned git repository
* build a Docker image for our Go application and push it to DockerHub

The required resource files can be found at:

  > Tekton demo repository: http://github.com/jporeilly/Tekton-demo.git

create a file called 01-task-test.yaml with the following content:
```
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: test
spec:
  resources:
    inputs:
      - name: repo
        type: git
  steps:
    - name: run-test
      image: golang:1.14-alpine
      workingDir: /workspace/repo/src
      command: ["go"]
      args: ["test"]`
```
The resources: block defines the inputs that our task needs to execute its steps. Our step (name: run-test) needs the cloned tekton-demo git repository as an input and we can create this input with a PipelineResource.  

The git resource type will use git to clone the repo into the /workspace/$input_name directory everytime the Task is run. Since our input is named repo the code will be cloned to /workspace/repo. If our input would be named foobar it would be cloned into /workspace/foobar.

The next block in our Task (steps:) specifies the command to execute and the Docker image in which to run that command. We're going to use the golang Docker image as it already has Go installed.

For the go test command to run we need to change the directory. By default the command will run in the /workspace/repo directory but in our tekton-demo repo the Go application is in the src directory. We do this by setting workingDir: /workspace/repo/src.

Next we specify the command to run (go test) but note that the command (go) and args (test) need to be defined separately in the YAML file.

create a file called 02-pipelineresource.yaml:
```
apiVersion: tekton.dev/v1alpha1
kind: PipelineResource
metadata:
  name: tekton-example
spec:
  type: git
  params:
    - name: url
      value: https://github.com/jporeilly/tekton-demo
    - name: revision
      value: master
```

apply the Task and the PipelineResource with kubectl:
```
kubectl apply -f 01-task-test.yaml
kubectl apply -f 02-pipelineresource.yaml
```
To run our Task we have to create a TaskRun that references the previously created Task and passes in all required inputs (PipelineResource).

create a file called 03-taskrun.yaml with the following content:
```
apiVersion: tekton.dev/v1beta1
kind: TaskRun
metadata:
  name: testrun
spec:
  taskRef:
    name: test
  resources:
    inputs:
      - name: repo
        resourceRef:
          name: tekton-example
```
This will take our Task (taskRef is a reference to our previously created task name: test) with our tekton-demo git repo as an input (resourceRef is a reference to our PipelineResource name:: tekton-example) and execute it.

Apply the file with kubectl and then check the Pods and TaskRun resources. The Pod will go through the Init:0/2 and PodInitializing status and then succeed:
```
kubectl apply -f 03-taskrun.yaml
```
check Pods:
```
kgpo
```
check taskrun:
```
kg taskrun
```
To see the output of the containers we can run the following command. Make sure to replace testrun-pod-pds5z with the the Pod name from the output above (it will be different for each run).
```
kubectl logs testrun-pod-pds5z --all-containers
```
Our tests passed and our task succeeded. Next we will use the Tekton CLI to see how we can make this whole process easier.

Instead of manually writing a TaskRun manifest we can run the following command which takes our Task (named test), generates a TaskRun (with a random name) and shows its logs:
```
tkn task start test --inputresource repo=tekton-example --showlog
```









clean up:
```
delete-local-cluster.sh
```

---
