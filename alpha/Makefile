.PHONY: godeps
.PHONY: clean
.PHONY: run
.PHONY: open

all: alpha

clean:
	rm alpha

alpha: main.go 
	docker run \
	-w /go/src/github.com/igneous-systems/alpha \
	-v $(PWD):/go/src/github.com/igneous-systems/alpha \
	gotooling \
	go-with-godep.sh build \
	-o /go/src/github.com/igneous-systems/alpha/alpha \
	/go/src/github.com/igneous-systems/alpha/main.go

run: alpha
	docker build -t alpha .
	docker run --link=etcd:etcd --link=postgres:postgres -p=`boot2docker ip`:8081:80 alpha 

open: alpha
	open http://`boot2docker ip`:8081/index.html

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
	-v $(PWD):/go/src/github.com/igneous-systems/alpha \
	-v $(PWD)/../external:/external \
	-e "GOPATH=/external:/go" \
	-e "IGNEOUS_PROG=alpha" \
	-w /go/src/github.com/igneous-systems/alpha \
	gotooling \
	inside-godep.sh


