#!/bin/bash
export DISPLAY=:0.0
sudo chmod 777 /dev/i2c-2
sudo setcap CAP_NET_RAW+epi ./gideon
xdotool mousemove 1000 1000
while true; do /home/chip/gideon; sleep 1; done
