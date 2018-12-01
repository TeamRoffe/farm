#!/bin/bash
ssh pi@192.168.2.155 killall farm
ssh pi@192.168.2.155 rm -f /home/pi/farm
GOOS=linux GOARCH=arm go build -o farm src/main.go
scp farm config.ini pi@192.168.2.155:/home/pi/
ssh pi@192.168.2.155 /home/pi/farm