# LECTERN: Getting you up to speed on deis

## Pre-requisite: Docker

### Docker on Mac
You have to have docker working on your machine.   In this directory there is 
shell script in this directory called `enable-lectern` that you can start 
from with for setting your environment variables.  These are set up for
a modern version of [boot2docker](http://boot2docker.io).  

Boot2docker is  superior to vagrant if you are going to have many source 
directories  since they can all "share" the same virtual machine running docker. 
You can bring up the virtual machine with 
```
boot2docker init
boot2docker up
```

It is hard, but not impossible, to get this convenient sharing with vagrant 
since the `Vagrantfile`'s directory is important to which virtual machine is used.  
On modern versions of boot2docker, your home directory should be "mounted" 
into the virtual machine. You can test this with 

```
boot2docker ssh
cd /Users/yourusername
ls
```
You can test that docker is working ok with `docker version` and make sure you
get no errors and a report that looks like this:
```

Note that the boot2docker mechanism of mounting your home directory doesn't use
NFS so the recent problems with NFS stale file handles is avoided at the cost
of slower reads/writes.

### Docker on Linux
On a linux box, no configuration of docker should be necessary.  All docker
commands should be using the unix domain socket.

## Pre-requisite: Deis

The deis command line tool `deis` is similar to 
[Heroku Toolbelt](https://toolbelt.heroku.com/auth/heroku).  You can download
the deis tool like this:

```
 curl -sSL http://deis.io/deis-cli/install.sh | sh
 ```

You need to move it somewhere to be in your PATH, perhaps into `/usr/local/bin`
like this `ln -fs $PWD/deis /usr/local/bin/deis`. 

## Creating Your User Account

Use `deis register http://deis.apps.iggy.buzz` to register yourself an
account.  Don't use any email address when creating the account or things can
go haywire.  When you have succeceded, you should see:
```
Registered iansmith
Logged in as iansmith
```

You need to register the *public* key associated with your account with deis 
so various git operations will work.  You can do this with
`deis keys:add ~/.ssh/id_rsa.pub` or wherever your public key is stored. You
must use the corresponding private key with various git operations in later
steps, so it's best to use your "default" identity.

You are now ready to go!

## Setting up the local tooling

In the main directory, `make setup`.  This builds some images with `docker build`
that are needed to test/run locally; it also starts two containers running.  

The two containers that are left running are called `etcd` and `postgres`, you
can see them with `docker ps` and they can be ignored once started.  Note that if
you do `make setup` again it destroys the contents of the database and
the etcd key/value.

All this tooling for local only to allow your local workstation to (cheaply)
simulate the staging/production cluster running deis.

## Creating the beta app

You must be logged into our deis cluster to execute these commands. You can
check your status with `deis auth:whoami` and if necessary you can login
with the identity created above with `deis auth:login`.

```
cd beta
deis apps:create
```

and you should get some output that includes your new application's name,
which is a pair of words connected by a hyphen. We'll refer to that name
as `app-name` in this document:

```
Creating application... done, created rubber-yachting
Git remote deis added
```

This also adds a git remote that you can get information about with 
`git remote -v`.

## Building the beta app

```
cd beta
make
```

This make command requires some explanation.  First, the go tooling and
such are located in `/opt/go` inside the image that is used to build the code;
the image name is `gotooling` if you look in `docker images`.  
All go code for use in the build is either mounted (a "volume mount" in docker) 
into the directory `/go` _or_ provided by vendored godeps.  These are included 
in the source in `beta/Godeps`. 

Note that care is taken to make it "look like" we are using the standard go project
layout with the code for beta in `/go/src/github.com/igneous-systems/beta`.
Again, that is done with volume mounts for the benefit of the build.

To use other "build tools" that build css or html pages, you'll need to
make sure the gotooling image has these tools (see the `Dockerfile` in
`images/gotooling`). These tools should be in the path and
you should invoke them in a way that generates an _output_ to the directory
they are being run in.  This is necessary since that directory is not really 
"inside" the container. This is done with `gopherjs` in this example program.

## Running the beta app 
### Running the beta app on OSX
Assuming you are running boot2docker, you should do `boot2docker ip` to get
the ip address of your virtual machine. The makefile will do this as well,
so you should have boot2docker in your path.

Then you can use:

```
make run
```
This will build a docker image that includes the beta binary, setup some
services that are needed in the container, and run it.  It exposes port
8080 to the ip address of the boot2docker vm (see the `makefile`), so you 
should be able to see something on `http://192.168.59.103:8080/` or whatever
IP your boot2docker is on.

Note that `make run` tries to build before it runs, so it can be just run 
whenever you want to restart.

`make open` is a shortcut to open the browser to the beta application page
on OSX (only).

### Running the beta app on linux or "by hand"
For some _hostip_ and _hostport_:

```
cd beta
make
docker build -t beta .
docker run --link etcd:etcd -v $PWD/static:/static -p=hostip:hostport:80 beta
```

### Configuration of the Database Params With Beta
-----------------------------------------------

"beta" is a simple AJAX app for setting the configuration parameters that 
be used for the database access by the alpha application (see below).  For
this demo, we are assuming you need to poke the username and password to use
into the `etcd` storage.  In a real application, all of this configuration
should probably be done with environment variables as this is both more
12-factorish as well as trackable in git (assuming you are using deis!).

You can type a username and password into the form provided to change the
settings that will be used by the "alpha" application.  Errors are reported
in the page.  _Set the username and password to *postgres* and *seekret*_. 
At least set it to that if you want the alpha application to work! You may
find it interesting to set this to "wrong" values as well and watch what
happens to alpha.  The username and password is burned into the 
database image in use with this demo and is not easily changed (see
`images/database/provision.sh` and `images/database/provision.sql`).

