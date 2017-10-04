#!/bin/bash
pushd ~/workspace/cf-deployment
  git co v0.10.0
  git pull
popd

cp -r ~/workspace/cf-deployment/* cf-deployment-master/

pushd ~/workspace/cf-deployment
  git co release-candidate
  git pull
popd

cp -r ~/workspace/cf-deployment/* cf-deployment-release-candidate/
