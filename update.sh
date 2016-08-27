#!/bin/bash

cd ~/gideon_src;
git pull;
cd ~/gideon_src/src/github.com/averrin/shodan;
git pull;
cd -;
cd ~/gideon_src/src/github.com/averrin/seker;
git pull;
cd -;
./build.sh;
killall `ps -aux | grep start.sh | grep -v grep | awk '{ print $1 }'`
killall gideon
sleep 1
cp ./gideon ~;
cp ./start.sh ~;
chmod +x ~/start.sh
cp ./update.sh ~;
chmod +x ~/update.sh
cp -r ./fonts ~;
cd ~;
~/start.sh
