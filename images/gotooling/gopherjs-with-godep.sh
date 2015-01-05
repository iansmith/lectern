#!/bin/sh
GOPATH=/go:`godep path` gopherjs $@
