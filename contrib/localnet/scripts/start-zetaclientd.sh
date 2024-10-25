#!/bin/bash

# This script is used to start ZetaClient for the localnet
# An optional argument can be passed and can have the following value:
# background: start the ZetaClient in the background, this prevent the image from being stopped when ZetaClient must be restarted

/usr/sbin/sshd

HOSTNAME=$(hostname)
OPTION=$1

# read HOTKEY_BACKEND env var for hotkey keyring backend and set default to test
BACKEND="test"
if [ "$HOTKEY_BACKEND" == "file" ]; then
    BACKEND="file"
fi

cp  /root/preparams/PreParams_$HOSTNAME.json /root/preParams.json
num=$(echo $HOSTNAME | tr -dc '0-9')
node="zetacore$num"

echo "Wait for zetacore to exchange genesis file"
sleep 40
operator=$(cat $HOME/.zetacored/os.json | jq '.ObserverAddress' )
operatorAddress=$(echo "$operator" | tr -d '"')
echo "operatorAddress: $operatorAddress"
echo "Start zetaclientd"
if [ $HOSTNAME == "zetaclient0" ]
then
    rm ~/.tss/*
    MYIP=$(/sbin/ip -o -4 addr list eth0 | awk '{print $4}' | cut -d/ -f1)
    zetaclientd init  --zetacore-url zetacore0 --chain-id athens_101-1 --operator "$operatorAddress" --log-format=text --public-ip "$MYIP" --keyring-backend "$BACKEND"
    zetaclientd start < /root/password.file
else
  num=$(echo $HOSTNAME | tr -dc '0-9')
  node="zetacore$num"
  MYIP=$(/sbin/ip -o -4 addr list eth0 | awk '{print $4}' | cut -d/ -f1)
  SEED=""
  while [ -z "$SEED" ]
  do
    SEED=$(curl --retry 10 --retry-delay 5 --retry-connrefused  -s zetaclient0:8123/p2p)
  done
  rm ~/.tss/*
  zetaclientd init --peer /ip4/172.20.0.21/tcp/6668/p2p/"$SEED" --zetacore-url "$node" --chain-id athens_101-1 --operator "$operatorAddress" --log-format=text --public-ip "$MYIP" --log-level 1 --keyring-backend "$BACKEND"
  zetaclientd start < /root/password.file
fi

if [ "$OPTION" == "background" ]; then
    sleep 3
    tail -f $HOME/zetaclient.log
fi

