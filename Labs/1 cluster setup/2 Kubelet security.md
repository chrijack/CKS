# [Kubelet Security](https://kubernetes.io/docs/reference/access-authn-authz/kubelet-authn-authz/)

A kubelet's HTTPS endpoint exposes APIs which give access to data of varying sensitivity, and allow you to perform operations with varying levels of power on the node and within containers.


### 1. Check Kubelet configuration

--config= will show where the configuration file is 

```bash
$ ps -ef | grep kubelet # check --config parameter

root        1067       1  2 04:05 ?        00:01:07 /usr/bin/kubelet --bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf --config=/var/lib/kubelet/config.yaml --container-runtime-endpoint=unix:///var/run/crio/crio.sock --pod-infra-container-image=registry.k8s.io/pause:3.9 --node-ip=10.0.0.10
```

### 2. View the kubelet configuration file `sudo /var/lib/kubelet/config.yaml` on Control plane

```yaml
apiVersion: kubelet.config.k8s.io/v1beta1
authentication:
  anonymous:
    enabled: false # Disable annoymous auth.
  webhook: # Set to webhook
    cacheTTL: 0s
    enabled: true
  x509:
    clientCAFile: /etc/kubernetes/pki/ca.crt
authorization:
  mode: Webhook # webhook, instead of AlwaysAllow
  webhook:
    cacheAuthorizedTTL: 0s
    cacheUnauthorizedTTL: 0s
clusterDNS:
- 172.17.0.10
clusterDomain: cluster.local
cpuManagerReconcilePeriod: 0s
evictionPressureTransitionPeriod: 0s
...
```

### 3. Check that key and certificate are set in `kube-apiserver.yaml` file on control-plane: `sudo cat /etc/kubernetes/manifests/kube-apiserver.yaml`

```bash
$ sudo cat /etc/kubernetes/manifests/kube-apiserver.yaml | grep kubelet-client    

- --kubelet-client-certificate=/etc/kubernetes/pki/apiserver-kubelet-client.crt
- --kubelet-client-key=/etc/kubernetes/pki/apiserver-kubelet-client.key
```    

### 4. Test authentication
Without authentication
```bash
$ curl -sk https://localhost:10250/pods/
# Unauthorized
```
With authentication
```bash
$ sudo curl -sk https://localhost:10250/pods/ --key /etc/kubernetes/pki/apiserver-kubelet-client.key --cert /etc/kubernetes/pki/apiserver-kubelet-client.crt

# {"kind":"PodList","apiVersion":"v1","metadata":{},"items":[{"metadata":{"name":"etcd-controlplane","namespace": ...
```