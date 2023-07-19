# mTLS example






apiVersion: apps/v1
kind: Deployment
metadata:
  name: meow
spec:
  replicas: 2
  selector:
    matchLabels:
      app: meow
  template:
    metadata:
      labels:
        app: meow
    spec:
      containers:
      - name: meow
        image: gcr.io/kubernetes-e2e-test-images/echoserver:2.1
        ports:
        - containerPort: 8080



apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/auth-tls-verify-client: \"on\"
    nginx.ingress.kubernetes.io/auth-tls-secret: \"default/my-certs\"
  name: meow-ingress
  namespace: default
spec:
  ingressClassName: nginx
  rules:
  - host: meow.com
    http:
      paths:
      - backend:
          service:
            name: meow-svc
            port:
              number: 80
        path: /
        pathType: Prefix
  tls:
  - hosts:
    - meow.com
    secretName: my-certs


    # Generate the CA Key and Certificate
$ openssl req -x509 -sha256 -newkey rsa:4096 -keyout ca.key -out ca.crt -days 356 -nodes -subj '/CN=Fern Cert Authority'
# Generate the Server Key, and Certificate and Sign with the CA Certificate
$ openssl req -new -newkey rsa:4096 -keyout server.key -out server.csr -nodes -subj '/CN=meow.com'
$ openssl x509 -req -sha256 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt
# Generate the Client Key, and Certificate and Sign with the CA Certificate
$ openssl req -new -newkey rsa:4096 -keyout client.key -out client.csr -nodes -subj '/CN=Fern'
$ openssl x509 -req -sha256 -days 365 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 02 -out client.crt
Creating the Kubernetes Secrets
We must store the certificates generated above in a Kubernetes Secret in order to use them in our Ingress-NGINX controller.

In this example, for simplicity, our Secret will contain both our Server Certificate and our CA Certificate. The Ingress Controller will understand which certs to use and where to use them. They can also be Split into Separate Secrets. See here for detailed info.

$ kubectl create secret generic my-certs --from-file=tls.crt=server.crt --from-file=tls.key=server.key --from-file=ca.crt=ca.crt
$ kubectl get secret my-certs
NAME       TYPE     DATA   AGE
my-certs   Opaque   3      1m
Deploying an Application with Mutual Authentication
Deploy our pods(containers), using deployments. A deployment, pretty much just manages the state of our pods, It’s out of the scope of this tutorial, but you can learn more here
$ echo "
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: meow
spec:
  replicas: 2
  selector:
    matchLabels:
      app: meow
  template:
    metadata:
      labels:
        app: meow
    spec:
      containers:
      - name: meow
        image: gcr.io/kubernetes-e2e-test-images/echoserver:2.1
        ports:
        - containerPort: 8080
" | kubectl apply -f -
# wait a min for the deployment to be created
$ kubectl get deploy
NAME       DESIRED   CURRENT   UP-TO-DATE   AVAILABLE
meow       2         2         2            2
# you should have 2 pods running
$ kubectl get pods
NAME                       READY     STATUS
meow-5557bc7c54-cw2ck     1/1       Running
meow-5557bc7c54-kfzm5     1/1       Running
gcr.io/kubernetes-e2e-test-images/echoserver:2.1 just responds with information about the request.

1. Expose our pods using Services. 

$ echo "
apiVersion: v1
kind: Service
metadata:
  name: meow-svc
spec:
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
    name: http
  selector:
    app: meow
" | kubectl apply -f -
# wait a min for the service to be created
$ kubectl get svc
NAME           TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)
meow-svc       ClusterIP   10.107.78.24    <none>        80/TCP
3. Setup the Ingress Rules, and add meow.com to /etc/hosts

$ echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/auth-tls-verify-client: \"on\"
    nginx.ingress.kubernetes.io/auth-tls-secret: \"default/my-certs\"
  name: meow-ingress
  namespace: default
spec:
  rules:
  - host: meow.com
    http:
      paths:
      - backend:
          serviceName: meow-svc
          servicePort: 80
        path: /
  tls:
  - hosts:
    - meow.com
    secretName: my-certs
" | kubectl apply -f -
$ kubectl get ing cat-nginx
NAME           HOSTS     ADDRESS     PORTS     AGE
meow-ingress   meow.com  10.0.2.15   80, 443   1m
# Add meow.com to /etc/hosts
$ sudo -- sh -c 10.0.0.10  meow.com >> /etc/hosts"
This allows us to access the service meow-svc via https://meow.com/.

TLS is enabled and it is using the tls.key and tls.crt provided in the my-certs secret.
The nginx.ingress.kubernetes.io/auth-tls-secret annotation uses ca.crt from the my-certs secret.
Testing it all out
Sending a request without the Client Certificate and Key, should give a 400 Error as follows:

$ curl https://meow.com/ -k
...
<center><h1>400 Bad Request</h1></center>
<center>No required SSL certificate was sent</center>
....
Sending a request with the Client Certificate and Key, should redirect you to the meow-svc:

$ curl https://meow.com/ --cert client.crt --key client.key -k
...
ssl-client-issuer-dn=CN=Fern Cert Authority
ssl-client-subject-dn=CN=Fern
ssl-client-verify=SUCCESS
user-agent=curl/7.54.0
...
If the above did not work for you, skip down to the troubleshooting section.

Other possible settings
nginx.ingress.kubernetes.io/auth-tls-verify-client enables verification of client certificates.
nginx.ingress.kubernetes.io/auth-tls-verify-depth sets the validation depth between the provided client certificate and the certification authority chain.
nginx.ingress.kubernetes.io/auth-tls-error-page sets the URL/Page that user should be redirected in case of a Certificate Authentication Error.
nginx.ingress.kubernetes.io/auth-tls-pass-certificate-to-upstream indicates if the received certificates should be passed or not to the upstream server.
For more information about these setting can be seen in the Ingress-NGINX documentation as well as the NGINX documentation for the directives set by Ingress-NGINX.

Troubleshooting
First thing we can do is check if the NGINX configuration was generated properly with the CA:

$ kubectl get pods -n kube-system | grep nginx-ingress-controller
nginx-ingress-controller-5984b97644-qbwtv   1/1       Running
$ kubectl exec -it -n kube-system nginx-ingress-controller-5984b97644-qbwtv cat /etc/nginx/nginx.conf | grep ssl_client_certificate -A 1
ssl_client_certificate                  /etc/ingress-controller/ssl/default-my-certs.pem;
ssl_verify_client                       on;
Another thing we can look at is logs, and see if theres anything going on with the CA:

$ kubectl get pods -n kube-system | grep nginx-ingress-controller
nginx-ingress-controller-5984b97644-qbwtv   1/1       Running
$ kubectl logs -n kube-system nginx-ingress-controller-5984b97644-qbwtv
....
I0831 21:45:34.090212       5 nginx.go:271] Starting NGINX process
....

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2 
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        resources:
          limits:
            cpu: "0.5"
            memory: "256Mi"
          requests:
            cpu: "0.25"
            memory: "128Mi"
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
  type: NodePort
  ports:
    - port: 80
      targetPort: 80
      nodePort: 32000