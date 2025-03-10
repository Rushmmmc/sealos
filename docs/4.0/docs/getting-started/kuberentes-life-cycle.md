---
sidebar_position: 3
---

# Kubernetes life cycle management

## Stand alone installation kuberentes

```shell
$ sealos run kubernetes:v1.25.0 --single
```

## Cluster installation kubernetes

```shell
$ sealos run kubernetes:v1.25.0 labring/calico:v3.22.1 \
     --masters 192.168.64.2,192.168.64.22,192.168.64.20 \
     --nodes 192.168.64.21,192.168.64.19 -p [your-ssh-passwd]
```

Parameter Description:

| Parameter name | Parameter value example | Parameter description |
| --- | --- | --- |
| --masters | 192.168.0.2 | kubernetes master node address list |
| --nodes | 192.168.0.3 | kubernetes node node address list |
| --ssh-passwd | [your-ssh-passwd] | ssh login password |
|image | kubernetes:v1.25.0 | kubernetes image |

Execute the above command directly on a clean server without any extra effort to start a highly available kubernetes cluster.

## Install various distributed applications

```shell
sealos run labring/helm:v3.8.2 # install helm
sealos run labring/openebs:v1.9.0 # install openebs
sealos run labring/minio-operator:v4.4.16 labring/ingress-nginx:4.1.0 \
   labring/mysql-operator:8.0.23-14.1 labring/redis-operator:3.1.4 # oneliner
```

Such highly available mysql redis, etc. are all available, so you don't need to care about all dependencies.

## Add node

add node :
```shell
$ sealos add --nodes 192.168.64.21,192.168.64.19 
```

add master:
```shell
$ sealos add --masters 192.168.64.21,192.168.64.19 
```

## Delete node

delete node :
```shell
$ sealos delete --nodes 192.168.64.21,192.168.64.19 
```

delete master:
```shell
$ sealos delete --masters 192.168.64.21,192.168.64.19  
```

## clean up the cluster

```shell
$ sealos reset
```