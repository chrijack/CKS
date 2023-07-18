# [RBAC authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)

- Role and Role bindings are namespace scoped (pods, deployments, configmaps, etc).
- Cluster Role and Cluster Role bindings are cluster scoped resources and not limited to namespaces (nodes, pv, etc).


### 1. Verify authorization used by the cluster 


Check `/etc/kubernetes/manifests/kube-apiserver.yaml` for `--authorization-mode=Node,RBAC`. This means the cluster will accept Node auth, which is for kubelet, and Role Based Access control for authorization.




## Role and Role Bindings



### 2. Create a role

 Create the role `pods-read` to `get, create, list and delete` `pods` in the default namespace.


```bash
kubectl create role pods-read --verb=get,create,list,delete --resource=pods
```

OR 

```yaml
cat << EOF > pods-read.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: pods-read
rules:
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - get
  - create
  - list
  - delete
EOF

kubectl apply -f pods-read.yaml
```
### 3. Verify role creation
```bash
# verify
kubectl get role pods-read
# NAME        CREATED AT
# pods-read   2021-12-13T01:35:10Z
``` 



### 4. Create a service account `sample-sa`


```bash
kubectl create sa sample-sa
```

OR

```yaml
cat << EOF > sample-sa.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  creationTimestamp: null
  name: sample-sa
EOF

kubectl apply -f sample-sa.yaml
```

```bash
# verify
kubectl get serviceaccount sample-sa
# NAME        SECRETS   AGE
# sample-sa   1         14s
```


### 4. Create role binding 
Create a role binding `sample-sa-pods-read-role-binding` binding service account `sample-sa` and role `pods-read`


```bash
$ kubectl create rolebinding sample-sa-pods-read-role-binding --serviceaccount=default:sample-sa --role=pods-read
```

OR 

```yaml
cat << EOF > sample-sa-pods-read-role-binding.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  creationTimestamp: null
  name: sample-sa-pods-read-role-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: pods-read
subjects:
- kind: ServiceAccount
  name: sample-sa
  namespace: default
EOF

kubectl apply -f sample-sa-pods-read-role-binding.yaml
```

```bash
# verify
kubectl get rolebinding sample-sa-pods-read-role-binding
# NAME                               ROLE             AGE
# sample-sa-pods-read-role-binding   Role/pods-read   18s
```


### 5. Check access for service account

 Verify service account `sample-sa` can get pods using the `auth can-i` command.


```bash
# verify
kubectl auth can-i get pods --as system:serviceaccount:default:sample-sa
# yes
```

## Cluster Role and Cluster Role Bindings


### 6. Create the following for a user `proxy-admin` (which does not exist)
 - Cluster role `proxy-admin-role` with permissions to `nodes` with `get, list,create, update` actions
 - Cluster role binding `proxy-admin-role-binding` to bind cluster role `proxy-admin-role` to user `proxy-admin`


```bash
kubectl create clusterrole proxy-admin-role --resource=nodes --verb=get,list,create,update
kubectl create clusterrolebinding proxy-admin-role-binding --user=proxy-admin --clusterrole=proxy-admin-role
```

```bash
# verify
kubectl auth can-i get nodes --as proxy-admin
# yes
```


### 7. Clean up 


```bash
rm sample-sa-pods-read-role-binding.yaml pods-read.yaml
kubectl delete rolebinding sample-sa-pods-read-role-binding
kubectl delete serviceaccount sample-sa
kubectl delete role pods-read
kubectl delete clusterrolebinding proxy-admin-role-binding
kubectl delete clusterole proxy-admin-role
```