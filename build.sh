#!/bin/bash
GOOS=linux GOARCH=arm go build -o farm src/main.go
scp farm pi@192.168.2.155:/home/pi/farm
ssh pi@192.168.2.155 /home/pi/farm