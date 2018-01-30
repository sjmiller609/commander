# Commander

Commander is the provisioning component of the Astronomer Platform. It is responsible for interacting with the underlying infrastructure layer. It takes care of deployments on different scheduling systems like Kubernetes and Marathon.

## Configuration
- `DEBUG_MODE`: Logs at DEBUG level.
- `PORT`: Port for service to listen on.
- `KUBE_CONFIG`: Path to a kubectl config. Typically ~/.kube/config in development. Left blank in production and assumes service role of node.
- `KUBE_NAMESPACE`: Kubernetes namespace to operate within.
