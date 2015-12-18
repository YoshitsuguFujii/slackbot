#!/bin/bash

cd `dirname $0`
cd ../
source .envrc
./slackbot
exit 0
