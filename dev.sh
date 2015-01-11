#!/bin/bash
export GOPATH=$PWD

if [ ! -f .git/hooks/pre-commit ]; then
  chmod +x hooks/pre-commit
  ln -s ../../hooks/pre-commit .git/hooks/
fi


