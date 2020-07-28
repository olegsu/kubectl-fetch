# kubectl-fetch
kubectl plugin to download Kubernetes manifests files and apply onto the cluster
`kubectl fetch git --repo https://github.com/kubernetes/examples --path guestbook/all-in-one | kubectl apply -f -`