# Commander

Commander is the provisioning component of the Astronomer Platform. It is responsible for interacting with the underlying infrastructure layer. It takes care of deployments on different scheduling systems like Kubernetes and Marathon.

## Configuration

* `DEBUG_MODE`: Logs at DEBUG level.
* `PORT`: Port for service to listen on.
* `KUBECONFIG`: Path to a kubectl config. Typically ~/.kube/config in development. Left blank in production and assumes service role of node.
* `KUBE_NAMESPACE`: Kubernetes namespace to operate within.
* `HELM_DEBUG`: true/false to enable/disable helm debugging
* `HELM_HOME`:
* `HELM_HOST`:
* `KUBECONFIG`:
* `TILLER_NAMESPACE`:

## gRPC functions

* CreateDeployment
* ~~FetchDeployment~~
* DeleteDeployment
* UpdateDeployment
* PatchDeployment

## Development

### Install protobuf compiler

Visit the protobuf [release page][1], at the bottom of the list there are `protoc-*` zips, download the one for your OS.

[1]: https://github.com/google/protobuf/releases/latest
