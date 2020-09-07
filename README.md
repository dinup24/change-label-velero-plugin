## Velero Plugin to change label of Kubernetes resources during restore action

### Build the plugin locally
_This requires golang to be installed locally_
```
go build .
```

### Build and create a docker image
_This doesn't require golang to be installed locally_  
_This requires docker to be installed locally_
```
docker build -t change-label-velero-plugin .
```

### Setting up the plugin
1. Add the plugin
```
velero plugin add <image>
```
eg. `velero plugin add change-label-velero-plugin:latest`
Ensure that velero pod has restarted successsfully `kubectl get pods -n velero`

2. Check whether plugin was installed successfully
```
velero plugin get
```
The plugin `velero.io/change-label` should be listed.

### Configuring the labels to be set using a configmap
Create a configmap
```
apiVersion: v1
kind: ConfigMap
metadata:
  name: change-label-config
  namespace: velero
  labels:
    velero.io/plugin-config: ""
    velero.io/change-label: RestoreItemAction
data:
  zone: dal13
  region: us-south
```

### Removing up the plugin
1. Remove the plugin
```
velero plugin remove <image>
```
eg. `velero plugin remove change-label-velero-plugin:latest`
Ensure that velero pod has restarted successsfully `kubectl get pods -n velero`

2. Check whether plugin was removed successfully
```
velero plugin get
```
The plugin `velero.io/change-label` should not be listed.

