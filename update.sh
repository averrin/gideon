#!/bin/bash

cd ~/gideon_src;
echo 'Update sources'
git pull;
cd ~/gideon_src/src/github.com/averrin/shodan;
echo 'Update Shodan'
git pull;
cd -;
cd ~/gideon_src/src/github.com/averrin/seker;
echo 'Update Seker'
git pull;
cd -;
echo 'Building...'
./build.sh
