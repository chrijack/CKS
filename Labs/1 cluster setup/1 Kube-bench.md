# [Kube-bench from Aquasecurity](https://github.com/aquasecurity/kube-bench)

Aqua Security Kube-bench is a tool that checks whether Kubernetes is deployed securely by running the checks documented in the [CIS Kubernetes Benchmark](https://www.cisecurity.org/benchmark/kubernetes/).

 

### 1. Installation

```bash
curl -L https://github.com/aquasecurity/kube-bench/releases/download/v0.6.15/kube-bench_0.6.15_linux_amd64.tar.gz -o kube-bench_0.6.15_linux_amd64.tar.gz
tar -xvf kube-bench_0.6.15_linux_amd64.tar.gz
```

 

### 2. Execute Kubebench on the cluster

 

```bash
sudo ./kube-bench --config-dir `pwd`/cfg --config `pwd`/cfg/config.yaml
# .....
== Summary total ==
65 checks PASS
10 checks FAIL
56 checks WARN
0 checks INFO
```
### 3. Check the failed tests on the cluster 

 

```bash
$ ./kube-bench --config-dir `pwd`/cfg --config `pwd`/cfg/config.yaml | grep "\[FAIL\] "

[FAIL] 1.1.12 Ensure that the etcd data directory ownership is set to etcd:etcd (Automated)
[FAIL] 1.2.5 Ensure that the --kubelet-certificate-authority argument is set as appropriate (Automated)
[FAIL] 1.2.17 Ensure that the --profiling argument is set to false (Automated)
[FAIL] 1.2.18 Ensure that the --audit-log-path argument is set (Automated)
[FAIL] 1.2.19 Ensure that the --audit-log-maxage argument is set to 30 or as appropriate (Automated)
[FAIL] 1.2.20 Ensure that the --audit-log-maxbackup argument is set to 10 or as appropriate (Automated)
[FAIL] 1.2.21 Ensure that the --audit-log-maxsize argument is set to 100 or as appropriate (Automated)
[FAIL] 1.3.2 Ensure that the --profiling argument is set to false (Automated)
[FAIL] 1.4.1 Ensure that the --profiling argument is set to false (Automated)
[FAIL] 4.1.1 Ensure that the kubelet service file permissions are set to 600 or more restrictive (Automated)
```


### 4. Fix failed test for 1.3.2. 
Edit `sudo nano /etc/kubernetes/manifests/kube-controller-manager.yaml` and add `--profiling=false`

```
1.3.2 Edit the Controller Manager pod specification file /etc/kubernetes/manifests/kube-controller-manager.yaml
on the control plane node and set the below parameter.
--profiling=false
```  
Example YAML:                                                                  
 ```
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    component: kube-controller-manager
    tier: control-plane
  name: kube-controller-manager
  namespace: kube-system
spec:
  containers:
  - command:
    - kube-controller-manager
    - --profiling=false
    - --allocate-node-cidrs=true
    - --authentication-kubeconfig=/etc/kubernetes/controller-manager.conf
...
```
### 5. Rerun kube-bench and verify 1.3.2 no longer fails.

```bash
sudo ./kube-bench --config-dir `pwd`/cfg --config `pwd`/cfg/config.yaml | grep "\[FAIL\] "
[FAIL] 1.1.12 Ensure that the etcd data directory ownership is set to etcd:etcd (Automated)
[FAIL] 1.2.5 Ensure that the --kubelet-certificate-authority argument is set as appropriate (Automated)
[FAIL] 1.2.17 Ensure that the --profiling argument is set to false (Automated)
[FAIL] 1.2.18 Ensure that the --audit-log-path argument is set (Automated)
[FAIL] 1.2.19 Ensure that the --audit-log-maxage argument is set to 30 or as appropriate (Automated)
[FAIL] 1.2.20 Ensure that the --audit-log-maxbackup argument is set to 10 or as appropriate (Automated)
[FAIL] 1.2.21 Ensure that the --audit-log-maxsize argument is set to 100 or as appropriate (Automated)
[FAIL] 1.4.1 Ensure that the --profiling argument is set to false (Automated)
[FAIL] 4.1.1 Ensure that the kubelet service file permissions are set to 600 or more restrictive (Automated)
```