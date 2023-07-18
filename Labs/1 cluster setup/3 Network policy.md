# [Network policy](https://kubernetes.io/docs/concepts/services-networking/network-policies/)

Network policy is a foundational component of K8S security. For the exam you will need to know how to crete specific policies based on provided requirements and be able to validate that your policy is working as expected.

Great resource for various example network policies
https://github.com/ahmetb/kubernetes-network-policy-recipes



### 1. Create the namespace and nginx service
Create a new namespace called netpolicy and deploy an nginx image to it on port 80

```
kubectl create ns netpolicy
kubectl create deployment --namespace=netpolicy nginx --image=nginx
kubectl expose --namespace=netpolicy deployment nginx --port=80
```

### 2. Verify access - allowed all ingress and egress
Open up a second shell and create a busybox pod to test policy access from the netpolicy namespace. Once run you will be at a shell prompt in the new container.

```bash
$ kubectl run --namespace=netpolicy access --rm -ti --image busybox /bin/sh

If you don't see a command prompt, try pressing enter.
/ #
```

From within the busybox "access" pod execute the following command to test access to the nginx service.
``````bash
wget -q --timeout=5 nginx -O -
``````

It should return the HTML of the nginx welcome page.

Still within the busybox "access" pod, issue the following command to test access to google.com.
``````bash
wget -q --timeout=5 google.com -O -
``````
It should return the HTML of google.com home page.

### 3. Deny all ingress traffic
Enable ingress isolation on the namespace by deploying a default deny all ingress traffic policy.
``````bash
kubectl create -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-ingress
  namespace: netpolicy
spec:
  podSelector:
    matchLabels: {}
  policyTypes:
  - Ingress
EOF
``````
Verify access - denied all ingress and allowed all egress
Because all pods in the namespace are now selected, any ingress traffic which is not explicitly allowed by a policy will be denied.

We can see that this is the case by switching over to our "access" pod in the namespace and attempting to access the nginx service.
``````bash
wget -q --timeout=5 nginx -O -
``````
It should return:

wget: download timed out

Next, try to access google.com.

``````bash
wget -q --timeout=5 google.com -O -
``````
It should return:
``````
<!doctype html><html itemscope="" item....
``````
We can see that the ingress access to the nginx service is denied while egress access to outbound internet is still allowed.

### 4. Allow ingress traffic to Nginx
Run the following to create a NetworkPolicy which allows traffic to nginx pods from any pods in the netpolicy namespace.
``````bash
kubectl create -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: access-nginx
  namespace: netpolicy
spec:
  podSelector:
    matchLabels:
      app: nginx
  ingress:
    - from:
      - podSelector:
          matchLabels: {}
EOF
``````
### 5. Verify access - allowed nginx ingress
Now ingress traffic to nginx will be allowed. We can see that this is the case by switching over to our "access" pod in the namespace and attempting to access the nginx service.
``````bash
wget -q --timeout=5 nginx -O -
``````
It should return:
``````
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>...
``````
After creating the policy, we can now access the nginx Service.

### 6. Deny all egress traffic
Enable egress isolation on the namespace by deploying a default deny all egress traffic policy.
``````bash
kubectl create -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-egress
  namespace: netpolicy
spec:
  podSelector:
    matchLabels: {}
  policyTypes:
  - Egress
EOF
``````
### 7. Verify access - denied all egress
Now any ingress or egress traffic which is not explicitly allowed by a policy will be denied.

We can see that this is the case by switching over to our "access" pod in the namespace and attempting to nslookup nginx or wget google.com.
``````bash
nslookup nginx
``````
It should return something like the following.

Server:    10.96.0.10
Address 1: 10.96.0.10
;; connection timed out; no servers could be reached
nslookup: can't resolve 'nginx'

Next, try to access google.com.
``````bash
wget -q --timeout=5 google.com -O -
``````
It should return:

wget: bad address 'google.com'

NOTE
The nslookup command can take a minute or more to timeout.

### 7. Allow DNS egress traffic
Run the following to create a label of name: kube-system on the kube-system namespace and a NetworkPolicy which allows DNS egress traffic from any pods in the netpolicy namespace to the kube-system namespace.
``````bash
$ kubectl label namespace kube-system name=kube-system
$ kubectl create -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-dns-access
  namespace: netpolicy
spec:
  podSelector:
    matchLabels: {}
  policyTypes:
  - Egress
  egress:
  - to:
    - namespaceSelector:
        matchLabels:
          name: kube-system
    ports:
    - protocol: UDP
      port: 53
EOF
``````
### 8. Verify access - allowed DNS access
Now egress traffic to DNS will be allowed.

We can see that this is the case by switching over to our "access" pod in the namespace and attempting to lookup nginx and google.com.
```bash
nslookup nginx
```
It should return something like the following.

Server:    10.0.0.10
Address 1: 10.0.0.10 kube-dns.kube-system.svc.cluster.local

Next, try to look up google.com.
```bash
nslookup google.com
```
It should return something like the following.

Name:      google.com
Address 1: 2607:f8b0:4005:807::200e sfo07s16-in-x0e.1e100.net
Address 2: 216.58.195.78 sfo07s16-in-f14.1e100.net

Even though DNS egress traffic is now working, all other egress traffic from all pods in the netpolicy namespace is still blocked. Therefore the HTTP egress traffic from the wget calls will still fail.

### 9. Allow egress traffic to nginx
Run the following to create a NetworkPolicy which allows egress traffic from any pods in the netpolicy namespace to pods with labels matching app: nginx in the same namespace.
``````bash
$ kubectl create -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-egress-to-advance-policy-ns
  namespace: netpolicy
spec:
  podSelector:
    matchLabels: {}
  policyTypes:
  - Egress
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: nginx
EOF
``````
### 10. Verify access - allowed egress access to nginx
We can see that this is the case by switching over to our "access" pod in the namespace and attempting to access nginx.
```bash
wget -q --timeout=5 nginx -O -
```
It should return the HTML of the nginx welcome page.
``````
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>...
``````
Next, try to retrieve the home page of google.com.
```bash
wget -q --timeout=5 google.com -O -
```
It should return:
```
wget: download timed out

Access to google.com times out because it can resolve DNS but has no egress access to anything other than pods with labels matching app: nginx in the netpolicy namespace.
```
### 11.  Clean up namespace
``````bash
$ kubectl delete ns netpolicy
``````