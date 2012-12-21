#!/bin/bash

export GOPATH=$(pwd)

for i in `ls src`
do
	echo $i
	go install $i
done
