#!/bin/bash
export GOPATH=$PWD

chmod +x hooks/pre-commit
ln -s ../../hooks/pre-commit .git/hooks/
