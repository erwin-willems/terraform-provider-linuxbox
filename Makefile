SHELL = /bin/bash
NAMESPACE = erwinwillems
NAME = linuxbox
BINARY = terraform-provider-${NAME}
VERSION = 0.1
terraform_version = 1.4.6
bin ?= ${PWD}/.bin
terraform = ${bin}/terraform
platform = $(shell uname -s|tr '[:upper:]' '[:lower:]')
arch = $(shell uname -m|tr '[:upper:]' '[:lower:]')
OS_ARCH = ${platform}_${arch}

.PHONY: help



help: ## Displays help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-z0-9A-Z_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

tidy: ## Refresh go.mod file
	go mod tidy

clean: ## Clean up environment and build files
	@rm -rf terraform/.terraform || true
	@rm -rf terraform/.terraform.lock.hcl || true
	@rm -rf terraform/terraform.d || true
	@rm -rf terraform/terraform.tfstate || true
	@rm -rf .stage || true

fix: fmt

fmt: ## Format code
	@gofmt -w .

build:
	go build -o ${BINARY}

dev: run

run:
	go run .

stage: \
	.stage/${NAMESPACE}/${NAME}/${VERSION}/linux_amd64/${BINARY} \
	.stage/${NAMESPACE}/${NAME}/${VERSION}/darwin_amd64/${BINARY} \
	.stage/${NAMESPACE}/${NAME}/${VERSION}/windows_amd64/${BINARY}.exe

.stage/${NAMESPACE}/${NAME}/${VERSION}/linux_amd64/${BINARY}:
	@[[ -d .stage/${NAMESPACE}/${NAME} ]] || mkdir -p .stage/${NAMESPACE}/${NAME}
	GOOS=linux GOARCH=amd64 go build -o ./.stage/${NAMESPACE}/${NAME}/${VERSION}/linux_amd64/${BINARY}
	@chmod 0755 .stage/${NAMESPACE}/${NAME}/${VERSION}/linux_amd64/${BINARY}

.stage/${NAMESPACE}/${NAME}/${VERSION}/darwin_amd64/${BINARY}:
	@[[ -d .stage/${NAMESPACE}/${NAME} ]] || mkdir -p .stage/${NAMESPACE}/${NAME}
	GOOS=darwin GOARCH=amd64 go build -o ./.stage/${NAMESPACE}/${NAME}/${VERSION}/darwin_amd64/${BINARY}
	@chmod 0755 .stage/${NAMESPACE}/${NAME}/${VERSION}/darwin_amd64/${BINARY}

.stage/${NAMESPACE}/${NAME}/${VERSION}/windows_amd64/${BINARY}.exe:
	@[[ -d .stage/${NAMESPACE}/${NAME} ]] || mkdir -p .stage/${NAMESPACE}/${NAME}
	GOOS=windows GOARCH=amd64 go build -o ./.stage/${NAMESPACE}/${NAME}/${VERSION}/windows_amd64/${BINARY}.exe

install: build
	mkdir -p terraform/terraform.d/plugins/local/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	cp ${BINARY} terraform/terraform.d/plugins/local/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}


terraform: ${terraform} ## Download and unpack required terraform binary
${terraform}: ${terraform}-${terraform_version}
	ln -sf $< $@
	@touch $@

# download and unpack required terraform binary
TMPDIR ?= /tmp

${terraform}-${terraform_version}:
	cd ${TMPDIR}; wget --no-verbose https://releases.hashicorp.com/terraform/${terraform_version}/terraform_${terraform_version}_${platform}_${arch}.zip
	cd ${TMPDIR}; unzip terraform_${terraform_version}_${platform}_${arch}.zip
	mkdir -p ${@D}
	mv ${TMPDIR}/terraform $@
	@touch $@

mrproper: clean
	@rm -rf ${bin}
	@rm -f test_containers/.built.*

# test targets
test_acc_all: \
	test_setup \
	test_ubuntu_18.04 \
	test_ubuntu_20.04 \
	test_ubuntu_22.04 \
	test_rockylinux_8 \
	test_rockylinux_9

clean_test: test_clean

test_clean:
	@podman stop tf_linux_test_sshd_ubuntu_18.04 > /dev/null 2>&1 || true
	@podman stop tf_linux_test_sshd_ubuntu_20.04 > /dev/null 2>&1 || true
	@podman stop tf_linux_test_sshd_ubuntu_22.04 > /dev/null 2>&1 || true
	@podman stop tf_linux_test_sshd_rockylinux_8 > /dev/null 2>&1 || true
	@podman stop tf_linux_test_sshd_rockylinux_9 > /dev/null 2>&1 || true
	@podman rm tf_linux_test_sshd_ubuntu_18.04 > /dev/null 2>&1 || true
	@podman rm tf_linux_test_sshd_ubuntu_20.04 > /dev/null 2>&1 || true
	@podman rm tf_linux_test_sshd_ubuntu_22.04 > /dev/null 2>&1 || true
	@podman rm tf_linux_test_sshd_rockylinux_8 > /dev/null 2>&1 || true
	@podman rm tf_linux_test_sshd_rockylinux_9 > /dev/null 2>&1 || true
	@rm -f test_containers/.run.*
	@rm -f test_containers/.built.*

test_setup: \
	build_test_containers \
	test_containers/.run.ubuntu_18.04 \
	test_containers/.run.ubuntu_20.04 \
	test_containers/.run.ubuntu_22.04 \
	test_containers/.run.rockylinux_8 \
	test_containers/.run.rockylinux_9
	@echo ""
	@echo "#######################"
	@echo "# Running containers: #"
	@echo "#######################"
	@podman ps | egrep --regexp "tf_linux|CONTAINER ID"

# test_setup-ubuntu:
# 	@podman run --cap-add AUDIT_WRITE --detach --rm --publish 5001:22 --name tf_linux_test_sshd rastasheep/ubuntu-sshd:18.04
# 	@export TF_LINUX_SSH_USER=root
# 	@export TF_LINUX_SSH_HOST=127.0.0.1
# 	@export TF_LINUX_SSH_PORT=5001
# 	@export TF_LINUX_SSH_PASSWORD=root

.PHONY: test_ubuntu_18.04

TF_VARS = \
	TF_LINUX_SSH_USER=root \
	TF_LINUX_SSH_PASSWORD=root \
	TF_LINUX_SSH_HOST=127.0.0.1 \
	TF_ACC=1

test_ubuntu_18.04: test_containers/.run.ubuntu_18.04
	${TF_VARS} TF_LINUX_SSH_PORT=5001; \
	go test ./linuxbox -v -timeout 120s ${test_args}

test_ubuntu_20.04: test_containers/.run.ubuntu_20.04
	${TF_VARS} TF_LINUX_SSH_PORT=5002; \
	go test ./linuxbox -v -timeout 120s ${test_args}

test_ubuntu_22.04: test_containers/.run.ubuntu_22.04
	${TF_VARS} TF_LINUX_SSH_PORT=5003; \
	go test ./linuxbox -v -timeout 120s ${test_args}

test_rockylinux_8: test_containers/.run.rockylinux_8
	${TF_VARS} TF_LINUX_SSH_PORT=5004; \
	go test ./linuxbox -v -timeout 120s ${test_args}

test_rockylinux_9: test_containers/.run.rockylinux_9
	${TF_VARS} TF_LINUX_SSH_PORT=5005; \
	go test ./linuxbox -v -timeout 120s ${test_args}


test_containers/.run.ubuntu_18.04:
	@podman run --cap-add AUDIT_WRITE --detach --rm --publish 5001:22 --name tf_linux_test_sshd_ubuntu_18.04 ubuntu_18.04
	@touch $@

test_containers/.run.ubuntu_20.04:
	@podman run --cap-add AUDIT_WRITE --detach --rm --publish 5002:22 --name tf_linux_test_sshd_ubuntu_20.04 ubuntu_20.04
	@touch $@

test_containers/.run.ubuntu_22.04:
	@podman run --cap-add AUDIT_WRITE --detach --rm --publish 5003:22 --name tf_linux_test_sshd_ubuntu_22.04 ubuntu_22.04
	@touch $@

test_containers/.run.rockylinux_8:
	@podman run --cap-add AUDIT_WRITE --detach --rm --publish 5004:22 --name tf_linux_test_sshd_rockylinux_8 rockylinux_8
	@touch $@

test_containers/.run.rockylinux_9:
	@podman run --cap-add AUDIT_WRITE --detach --rm --publish 5005:22 --name tf_linux_test_sshd_rockylinux_9 rockylinux_9
	@touch $@

#build_test_container_%: test_containers/.$*

test_containers/.built.%:
	@docker build --tag $* --file test_containers/$*/Dockerfile
	@touch $@

#build_test_containers: build_test_container_ubuntu_18.04 build_test_container_rockylinux_8 build_test_container_rockylinux_9
build_test_containers: \
	test_containers/.built.ubuntu_18.04 \
	test_containers/.built.ubuntu_20.04 \
	test_containers/.built.ubuntu_22.04 \
	test_containers/.built.rockylinux_8 \
	test_containers/.built.rockylinux_9

# terraform targets
init: terraform install ## Run terraform init
	cd terraform; ${terraform} init

plan: terraform install ## Run terraform plan
	cd terraform; ${terraform} plan

apply: terraform install ## Run terraform apply
	cd terraform; ${terraform} apply
