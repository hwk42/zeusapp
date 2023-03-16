IMG=10.252.39.13:5000/aiscene/zeusapp-controller:v1
make docker-build docker-push IMG=${IMG}
make update IMG=${IMG}