# Verify Platform Binaries

Kubernetes provides binaries and their checksum to verify authenticity of the files.

Check the Kubernetes downloads page (https://www.downloadkubernetes.com/)


### 1. Download binary file and checksum 

Download kubectl https://dl.k8s.io/v1.27.3/bin/darwin/amd64/kubectl and verify checksum matches 

### 2. Validate checksum
```
$ curl -LO https://dl.k8s.io/release/v1.27.3/bin/linux/amd64/kubectl
$ curl -LO https://dl.k8s.io/v1.27.3/bin/darwin/amd64/kubectl.sha256
$ echo "$(cat kubectl.sha256)  kubectl" | sha256sum --check
# kubectl: OK
```
https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/
