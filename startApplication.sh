#!/bin/bash  
# 웹애플리케이션 실행

pushd application

rm -rf wallet/*

./ccp/ccp-generate.sh

sleep 2

node enrollAdmin.js

sleep 2

node registerUser.js

sleep 2

node server.js

popd

