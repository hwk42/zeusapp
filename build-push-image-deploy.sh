IMG=10.252.39.13:5000/aiscene/zeusapp-controller:v1
docker build --build-arg https_proxy=10.1.180.122:20171 -t ${IMG} .
docker push ${IMG}
kubectl -nzeusapp-system delete deploy zeusapp-controller-manager
kubectl apply -f config/default