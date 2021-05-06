## <font color='red'>2.4 Tekton</font>
Argo CD is a declarative, GitOps continuous delivery tool for Kubernetes.

In this lab we're going to:
* install Tekton
* install Tekton CLI
* install Tekton Dashboard on Kubernetes

* run the application tests inside the cloned git repository
* build a Docker image for our Go application and push it to DockerHub

**The second part requires a Docker Hub account.**

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

#### <font color='red'>2.4.1 Install Tekton + CLI</font>

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

  > view Tekton dashboard: http://localhost:9097

---

#### <font color='red'>2.4.2 Tekton Taskrun</font>
In our first tekton pipeline a Go application simply prints the sum of two integers.
* run the application tests inside the cloned git repository

The required resource files can be found at:

  > Tekton demo repository: http://github.com/jporeilly/Tekton-demo

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
      args: ["test"]
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
      value: https://github.com/jporeilly/Tekton-demo
    - name: revision
      value: main
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

apply the file with kubectl and then check the Pods and TaskRun resources. The Pod will go through the Init:0/2 and PodInitializing status and then succeed:
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

---

#### <font color='red'>2.4.3 Tekton Pipeline - Docker</font>
In our second tekton pipeline a Go application simply prints the sum of two integers.
* build a Docker image for our Go application and push it to DockerHub

**You will need a Docker Hub account**

To build and push our Docker image we use Kaniko, which can build Docker images inside a Kubernetes cluster without depending on a Docker daemon.

Kaniko will build and push the image in the same command. This means before running our task we need to set up credentials for DockerHub so that the docker image can be pushed to the registry.

The credentials are saved in a Kubernetes Secret. 

create a file named 04-secret.yaml with the following content and replace myusername and mypassword with your DockerHub credentials:
```
apiVersion: v1
kind: Secret
metadata:
  name: basic-user-pass
  annotations:
    tekton.dev/docker-0: https://index.docker.io/v1/
type: kubernetes.io/basic-auth
stringData:
    username: [myusername]
    password: [mypassword]
```
Note: the tekton.dev/docker-0 annotation in the metadata which tells Tekton the Docker registry these credentials belong to.

WARNING: This will write the unencypted credentials to /tekton/home.docker

Next we create a ServiceAccount that uses the basic-user-pass Secret. 

create a file named 05-serviceaccount.yaml with the following content:
```
apiVersion: v1
kind: ServiceAccount
metadata:
  name: build-bot
secrets:
  - name: basic-user-pass
```
apply both files with kubectl:
```
kubectl apply -f 04-secret.yaml
kubectl apply -f 05-serviceaccount.yaml
```
Now that the credentials are set up we can continue by creating the Task that will build and push the Docker image.

create a file called 06-task-build-push.yaml with the following content:
```
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: build-and-push
spec:
  resources:
    inputs:
      - name: repo
        type: git
  steps:
    - name: build-and-push
      image: gcr.io/kaniko-project/executor:v1.3.0
      env:
        - name: DOCKER_CONFIG
          value: /tekton/home/.docker
      command:
        - /kaniko/executor
        - --dockerfile=Dockerfile
        - --context=/workspace/repo/src
        - --destination=jporeilly/tekton-test:v1
```
Similarly to the first task this task takes a git repo as an input (the input name is repo) and consists of only a single step since Kaniko builds and pushes the image in the same command.

Make sure to create a DockerHub repository and replace jporeilly/tekton-test with your repository name. In this example it will always tag and push the image with the v1 tag.

Tekton has support for parameters to avoid hardcoding values like this. However to keep this tutorial simple I've left them out...  :)

The DOCKER_CONFIG env var is required for Kaniko to be able to find the Docker credentials.

apply the file with kubectl:
```
kubectl apply -f 06-task-build-push.yaml
```
There are two ways we can test this Task, either by manually creating a TaskRun definition and then applying it with kubectl or by using the Tekton CLI (tkn).

To run the Task with kubectl we create a TaskRun that looks identical to the previous with the exception that we now specify a ServiceAccount (serviceAccountName) to use when executing the Task.

create a file named 07-taskrun-build-push.yaml with the following content:
```
apiVersion: tekton.dev/v1beta1
kind: TaskRun
metadata:
  name: build-and-push
spec:
  serviceAccountName: build-bot
  taskRef:
    name: build-and-push
  resources:
    inputs:
      - name: repo
        resourceRef:
          name: tekton-example
```
apply the task and check the log of the Pod by listing all Pods that start with the Task name build-and-push:
```
kubectl apply -f 07-taskrun-build-push.yaml
```
check Pods:
```
kubectl get pods | grep build-and-push
```

To see the output of the containers we can run the following command. Make sure to replace build-and-push-pod-c698q with the the Pod name from the output above (it will be different for each run).
```
kubectl logs --all-containers build-and-push-pod-c698q --follow
```
the task executed without problems and we can now pull/run our Docker image:
```
docker run [docker-hub-username]/tekton-test:v1
```
Running the Task with the Tekton CLI is more convenient. With a single command it generates a TaskRun manifest from the Task definition, applies it, and follows the logs.

```
tkn task start build-and-push --inputresource repo=tekton-example --serviceaccount build-bot --showlog
```
Now that we have both of our Tasks ready (test, build-and-push) we can create a Pipeline that will run them sequentially: First it will run the application tests and if they pass it will build the Docker image and push it to DockerHub.

create a file named 08-pipeline.yaml with the following content:
```
apiVersion: tekton.dev/v1beta1
kind: Pipeline
metadata:
  name: test-build-push
spec:
  resources:
    - name: repo
      type: git
  tasks:
    # Run application tests
    - name: test
      taskRef:
        name: test
      resources:
        inputs:
          - name: repo      # name of the Task input (see Task definition)
            resource: repo  # name of the Pipeline resource

    # Build docker image and push to registry
    - name: build-and-push
      taskRef:
        name: build-and-push
      runAfter:
        - test
      resources:
        inputs:
          - name: repo      # name of the Task input (see Task definition)
            resource: repo  # name of the Pipeline resource
```
The first thing we need to define is what resources our Pipeline requires. A resource can either be an input or an output. In our case we only have an input: the git repo with our application source code. We name the resource repo.

Next we define our tasks. Each task has a taskRef (a reference to a Task) and passes the tasks required inputs.

apply the file with kubectl:
```
kubectl apply -f 08-pipeline.yaml
```
Similar to how we can run as Task by creating a TaskRun, we can run a Pipeline by creating a PipelineRun.

This can either be done with kubectl or the Tekton CLI. In the following two sections I will show both ways.

To run the file with kubectl we have to create a PipelineRun. 
create a file named 09-pipelinerun.yaml with the following content:
```
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata:
  name: test-build-push-pr
spec:
  serviceAccountName: build-bot
  pipelineRef:
    name: test-build-push
  resources:
  - name: repo
    resourceRef:
      name: tekton-example
```
apply the file, get the Pods that are prefixed with the PiplelineRun name, and view the logs to get the container output:
```
kubectl apply -f 09-pipelinerun.yaml
kubectl get pods | grep test-build-push-pr
```
To see the output of the containers we can run the following command. Make sure to replace test-build-push-pr-build-and-push-gh4f4-pod-nn7k7 with the the Pod name from the output above (it will be different for each run).
```
kubectl logs test-build-push-pr-build-and-push-gh4f4-pod-nn7k7 --all-containers --follow
```
When using the CLI we don't have to write a PipelineRun, it will be generated from the Pipeline manifest. By using the --showlog argument it will also display the Task (container) logs:
```
tkn pipeline start test-build-push --resource repo=tekton-example --serviceaccount build-bot --showlog
```

---


#### <font color='red'>2.4.4 Further Tekton Tasks</font>
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
tkn task start --showlog hello -n tasks
tkn task start --showlog -p person=James hello -n tasks
```
* view in Tekton dashboard

---

**multiple tasks**  
Your tasks can have more than one step. In this next example, you will change this Task to use two steps. The first one will write to a file, and the second one will output the content of that file. The steps will run in the order in which they are defined in the steps array.

First, start by adding a new step called write-hello. In here, you will use the same UBI base image. Instead of using a single command, you can also write a script. You can do this with a script parameter, followed by a | and the actual script to run. In this script, start by echoing "Preparing greeting", then echo the "Hello $(params.person)" that you had in the previous example into the ~/hello.txt file. Finally, add a little pause with the sleep command and echo "Done".

For the second step, you can create a new step called say-hello. This second step will run in its container but share the /tekton folder from the previous step. In the first step, you created a file in the "~" folder, which maps to "/tekton/home". For this second step, you can use an image node:14, and the file you created in the first step will be accessible. You can also run a NodeJS script as long as you specify the executable in the #! line of your script. In this case, you can write a script that will output the content of the ~/hello.txt file.
```
kubectl apply -f 03_multi-steps/step.yaml -n tasks
tkn task start --showlog hello -n tasks
```
* view in Tekton dashboard

---


#### <font color='red'>2.4.5 More Tekton Pipelines</font>
Tasks are useful, but you will usually want to run more than one Task. In fact, tasks should do one single thing so you can reuse them across pipelines or even within a single pipeline. For this next example, you will start by writing a generic task that will echo whatever it receives in the parameters.

List of Tekton Pipelines:
* hello
* run sequentially or parallel
* resources

> lots more examples: https://github.com/tektoncd/pipeline/tree/main/examples/v1beta1/pipelineruns

* ensure you're in the Pipelines directory..

create a namespace to run pipelines:
```
k create namespace pipelines 
```

---

**hello**  
A pipeline is a series of tasks that can run either in parallel or sequentially. In this Pipeline, you will use the say-something tasks twice with different outputs.

You can now apply the Task and this new Pipeline to your cluster and start the Pipeline. Using tkn pipeline start will create a PipelineRun with a random name. You can also see the logs of the Pipeline by using the --showlog parameter.

```
kubectl apply -f 01_hello/tasks.yaml -n pipelines
kubectl apply -f 01_hello/pipeline.yaml -n pipelines
tkn pipeline start say-things --showlog -n pipelines
```

* view in Tekton dashboard

---

**run sequentially or parallel**  
For Tasks to run in a specific order, the runAfter parameter is needed in the task definition of your Pipeline.
The runAfter parameter is being applied to specific numbered tasks, and after applying this Pipeline to our cluster, we’ll be able to see logs from each task, but ordered:

```
kubectl apply -f 02_para-seq/pipeline-order.yaml -n pipelines
tkn pipeline start say-things-in-order --showlog
```

* view in Tekton dashboard

---

**resources**  
The last object that will be demonstrated in this lab is PipelineResources. When you create pipelines, you will want to make them as generic as you can. This way, your pipelines can be reused across various projects. In the previous examples, we used pipelines that didn't do anything interesting. Typically, you will want to have some input on which you will want to perform your tasks. Usually, this would be a git repository. At the end of your Pipeline, you will also typically want some sort of output. Something like an image. This is where PipelineResources will come into play.

In this next example, you will create a pipeline that will take any git repository as a PipelineResource and then count the number of files.

* First, you can start by creating a task. This Task will be similar to the ones you've created earlier but will also have an input resource.
* Next, you can create a pipeline that will also have an input resource. This Pipeline will have a single task, which will be the count-files task you've just defined.
* Finally, you can create a PipelineResource. This resource is of type git, and you can put in the link of a Github repository in the url parameter. You can use the repo for this project.

```
kubectl apply -f pipeline-resource.yaml -n pipeline
tkn pipeline start count --showlog
tkn pipeline start count --showlog --resource git-repo=git-repo
```

clean up:
```
minikube stop
minikube delete
```

---
