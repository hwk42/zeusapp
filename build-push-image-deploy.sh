IMG=10.252.39.13:5000/aiscene/zeusapp-controller:v1
make docker-build docker-push IMG=${IMG}
LOCALBIN ?= $(shell pwd)/bin
KUSTOMIZE ?= $(LOCALBIN)/kustomize
cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
cd - &&	$(KUSTOMIZE) build config/default | kubectl apply -f -