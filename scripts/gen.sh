#!/bin/bash

svcName=${1}

if [ ! -d "app/${svcName}" ]; then
  echo "Directory app/${svcName} does not exist. Creating it..."
  mkdir -p app/${svcName}
fi
cd app/${svcName}
cwgo client -I ../../idl --type RPC --service ${svcName} --module gomall/app/${svcName} --idl ../../idl/${svcName}.proto
cwgo server -I ../../idl --type RPC --service ${svcName} --module gomall/app/${svcName} --idl ../../idl/${svcName}.proto
