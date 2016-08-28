#!/bin/bash
killall `ps -aux | grep start.sh | grep -v grep | awk '{ print $2 }'`
killall gideon
sleep 1
cd /home/chip/gideon_src;
cp ./gideon ~;
cp ./start.sh ~;
chmod +x ~/start.sh
cp ./update.sh ~;
chmod +x ~/update.sh
cp -r ./fonts ~;
cd ~;
export DISPLAY=:0.0
sudo chmod 777 /dev/i2c-2
sudo setcap CAP_NET_RAW+epi ./gideon
xdotool mousemove 1000 1000
while true; do /home/chip/gideon; sleep 1; done
