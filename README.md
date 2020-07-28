# kubectl-fetch
kubectl plugin to download Kubernetes manifests files and apply onto the cluster

## Examples
* Official example guestbook app  
`kubectl fetch git --repo https://github.com/kubernetes/examples --path guestbook/all-in-one | kubectl apply -f -`

* Use SSH key to get files from private repo  
`kubectl fetch git --repo git@github.com:$OWNER/$NAME.git --key-file $HOME/.ssh/id_rsa --user $USER | kubectl apply -f -`

* Use Github token  
`kubectl fetch git --repo https://github.com/$OWNER/$NAME --token $GITHUB_TOKEN --user $USER | kubectl apply -f -`