CLI tool to list K8ts ConfigMaps that particular user has access to
===================================================================

How to install?
---------------

First of all you need following tools:

 - [glide](https://github.com/Masterminds/glide)

The just use make:

```bash
    make
```

Why Glide?
----------

Kubernetes client does not support deb tool

How to use?
-----------

At the first run try use `-h` for help:
```bash
./kube-configmaps -h
Usage of ./kube-configmaps:
  -alsologtostderr
    	log to standard error as well as files
  -group string
    	k8ts user group
  -kubeconfig string
    	(optional) absolute path to the kubeconfig file (default "/Users/denismakogon/.kube/config")
  -log_backtrace_at value
    	when logging hits line file:N, emit a stack trace
  -log_dir string
    	If non-empty, write log files in this directory
  -logtostderr
    	log to standard error instead of files
  -namespace string
    	k8ts namespace (default "default")
  -stderrthreshold value
    	logs at or above this threshold go to stderr
  -user string
    	k8ts user
  -v value
    	log level for V logs
  -vmodule value
    	comma-separated list of pattern=N settings for file-filtered logging
```

All you'd need from this opts is:

 - `-kubeconfig` or `--kubeconfig`
 - `-namespace` or `--namespace`
 - `-user` or `--user`
 - `-group` or `--group`


Sample templates to use
-----------------------

Use [templates/configmap.yml](templates/configmap.yml) to create configmaps using the following command:
```bash
kubectl create -f templates/configmap.yml 
configmap "configmap-1" created
configmap "config" created
configmap "nginx" created
configmap "caddy" created
configmap "nginx-different" created
```

Use [templates/roles.yml](templates/roles.yml) to create few roles and their bindings using the following command:
```bash
kubectl create -f templates/roles.yml 
role "test-role" created
rolebinding "ann-binding" created
rolebinding "denis-binding" created
role "test-role-2" created
role "test-role-3" created
rolebinding "lary-bidning" created
rolebinding "dave-binding" created
```

After those resources got created you can query using different users and groups to see configmaps available for the access:
```bash
 ./kube-configmaps -user=denis
User 'denis' granted with access according to its role 'test-role' to the following configmaps: [caddy config configmap-1 nginx nginx-different]
```

```bash
./kube-configmaps -group=system:unauthenticated
Group 'system:unauthenticated' granted with access according to its role 'test-role-2' to the following configmaps: [configmap-1 config nginx caddy]
```

Edge cases
----------

You can't set both group and user at the same time.
You can't left unset both group and user at the same time.

I definitely need to learn more about RBAC for K8ts =/
