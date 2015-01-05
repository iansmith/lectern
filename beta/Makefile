.PHONY: godeps
.PHONY: clean
.PHONY: run
.PHONY: open

all: beta static/client.js

clean:
	rm beta
	rm static/client.js

beta: main.go 
	docker run \
	-w /go/src/github.com/igneous-systems/beta \
	-v $(PWD):/go/src/github.com/igneous-systems/beta \
	gotooling \
	go-with-godep.sh build \
	-o /go/src/github.com/igneous-systems/beta/beta \
	/go/src/github.com/igneous-systems/beta/main.go

## run the gopherjs compiler
static/client.js: client/clientmain.go
	docker run \
	-v $(PWD):/go/src/github.com/igneous-systems/beta \
	-w /go/src/github.com/igneous-systems/beta \
	gotooling \
	gopherjs-with-godep.sh build \
	-o static/client.js -m \
	github.com/igneous-systems/beta/client

run: beta static/client.js
	docker build -t beta .
	docker run --link=etcd:etcd -v $(PWD)/static:/static -p=`boot2docker ip`:8080:80 beta 

open: beta static/client.js
	open http://`boot2docker ip`:8080/index.html

##
## You only need this target if you wante to *CREATE* a new vendored
## snapshot.  If you want to do that, you have to have all the
## dependencies at your preferred version retrieved into the external
## directory (up one level from here).  Usually you can use "go get" to
## to set that up the way you want.  After everything is set, you can
## use this make target to copy all the content into the Godeps directory
## and re-write the Godeps/Godeps.json file.
##
godep:
	-rm -rf Godeps
	docker run \
	-v $(PWD):/go/src/github.com/igneous-systems/beta \
	-v $(PWD)/../external:/external \
	-e "GOPATH=/external:/go" \
	-e "IGNEOUS_PROG=beta" \
	-w /go/src/github.com/igneous-systems/beta \
	gotooling \
	inside-godep.sh


