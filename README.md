English|[Russian](README_RU.md)

# Hello-operator

Simple kubernetes operator implemented using [client-go](https://github.com/kubernetes/client-go) and [code-generator](https://github.com/kubernetes/code-generator)  

When applying our CRD controller create job, that launches pod, which outputs a string specified in our CRD in parameter `message`

## Set up development environment

* `go mod vendor` to generate a `vendor/` folder with dependencies - required for use code-generator
* Run `/hack/update_codegen.sh`  to generate the CRD boilerplate code. In our case it is generated and located in the `/pkg/generated` folder  

## Launch
To set up a local cluster, you can use [Kind](https://kind.sigs.k8s.io/)
1. Build controller: `go build -o hello-controller . `
2. Add CRD to our k8s cluster `kubectl create -f crds/echo.yaml`
3. Run controller `./hello-controller -kubeconfig=path-to-your-cluser-config.yaml`   
Runs not as a pod in cluster, run directly on our local machine, remotely connecting to the cluster via config.
4. Apply our example CRD `kubectl create -f crds/examples/echo.yaml`

## Helpful links
[Official controller sample](https://github.com/kubernetes/sample-controller)  
[Operator Sample](https://github.com/mmontes11/echoperator)  
[Kubernetes notes (Chinese)](https://github.com/huweihuang/kubernetes-notes)  
[Operator SDK](https://github.com/kubernetes-sigs/kubebuilder)  
[Deep dive kubernetes code-generat0r](https://www.redhat.com/en/blog/kubernetes-deep-dive-code-generation-customresources)
