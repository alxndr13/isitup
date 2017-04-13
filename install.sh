#!/bin/bash
echo 'Copying Servicefile to /etc/systemd/system/'
sudo mv isitup.service /etc/systemd/system/
export GOPATH=/root/go
echo 'You need the go runtime installed to build your Version of isitup.'
echo 'trying to build...'
go build
echo 'reloading systemd daemon'
sudo systemctl daemon-reload
echo 'Setting config Directory..'
sudo mkdir -p /etc/isitup/
sudo mv settings.toml /etc/isitup/
