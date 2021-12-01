## Custom Kubernetes Scheduler example

### Build Scheduler

```bash
# for local os/arch, please run
make

# for linux/amd64, please run
GOOS=linux GOARCH=amd64 make
```

The binary will be put in `bin` folder.

### Run our custom scheduler

#### Disable system scheduler

Prepare a Kubernetes cluster with one node, then disable the system kube-scheduler with the following command:

```bash
mv /etc/kubernetes/manifests/kube-scheduler.yaml /etc/kubernetes/
```

#### Run the custom scheduler with the example config

To simplify permission configuration, admin kubeconfig is used directly here:

```bash
bin/scheduler --config example-config.yaml --authentication-kubeconfig=/etc/kubernetes/admin.conf --authorization-kubeconfig=/etc/kubernetes/admin.conf
```

### Verify that Pod will not be scheduled when node load is high

Apply nginx deploy with memory request（let memory usage exceeds the node memory usage limit）

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment-big-memory
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: 18Gi
```

Then apply more nginx Pod, these Pods will be Pending state.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment-pending
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: 1Gi
```

Scheduler logs like:

```bash
I1201 12:52:13.405370   10119 NodeMemoryUsageLimit.go:82] Current Node memory usage 60.45660111469591, limit 20
I1201 12:52:13.405540   10119 NodeMemoryUsageLimit.go:82] Current Node memory usage 60.45660111469591, limit 20
```

Pod event like:

```bash
Events:
  Type     Reason            Age   From               Message
  ----     ------            ----  ----               -------
  Warning  FailedScheduling  11m   default-scheduler  0/1 nodes are available: 1 node memory usage reach the limit 20.
  Warning  FailedScheduling  10m   default-scheduler  0/1 nodes are available: 1 node memory usage reach the limit 20.
```
