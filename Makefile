all: scheduler

# Build scheduler binary
scheduler: generate fmt
	go build -o bin/scheduler cmd/scheduler/scheduler.go

# Generate code
generate: controller-gen conversion-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."
	$(CONVERSION_GEN) -h hack/boilerplate.go.txt \
	-O zz_generated.conversion \
	-i github.com/SataQiu/k8s-scheduler-example/api/v1beta2 \
	-p github.com/SataQiu/k8s-scheduler-example/api/v1beta2/

fmt:
	go fmt ./...
# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

# find or download conversion-gen
# download conversion-gen if necessary
conversion-gen:
ifeq (, $(shell which conversion-gen))
	@{ \
	set -e ;\
	CONVERSION_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONVERSION_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get k8s.io/code-generator/cmd/conversion-gen@v0.22.4 ;\
	rm -rf $$CONVERSION_GEN_TMP_DIR ;\
	}
CONVERSION_GEN=$(GOBIN)/conversion-gen
else
CONVERSION_GEN=$(shell which conversion-gen)
endif