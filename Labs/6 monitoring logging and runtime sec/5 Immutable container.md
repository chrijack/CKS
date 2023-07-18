You are an administrator of a Kubernetes cluster running a couple of existing Pods. It's your job to inspect the containers defined by the Pods for immutability. Delete all Pods that do not follow typical immutability best practices.

1. Create the objects from the file [`setup.yaml`](./setup.yaml).
2. List the Pods in the namespace `h92`.
3. Identify the Pods in the namespace that cannot be considered to run immutable containers.
4. Delete the Pods with mutable containers from the namespace. Which of the Pods are left running?

apiVersion: v1
kind: Namespace
metadata:
  name: h92
---
apiVersion: v1
kind: Pod
metadata:
  name: loop
  namespace: h92
spec:
  containers:
  - name: loop
    image: alpine:3.13.4
    args:
    - /bin/sh
    - -c
    - while true; do sleep 5; done;
    securityContext:
      privileged: true
  restartPolicy: Never
---
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  namespace: h92
spec:
  containers:
  - name: nginx
    image: nginx:1.21.6
    securityContext:
      readOnlyRootFilesystem: true
    volumeMounts:
    - name: nginx-run
      mountPath: /var/run
    - name: nginx-cache
      mountPath: /var/cache/nginx
    - name: nginx-data
      mountPath: /usr/local/nginx
  volumes:
  - name: nginx-run
    emptyDir: {}
  - name: nginx-data
    emptyDir: {}
  - name: nginx-cache
    emptyDir: {}
  restartPolicy: Never
---
apiVersion: v1
kind: Pod
metadata:
  name: hello-world
  namespace: h92
spec:
  securityContext:
    runAsUser: 0
  containers:
  - name: nginx
    image: hello-world:linux
  restartPolicy: Never
---
apiVersion: v1
kind: Pod
metadata:
  name: busybox
  namespace: h92
spec:
  hostNetwork: true
  containers:
  - name: busybox
    image: yauritux/busybox-curl:latest
    args:
    - /bin/sh
    - -c
    - curl google.com;
  restartPolicy: Never

  # Solution

## Setting up the Pods

Start by creating the objects from the existing YAML manifest.

```
$ kubectl apply -f setup.yaml
namespace/h92 created
pod/loop created
pod/nginx created
pod/hello-world created
pod/busybox created
```

List all Pods in the namespace `h92`.

```
$ kubectl get pods -n h92
NAME          READY   STATUS      RESTARTS   AGE
busybox       0/1     Completed   0          14s
hello-world   0/1     Completed   0          14s
loop          1/1     Running     0          14s
nginx         1/1     Running     0          14s
```

## Identifying Pods with Mutable Containers

You can inspect the configuration of the live objects using the `kubectl get pod` command with the option `-o yaml`. For example, to inspect the configuration of the `busybox` Pod, run the command `kubectl get pod busybox -n h92 -o yaml`.

The following table explains why a specific Pod is immutable or mutable.

| Pod             | Mutability  | Reason |
| :-------------- | :---------- | :----- |
| `busybox`       | mutable     | The container uses `hostNetwork` with the value `true`. |
| `hello-world`   | mutable     | The container uses `securityContext.runAsUser` with the value `0`. |
| `loop`          | mutable     | The container uses `securityContext.privileged` with the value `true`. |
| `nginx`         | immutable   | The container uses `securityContext.readOnlyRootFilesystem` with the value `true`. Directories required by nginx have been mounted as Volumes. |

## Deleting Mutable Pods

Delete the Pods that define mutable containers.

```
$ kubectl delete pod busybox -n h92
pod "busybox" deleted
$ kubectl delete pod hello-world -n h92
pod "hello-world" deleted
$ kubectl delete pod loop -n h92
pod "loop" deleted
```

You should be left with a single Pod.

```
$ kubectl get pods -n h92
NAME    READY   STATUS    RESTARTS   AGE
nginx   1/1     Running   0          41m
```