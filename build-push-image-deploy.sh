IMG=10.252.39.13:5000/aiscene/zeusapp-controller:v1
docker build --build-arg https_proxy=10.1.180.122:20171 -t ${IMG} .
docker push ${IMG}
kubectl -nzeusapp-system get po |grep zeusapp-controller-manager | awk '{print $1}' | xargs kubectl -nzeusapp-system delete po 