# kubedebugger

A simple tool that allows launching ephemeral containers from a YAML file to debug the running GO processes in the containers which runs on any Kubernetes environment.

## Installation

`kubedb` is currently available from this repository. A snap package will be available soon.


### Snap

TODO

### Manual

Get the binary `build/kubedb` and copy it under `/usr/local/bin/` or your preferred location.

```sh
$ make build
$ chmod +x build/kubedb
$ sudo cp build/kubedb /usr/local/bin/kubedb 
```

## Prerequisites

The following requirements should be met:
1. `kubeconfig` file of your Kubernetes cluster should be saved under $HOME/.kube/config: `sudo microk8s kubectl config view --raw > $HOME/.kube/config`
2. The target Go application should be built with `-gcflags='all=-N -l'` before running it.
3. The target Go application should be running and target container should allow you to execute shell commands to get process ID (PID).
4. A YAML file that describes the ephemeral container to launch is required and template is provided in this repository (pls see `/KubeDebugger/ephemeral.yaml`).

## Usage

1. Make sure that target pod is `Running`.

```sh
$ kubectl get pod < pod name > -n < namespace >
```

Sample execution:

```sh
$ kubectl get pod pcf-0 -n sdcore
pod/pcf-0                            2/2     Running   0              68m
service/pcf                                  ClusterIP      10.152.183.113   <none>        65535/TCP,8080/TCP,29507/TCP            25h
```

2. Exec into target pod to get the target PID to debug.

```sh
$ kubectl exec -it < pod name > -n < namespace > -c < container name > -- sh
```

Sample execution:
```sh
$ kubectl exec -it pcf-0 -n sdcore3 -c pcf -- sh
$ ps -ef
UID          PID    PPID  C STIME TTY          TIME CMD
root           1       0  0 10:12 ?        00:00:02 /charm/bin/pebble run --create-dirs --hold --http :38813 --verbose
root          54       1  0 10:13 ?        00:00:01 /bin/pcf --pcfcfg /etc/pcf/pcfcfg.yaml
```

The target process ID is `54` according to this environment.

3. Prepare an ephemeral pod YAML file by inserting the correct process ID belongs to a Go binary.
Ephemeral pod template named `ephemeral.yaml` is provided in this repository.

Sample pod YAML:

```yaml
name: delve
image: gatici/delve:1.23
securityContext:
  privileged: true
command:
  - dlv
  - --listen=127.0.0.1:2345
  - --headless=true
  - --accept-multiclient
  - --api-version=2
  - attach
  - '54'
```

4. Run the following command by providing the necessary inputs.

```sh
$ kubedb <target pod name> -f <path to ephemeral container>.yaml -c <target container name> -n <namespace>
```

Sample execution:

```sh

$ kubedb pcf-0 -f ephemeral.yaml -c pcf -n sdcore
EphemeralContainer/delve-xh87wm created
```

5. Check the target container to see the Delve process which listens the provided PID.

```sh
$ kubectl exec -it pcf-0 -n sdcore3 -c pcf -- sh
$ ps -ef
UID          PID    PPID  C STIME TTY          TIME CMD
root           1       0  0 10:12 ?        00:00:03 /charm/bin/pebble run --create-dirs --hold --http :38813 --verbose
root          54       1  0 10:13 ?        00:00:01 /bin/pcf --pcfcfg /etc/pcf/pcfcfg.yaml
root         161       0  0 10:23 ?        00:00:02 dlv --listen=127.0.0.1:2345 --headless=true --accept-multiclient --api-version=2 attach 54
```

6. Do kubectl port forwarding to access Delve server which is running in a Kubernetes cluster from our local machine or another external system.

```sh
$ kubectl port-forward  < pod name >  < Delve server port >  < Target application port > -n < namespace >
```

Sample execution:

```sh
$ kubectl port-forward  pcf-0  2345:2345 29507:29507 -n sdcore3
Forwarding from 127.0.0.1:2345 -> 2345
Forwarding from [::1]:2345 -> 2345
Forwarding from 127.0.0.1:29507 -> 29507
Forwarding from [::1]:29507 -> 29507
Handling connection for 2345
```

7. Connect to Delve debugger using an IDE or any CLI.

```sh
$ dlv connect 127.0.0.1:2345
Type 'help' for list of commands.
(dlv) b main.main
Breakpoint 88 set at 0xfdee16 for main.main() /root/parts/pcf/build/pcf.go:71
(dlv) 
```
