# [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)

- Ingress manages external access to the services in a cluster, typically HTTP.
- Ingress may provide load balancing, SSL termination and name-based virtual hosting.



### 1. Create the following
 - Deployment `web` with image `gcr.io/google-samples/hello-app:1.0` with 3 replicas. 
 - Service `web` to expose the deployment as Node Port
 - Ingress `web-ingress` to point to the `web` service using host `hellow-world.info`.


```bash
kubectl create deployment web --image=gcr.io/google-samples/hello-app:1.0
kubectl expose deployment web --type=NodePort --port=8080
kubectl get service web
#NAME   TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
#web    NodePort   172.17.23.11   <none>        8080:31646/TCP   1s
```

### 2. Create Ingress with the below specs and apply using `kubectl apply -f web-ingress.yaml`

```yaml
cat << EOF > web-ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: web-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
    - host: hello-world.info
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: web
                port:
                  number: 8080
EOF

kubectl apply -f web-ingress.yaml
```


```bash
# verification
kubectl get nodes -o wide # get node ip
kubectl get deploy web # check status
kubectl get svc web # check node port ip

local machine terminal
port forward terminal 
kubectl port-forward --namespace=ingress-nginx service/ingress-nginx-controller 8080:80

second terminal
curl --resolve hello-world.info:8080:127.0.0.1 http://hello-world.info:8080
# Hello, world!
# Version: 1.0.0
# Hostname: web-548f6458b5-xwbld

kubectl get ingress web-ingress # you will get an ip address of the ingress controller if installed
# NAME          CLASS    HOSTS              ADDRESS   PORTS   AGE
# web-ingress   <none>   hello-world.info             80      11s
```

## Ingress Security


### 3. Create a tls secret `testsecret-tls` using tls.crt from file `../data/tls.crt` and `../data/tls.key`. Enable tls for the ingress below.

```bash
$ openssl req -nodes -new -x509 -keyout tls.key -out tls.crt -subj "/CN=hello-world.tls"
```
```bash
$ kubectl create secret tls hello-world-secret --cert=tls.crt --key=tls.key --dry-run=client -o yaml
```

```yaml
apiVersion: v1
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURGVENDQWYyZ0F3SUJBZ0lVVytLYVhaTUJQT3YvMTJmTG5CQ1lqU2Ixcmljd0RRWUpLb1pJaHZjTkFRRUwKQlFBd0dqRVlNQllHQTFVRUF3d1BhR1ZzYkc4dGQyOXliR1F1ZEd4ek1CNFhEVEl6TURjeE5UQTFNelEwT0ZvWApEVEl6TURneE5EQTFNelEwT0Zvd0dqRVlNQllHQTFVRUF3d1BhR1ZzYkc4dGQyOXliR1F1ZEd4ek1JSUJJakFOCkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQTF5WGJGdThEZm1CUTdCUC9MWnpHQ08rMURZQWoKZmxxYXVzY0NkY2ZLb0g2OVAxZnJZWThBdlVqQlNhZ1hGckovSjhMSVJlVXR0TndzMGNvLzF2VVQ3Tks1VExVaAo5WXM0MnROMy9DWmQwWEpNYWF6aDhoM0xqcnpkRFlrWHIzUFF6QTRGb2lnOU5JYUdscld6ODNqbkthU0twMUNXClI2QjJEMzZlU1p0Rk9zZ0VjUjFuZVVncEsvNXRYUjVoMXkxR21TTWF2TGcwaGhlR09pYVdWV1JjZG0vV2NsYXcKMkNsVkZEY3c0UGwxSUI2RCtEY0ZyMzJSemxiVHg3eG9XZmxWRC8xRlYvb0V5U1BiL1F6ZHJCb0Y4M3RwNXJsbwoxZ1lUYnRWdVNsemZIbmNZL09MWmtBTFk3K1J1UXJMZjhvNWJVWnZLOTZkbFc4TUhycUJnaDFpV3lRSURBUUFCCm8xTXdVVEFkQmdOVkhRNEVGZ1FVNCtURTBDWlVvOE5lL2t1TExuWGMzcSsxQVpnd0h3WURWUjBqQkJnd0ZvQVUKNCtURTBDWlVvOE5lL2t1TExuWGMzcSsxQVpnd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBTkJna3Foa2lHOXcwQgpBUXNGQUFPQ0FRRUFPdXJLOGxsQm5OaU9mcExnTlZaYkxPN085MzdRZEhzVmg2SEkwQTk2Q2VhKzZ2UlJ2V3lnClIwY0p3RDJncVg0M0NZVzlIdmdwa0Q0eWx1amhNdnZ2V1JtK01pTU5kQkk0dzRpTHB1TzJDQWJiUUxGY2Y0cXAKSFJreHlvWkRCeEhLSlFabHlxTFhWeGhNdUJxcUxad2trT0d5VElSWHhLcHAwd0JlMUJ6WUVtQlY3NlBzQjNRdgpmV2VRZjExay84TmZWNk9rVFFjWkF5RWp6bGJXT2dTL2VHOHlKd1BPaDZrUURYM2FibEVoRFNQNlkxZHZQWnRNCmovQ2UwY1puWkNBbUI0UWdwNi96dm1ZV3R5anlrN2hLVFRtdm1veGxYQ3V0a3Qyc2RvYk8vTjI3K3ZPWkduU2cKbUVtYWY4bUhrQXJpNG1GeVpDM1ZGRzAxeDZ5Q3B1S2plZz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
  tls.key: LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRRFhKZHNXN3dOK1lGRHMKRS84dG5NWUk3N1VOZ0NOK1dwcTZ4d0oxeDhxZ2ZyMC9WK3RoandDOVNNRkpxQmNXc244bndzaEY1UzIwM0N6Ugp5ai9XOVJQczBybE10U0gxaXpqYTAzZjhKbDNSY2t4cHJPSHlIY3VPdk4wTmlSZXZjOURNRGdXaUtEMDBob2FXCnRiUHplT2NwcElxblVKWkhvSFlQZnA1Sm0wVTZ5QVJ4SFdkNVNDa3IvbTFkSG1IWExVYVpJeHE4dURTR0Y0WTYKSnBaVlpGeDJiOVp5VnJEWUtWVVVOekRnK1hVZ0hvUDROd1d2ZlpIT1Z0UEh2R2haK1ZVUC9VVlgrZ1RKSTl2OQpETjJzR2dYemUybm11V2pXQmhOdTFXNUtYTjhlZHhqODR0bVFBdGp2NUc1Q3N0L3lqbHRSbThyM3AyVmJ3d2V1Cm9HQ0hXSmJKQWdNQkFBRUNnZ0VBRUZyUGRpVjJhcWFMOUR3N2FVWlorLzhyMEZKSnNHN1hCUVNBakg5QzZSTWwKVS9HOUlnSUd4SHpKdTYySVRzaUVGN3ZMOFZzLzc1Z09uS1ovRnBveitPeFdXMjFIL3ZSRnZzVzFRT0RTUUN3TQpBTHI4aElVTkRxVk8rUHJQUnY0NjlLNEtzZFpVKzcrUHJ0aHliTnRIbUQvcXJkZ1hpbXVaV2VsdVF5cG5WM1hPCkNCOG1TUmdGUW94aC93TDBScVpUNExNMnFRdWlDVmoxUUlUUDA5QUR3V2QyN2NzNGg3Ylhqb2M1UXdsY0pxcmMKbll2M3BqUmZIUzlyMHI1L3lIOWlnekh4Rzd4cDA4UlE2blZsNFVmNDVSOWQyZ0dpMTNLZHhEbTg1cmhzMXMwVwo4Nk9DWnNuV0VqYTZMS3g1R25ZWVo0ZFA1MUhETUpKeVZHU00xK1ZOc1FLQmdRRHVySFVDTTM3a1gyVmFxRjdGCkRBdTFJMEpsR2JHVm5EYXdUSjlHZkFwbFN3azBlczc0ZmJrY2s4ZkxROFU0K3FZSWJTdmtubU5iL3RrRnZEU3AKdThZWFpyUkZSTkI2TDUza1dEcG5lYmxIMUlEd0NxbEpSL3NHeHZkT29NZ1IxaVM2WlptTHI3algyWFNXRVlRRQp1bGlBQlB0ZVBwS0Q1MDdoeWlydXRaNmlXUUtCZ1FEbXhERkM0NXdlM2IyamQzUWhhb2Fwc2EvTUJZRmtKZURoCnhnanhFdERadjZiZnhYV2FTR1UvcUZEVk45YWVhSXE2dlhmcXJBbVFYcFpwZXVNVGJMMFl6T1JMRlZBZHpKejQKL0N0STVhSTREK3Z4QXZ0M3BWbFlQWDl6WjZsZHZrN3RMUmVaU0hURHpPN0poUjhSQmoxbmk1MTNrSzFQcVVqZAoxbnpsZGQycDhRS0JnR25vV3MrTVBEMW9XMWc4S0RsMTIrZ1g2a2xrZWxteUJNNThZWnpzUTh5bzNEL09Vdk1QCkRzb3doVktjRXZMeXZpUTFGV29RSW5UYkJFQTRRSUlrMFlUbXpRMkR1S0xQYUZmQUVzS0IwQWFndEVwUDRNaWsKeVM0N3NkZlJIcFpUbU42aFlONE1McENSYk50K2tEOXMxUUtSczEwaGxQQTAvdHJRVS9heTN0NlpBb0dCQU5mSwpkdFk1eVhBRG9xWU51Z2JWZW4rTXJQdDMybnN3QUk2ZmhSVUpwMHMzc0hCb1VZU29WaVdrMjVIRzJXYmdFM3AwCldvS1R4WTIvSmFBVlRlcUlNSHZVZlZiSTQxWTZWRDl3YjNtSFlwNVNKU3dHd2Vnc010UVZEZTE0N2lPL0prY2IKZDRuYSszVnRXWTRRY2ZDVmlTNGVuRzJHV01xNVhtNHplQytLZkhIQkFvR0FZREJpKzBPbXNlNFRSbk9wUTNiSApQYjBEaTlQMGxBaHRQVERIVFRmV3UrS3hiSjhJM1o0NU5mVTdhS0ZGVi8xWEE2SDZaelRITVd6WkEyUG1SZnBWClJhSGtwc0o5UDZTdUdJWXZuQlo5ZHI1ZXhEZXVSU3hyUENnZW1LR2N0c0ZITTA5Y0NBUlpVY2lxbFBiSktKL0QKbDVibG5yZnpCZzlGSlJXeHJJL2syeTQ9Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K
kind: Secret
metadata:
  creationTimestamp: "2023-07-15T05:36:01Z"
  name: hello-world-secret
  namespace: default
  resourceVersion: "70018"
  uid: 43ca4f58-f9c5-41f8-8fd9-164d2d09fc70
type: kubernetes.io/tls
```


<br />

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: tls-hello-world-ingress
spec:
  rules:
  - host: hello-world.info
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: service1
            port:
              number: 80
```

```bash
$ kubectl create secret tls testsecret-tls --cert=tls.crt --key=tls.key
```

```bash
$ kubectl create ingress hello-world-ingress --rule="hello-world.info/*=web-ingress:80,tls=hello-world-secret"
```
```bash

cat << EOF > hello-world-ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  creationTimestamp: null
  name: hello-world-ingress
spec:
  rules:
  - host: hello-world.info
    http:
      paths:
      - backend:
          service:
            name: web-ingress
            port:
              number: 80
        path: /
        pathType: Prefix
  tls:
  - hosts:
    - hello-world.info
    secretName: hello-world-secret
status:
  loadBalancer: {}
EOF
```
```bash
$ kubectl apply -f hello-world-ingress.yaml

```

```bash
# verification
$ kubectl get secret hello-world-secret
$ kubectl get ingress tls-example-ingress
```
```
NAME                  CLASS    HOSTS              ADDRESS   PORTS     AGE
hello-world-ingress   <none>   hello-world.info             80, 443   6m28s
```

## 4. Clean up

<br />

```bash
$ kubectl delete secret hello-world-secret
$ kubectl delete ingress hello-world-ingress web-ingress
$ kubectl delete svc web
$ kubectl delete deployment web
```