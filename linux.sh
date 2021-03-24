#!/bin/bash

export GOOS=linux
go build .
export GOOS=windows

sleep 5