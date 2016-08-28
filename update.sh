#!/bin/bash

cd ~/gideon_src;
git pull;
cd ~/gideon_src/src/github.com/averrin/shodan;
git pull;
cd -;
cd ~/gideon_src/src/github.com/averrin/seker;
git pull;
cd -;
./build.sh
