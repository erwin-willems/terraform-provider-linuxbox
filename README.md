# Introduction

This provider can configure a linux VM.

Inspired by:
- https://github.com/numtide/terraform-provider-linuxbox
- https://github.com/mavidser/terraform-provider-linux

# Development

## Check if code builds
```sh
make dev
```

Note: output should be:
```
This binary is a plugin. These are not meant to be executed directly.
Please execute the program that consumes these plugins, which will
load any plugins automatically
exit status 1
make: *** [run] Error 1
```

## Test with real terraform
```sh
make clean init apply
```

# Run tests

## Acceptance test
```sh
$ make test_acc_all
```

## Acceptence test on a specific system:
Acceptance test can be run on a specific system. Available systems:

- test_ubuntu_18.04
- test_ubuntu_20.04
- test_ubuntu_22.04
- test_rockylinux_8
- test_rockylinux_9

```sh
$ make test_rockylinux_8 
```

## Run single test
```sh
$ make test_acc_all test_args="-run TestAccUserCreation
```

## Run multiple tests
In this example, run all Acceptance tests for the User resource
```sh
$ make test_acc_all test_args="-run TestAccUser
```



## Run single test on a specific system
```sh
$ make test_rockylinux_8 test_args="-run TestAccUserCreation"
```

## Run multiple tests
In this example, run all Acceptance tests for the User resource
```
$ make test_rockylinux_8 test_args="-run TestAccUser"
```


## Enable verbose logging in Terraform:
```sh
$ export TF_LOG="info" 
```
