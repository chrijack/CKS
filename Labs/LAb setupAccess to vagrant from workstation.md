# CKS Crash Course Lab Setup

This lab environment is intended to give you a working kubernetes cluster on your local machine that can be used for study and experimentation. It uses Vagrant to get you up and running quickly. It has been tested on Intel mac, Linux, and Windows 11. Apple silicon will work as well, but I have not been able to throughly test this as I don't have one yet.


To get started, the First thing you will need to do is pull down the needed open source software.

Latest Vagrant
https://developer.hashicorp.com/vagrant/downloads

Latest virtual box
https://www.virtualbox.org/wiki/Downloads

## Mac

Use Home brew to install

 ```bash
$ brew install hashicorp/tap/hashicorp-vagrant
$ brew install kubectl
```

## Windows
install python
https://www.python.org/ftp/python/3.11.4/python-3.11.4-amd64.exe

install pywin32 python libraries from PS prompt
```powershell
py -m pip install pywin32
```
install virtual box
https://download.virtualbox.org/virtualbox/7.0.8/VirtualBox-7.0.8-156879-Win.exe

install vagrant
https://releases.hashicorp.com/vagrant/2.3.7/vagrant_2.3.7_windows_amd64.msi

Install git
https://github.com/git-for-windows/git/releases/download/v2.41.0.windows.1/Git-2.41.0-64-bit.exe

Download Kubectl binary
```powershell
 curl.exe -LO "https://dl.k8s.io/release/v1.27.3/bin/windows/amd64/kubectl.exe"
 ```
or install via windows package manager [(Chocolatey, Scoop, or winget
)](https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/#install-nonstandard-package-tools)

You should have already cloned the github repository for class. Launch Vagrant.

1. clone class github repo

2. cd to vagrant files

3. vagrant up

## Connect workstation kubectl to vagrant cluster

Once Vagrant execution is successful, you will see a configs folder with a few files (config, join.sh, and token) inside the cloned repo. These are generated during the run time.

Copy the config file to your $HOME/.kube folder if you want to interact with the cluster from your workstation terminal. You should have kubectl installed on your workstation. This saves the current config of Vagrant, but will need to be replaced if you reprovision the cluster.

```bash
mkdir -p $HOME/.kube
cp configs/config $HOME/.kube
```

I prefer to use a Kubeconfig env variable as shown below. Make sure you execute the commands below from the CKS lab main directory.

```bash
export KUBECONFIG=$(PWD)/configs/config
```
Power shell
```powershell
$Env:KUBECONFIG="$PWD\configs\config"
```

Verify the config by listing the cluster nodes.
```bash
$ kubectl get nodes
```

This Vagrant install script was updated and modified from the below post and github link 

## [Devopscube Install vagrant and virtualbox tutorial](https://devopscube.com/kubernetes-cluster-vagrant/)

## [Github link](https://github.com/techiescamp/vagrant-kubeadm-kubernetes)