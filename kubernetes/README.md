
### Kubernetes 

#### Locally
using kind to deploy image locally
```bash
docker build --pull --rm -f "Dockerfile" -t gym:1.0.0 "."
kind load docker-image gym-pod:1.0.0
```


#### Using gitlab registry as the app img source
Before you can pull from the private repository, in order for the pulling to go through, you need to create a secret for Kubernetes.
According to the Kubernetes documentation 22, you can create a new secret by executing the following with your username (below, k8s) and token:

```bash
kubectl create secret docker-registry regcred --docker-server=registry.gitlab.com --docker-username=k8s --docker-password=<token>
```