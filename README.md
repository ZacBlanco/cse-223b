# cse-223b

A Stateful, Serverless Actor System

Built upon a modified version of [Apache Openwhisk](https://openwhisk.apache.org)

- Openwhisk commit: `3802374d58d87fc6a95477929fc67269d6dcfe2c`
- Openwhisk Go Runtime (1.17.0) commit: `13be9e24885fc2b5c4c5093bb204da8395f0f401`


## Running the system

- This system was built and run with the following environment

- Ubuntu 20.04.2 LTS
  - Linux 5.4.0-74-generic
- docker 20.10.7
- containerd 1.5.2
- CRIU 3.15

Notes:

- The kernel MUST be at least 5.4.0-74. There is a critical patch for the
  overlayfs driver. Otherwise performance is awful and the `vfs` storage driver
  must be enabled for docker
- docker 20.10.7 uses `containerd` 1.4.6. Checkpoint and restore functionality
  is broken in this environment. To fix this, [download
  `containerd-1.5.2`](https://github.com/containerd/containerd/releases/tag/v1.5.2)
  from the github releases and install it under `/opt/containerd-1.5.2/` and
  then change the systemd `containerd` unit file
  (`/lib/systemd/system/containerd.service`)to point to the new executable.

For more information on these compatiblity issues, please read the following links:

- https://github.com/checkpoint-restore/criu/issues/1223
- https://github.com/checkpoint-restore/criu/issues/860
- https://github.com/checkpoint-restore/criu/issues/1316
- https://git.launchpad.net/~ubuntu-kernel/ubuntu/+source/linux/+git/focal/commit/?h=master-next&id=28eab192cf0e37156fc41b36f06790d5ca984834
- https://github.com/moby/moby/issues/41602
- https://github.com/moby/moby/issues/37344

Assuming all of the pre-requisites are satisfied:

- Build the actor container

```console
$ cd openwhisk-runtime-go && ./gradlew distDocker
```

- Build and run the standalone openwhisk system

```console
$ cd openwhisk && ./gradlew :core:standalone:bootRun
```

You can test a basic stateful actor under `tests/basic-actor`. You'll need the
openwhisk CLI, `wsk`, installed.
