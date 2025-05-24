## Simple Kubernetes Provider

This provider starts a privileged Docker-in-Docker pod in a Kubernetes cluster and connects the workspace to it. To use this provider simply clone the repo and run:
```
devspace provider add ./examples/provider.yaml
```

Then start a new workspace via:
```
devspace up github.com/microsoft/vscode-course-sample
```
