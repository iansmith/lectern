#!/bin/sh -x

if [ ! -n "$IGNEOUS_PROG" ]; then
	echo "IGNEOUS_PROG environment variable not set"
	exit 1
fi

if [ ! -d "/go/src/github.com/igneous-systems/$IGNEOUS_PROG" ]; then
	echo "can't understand IGNEOUS_PROG environment variable"
	exit 1
fi

if [ -d "/go/src/github.com/igneous-systems/$IGNEOUS_PROG/.git" ]; then
	echo "found a .git in the IGNEOUS_PROG directory, so aborting..."
	exit 1
fi

## this horrible hack is necessary because of the fact that we are doing
## tricks with the gopath that confuse godeps.  it wants there to be a
## a version control system repository at the level of the project
## that you are running godeps at. 

cd /go/src/github.com/igneous-systems/$IGNEOUS_PROG/
git init . > /dev/null
git config --global user.name "Godep NeedsWork"
git config --global user.email "brokenhack@example.com"
git commit --allow-empty -m "no message" > /dev/null

cd /go/src/github.com/igneous-systems/$IGNEOUS_PROG/
godep save ./...

OK="y"
if [ "$?" != "0" ]; then
		echo "******* "
		echo "******* godep failed!"
		echo "******* "
		OK="n"
fi

#scary
rm -rf .git

##show us what happened
if [ "$OK" = "n" ]; then
	exit 1
fi

cat /go/src/github.com/igneous-systems/$IGNEOUS_PROG/Godeps/Godeps.json
