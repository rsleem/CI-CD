## <font color='red'> 1.1 Kustomize </font>
Kustomize lets you lets you create an entire Kubernetes application out of individual pieces — without touching the YAML configuration filesfor the individual components.  For example, you can combine pieces from different sources, keep your customizations — or kustomizations, as the case may be — in source control, and create overlays for specific situations. 


And it is part of Kubernetes 1.14 or later. Kustomize enables you to do that by creating a file that ties everything together, or optionally includes “overrides” for individual parameters.

In this lab we're going to:
* install minikube
* check kustomize
* deploy nginx app with kubectl
* deploy app with kustomize

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


#### <font color='red'> 1.1.1 K8s Cluster </font>
The next step is to create Kubernetes cluster: 
* install minikube
* check kustomize
* deploy nginx app with kubectl

start minikube:
```
minikube start
```
minikube tunnel:
```
minikube tunnel
```
verify kustomize:
```
kustomize version
```

deploy nginx:
```
kubectl apply -f 01_nginx-deployment.yaml
```
show labels:
```
kubectl get all --show-labels
```
Note: the labels..

clean up:
```
kubectl delete -f 01_nginx-deployment.yaml
```

---

#### <font color='red'> 1.1.2 Kustomize Base </font>
A base is a kustomization referred to by some other kustomization.  
Any kustomization, including an overlay, can be a base to another kustomization.  
A base has no knowledge of the overlays that refer to it.  

For simple GitOps management, a base configuration could be the sole content of a git repository dedicated to that purpose. Same with overlays. Changes in a repo could generate a build, test and deploy cycle.

In this lab we're going to:
* configure kustomization.yaml
* deploy helloworld app
* verify deployment

* switch to hello-world directory.

tree the base directory:
```
tree base-hello
```
Examine each of the files to understand their relationship.


to view the concatenated output:
```
kustomize build base-hello
```
deploy app:
```
kustomize build base-hello | k apply -f -
```
verify deployment:
```
kg all -n hello
```
Note: make a note of the External IP of the service.  Each resource starts with the- prefix.

 > open in browser: http://Service-External-IP:8666

Note: Version 1: Good Morning!  These values are being pulled from the configmap. 


clean up:
```
kustomize build base-hello | k delete -f -
```
reset tunnel:
```
minikube tunnel cleanup
```

--- 

#### <font color='red'> 1.1.3 Kustomize Base + Overlays + Variants - helloworld </font>
An overlay is a kustomization that depends on another kustomization.  
The kustomizations an overlay refers to (via file path, URI or other method) are called bases.  
An overlay is unusable without its bases.  
An overlay may act as a base to another overlay.  

Overlays make the most sense when there is more than one, because they create different variants of a common base - e.g. development, QA, staging and production environment variants.  

These variants use the same overall resources, and vary in relatively simple ways, e.g. the number of replicas in a deployment, the CPU to a particular pod, the data source used in a ConfigMap, etc.  

in this lab we're going to:
* configure overlays - staging & production
* configure variants
* patching

**staging overlay**  

Patches can be used to apply different customizations to Resources. Kustomize supports different patching mechanisms through patchesStrategicMerge and patchesJson6902. patchesStrategicMerge is a list of file paths. Each file should be resolved to a strategic merge patch. 

The names inside the patches must match Resource names that are already loaded. Small patches that do one thing are recommended. For example, create one patch for increasing the deployment replica number and another patch for setting the memory limit.

switch to helloworld directory and tree:
```
tree
```
* view the map.yaml <font color='green'><em>#changing the configmap values</em></font>
* view the kustomization.yaml <font color='green'><em>#matched on labels and the patched & mergred with base/deployment.yaml</em></font>  

deploy staging variant:
```
kustomize build overlays/staging | kubectl apply -f -
```
verify deployment:
```
kg all -n hello
```
Note: make a note of the External IP of the service. Prefix is staging-.

 > open in browser: http://Service-External-IP:8666

Note: Version 1: <em>I have a pineapple! </em> These values are being pulled from the configmap. 


clean up:
```
kustomize build overlays/staging | k delete -f -
```
reset tunnel:
```
minikube tunnel cleanup
```


**production overlay**  

switch to helloworld directory and tree:
```
tree
```
* view the deployment.yaml    <font color='green'><em>#changing the replica count</em></font>
* view the kustomization.yaml <font color='green'><em>#matched on labels and the patched & merged with base/deployment.yaml</em></font>

deploy production variant:
```
kustomize build overlays/production | kubectl apply -f -
```
verify deployment:
```
kg all -n hello
```
Note: make a note of the External IP of the service. # of replicas = 5. Prefix is production- 

 > open in browser: http://Service-External-IP:8666

Note: Version 1: Good Morning!  These values are being pulled from the configmap. 

clean up:
```
kustomize build overlays/production | k delete -f -
```
reset tunnel:
```
minikube tunnel cleanup
```

compare the output directly to see how staging and production differ:
```
diff \
  <(kustomize build overlays/staging) \
  <(kustomize build overlays/production) |\
  more
```

  > For further examples: https://github.com/kubernetes-sigs/kustomize/tree/master/examples


---

#### <font color='red'> 1.1.4 Transformers </font>
Kustomize creates new resources by applying a series of transformations to an original set of resources. Kustomize provides the following default transformers:
  * annotations
  * images
  * labels
  * name reference
  * namespace
  * prefix/suffix
  * variable reference

Run through the transformers project and note how the transformer values declared in the patch file are mapped to the base.  This enables apps to be deployed against a 'template' base configuration.


**transformers**  
switch to transformers directory and tree:
```
tree
```
deploy transformers app:
```
kustomize build nginx | kubectl apply -f -
```
verify deployment:
```
kg all -n nginx
```
Note: you can expose the service or port-forward to the service.


clean up:
```
kustomize build nginx | kubectl delete -f -
```

---