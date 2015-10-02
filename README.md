# CHAOS GOPHER

A collection of unix style tools in GO to do chaos engineering or testing.

### TheSlowness

Coming soon, tool to simulate slowness


## MYSQL Integration Failover tests 


### MYSQL M-S

Run the script in the following order to setup a master-slave MySQL cluster via docker:

1. `cd mysql` - Goes in to the `mysql` folder first.
2. `./start.sh` - Builds and starts 2 MySQL Docker containers. One master, one slave.
3. `docker logs -f mysql-master` - And **wait for master MySQL instance to finish initialization.**
4. `./setup.sh` - Sets up the master-slave relation, adds data, and do a basic query test.
5. `./failover.sh` - To start the interactive `mysqlfailover` monitor and auto-failover process.

### MYSQL M-M

The scripts are the same as the basic MySQL setup. Except that there is no failover step
as all instances are considered master.

### ETCD (dockerized)

1. `cd etcd` - Go to etcd folder.
2. `./build.sh` - Builds required Docker images.
3. `./start.sh` - Starts the main ETCD instance.
4. `./seed.sh` - Seed initial data to ETCD. Run **again** if you encounter a timeout
   error.

Use `./ctl.sh` script when you need to run `etcdctl`. Example:

```sh
$ ./ctl.sh get /chaostesting/datasources
[ { "name": ...
```

### ETCD (local)

* Install via Homebrew `brew install etcd`
* Gets etcdctl via go get `go get github.com/coreos/etcd/client`
* Uses `seed-local.sh` instead of `seed.sh`

### SERVER (dockerized)

1. `cd app` - Go to the app folder.
2. `./build.sh` - Build app's Docker image.
3. `./start.sh` - Start app in Docker container. Keep the foreground process running.
4. `./start-tester.sh` - Start tester app in Docker container.

**NOTE:** The tester app should *not* crash. If it crashes, this mean something is not
configured correctly. For example, ETCD stalling may cause configuration settings to be
inconsistent inside the main app.

### SERVER (local)

Use `/run*.sh` script for local runs. Configurations needs to be updated since most
components will use the `docker0` network interface otherwise.

1. Start MySQL and ETCD containers as usual.
2. `docker pull hyperworks/disk-filler` - This image is required by the tester.
3. `cd etcd` - Go to ETCD folder.
4. `./seed-local.sh` - Seed configuration for local runs.
5. `cd ../app` - Go to app folder.
6. `export API_ENDPOINT=http://0.0.0.0:8080` - Points API endoint on local machine.
7. `export ETCD_ENDPOINT=http://0.0.0.0:2379` - Points ETCD endpoint on local machine.
8. `./run.sh` - Starts the app.
9. `./run-loader.sh` - Starts the loader (or tester with `run-tester.sh`).

### RESET

To destroy all containers in your system and start over from scratch.

```
docker rm -fv `docker ps -aq`
```
