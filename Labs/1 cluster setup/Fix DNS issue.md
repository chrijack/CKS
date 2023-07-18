kubectl get pods -n kube-system -oname |grep coredns |xargs kubectl delete -n kube-system

sudo nano /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
add below environmet Environment="KUBELET_EXTRA_ARGS=--resolv-conf=/run/systemd/resolve/resolv.conf"

systemctl daemon-reload && systemctl restart kubelet

sudo systemctl daemon-reload

kubectl get pods -n kube-system -oname |grep coredns |xargs kubectl delete -n kube-system

Environment="KUBELET_EXTRA_ARGS=--resolv-conf=/run/systemd/resolve/resolv.conf"

kubectl create deployment multitool --image=wbitt/network-multitool --name=multitool

kubectl run access -ti -rm --image=wbitt/network-multitool /bin/sh
kubectl exec -i -t multitool -- /bin/sh
Server:		172.17.0.10
Address:	172.17.0.10#53

Name:	nginx.netpolicy.svc.cluster.local
Address: 172.17.9.90


https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/kubelet-integration/#the-kubelet-drop-in-file-for-systemd



minikube start --nodes 2 --memory 4096 --cpus 2 --driver=docker --kubernetes-version v1.27.3



Install ingress via helm
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace


  https://kubernetes.github.io/ingress-nginx/deploy/#bare-metal-clusters


Pre-flight check¶
A few pods should start in the ingress-nginx namespace:


kubectl get pods --namespace=ingress-nginx
After a while, they should all be running. The following command will wait for the ingress controller pod to be up, running, and ready:


kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s
Local testing¶
Let's create a simple web server and the associated service:


kubectl create deployment demo --image=httpd --port=80
kubectl expose deployment demo
Then create an ingress resource. The following example uses a host that maps to localhost:


kubectl create ingress demo-localhost --class=nginx \
  --rule="demo.localdev.me/*=demo:80"
Now, forward a local port to the ingress controller:


kubectl port-forward --namespace=ingress-nginx service/ingress-nginx-controller 808

curl --resolve demo.localdev.me:8080:127.0.0.1 http://demo.localdev.me:8080

k delete demo-localhost
k delete ingress demo-localhost
k delete svc demo
k delete deployment demo