#!/bin/bash

if ! command -v curl &> /dev/null
then
    echo -e "Could not run script cause cURL command not found.\nInstall curl to run this script."
    exit 1
fi;

if ! command -v jq &> /dev/null
then
    echo -e "Could not run script cause jq command not found.\nInstall curl to run this script."
    exit 1
fi;

# build project
make;

# start server
echo "Building server...";
./xkcd-server -p 8080 &> /dev/null &

sleep 5;
# login
echo "Requesting token...";
token=$(curl -s --request POST \
      --data '{"name": "admin", "password": "admin"}'\
      localhost:8080/login \
      | jq -r '.token');
echo "> Got token: $token";

# update database
echo "Requesting update...";
added=$(curl -s --request POST \
      -H "Authorization: $token" \
      localhost:8080/update);
echo "> Update response: $added"

# search pics
echo "Requesting pics...";
response=$(curl -s --request GET \
           -H "Authorization: $token" \
           localhost:8080/pics?search="apple,doctor");

# check comic presence
comic="https://imgs.xkcd.com/comics/an_apple_a_day.png";
if [[ $response == *"comic"* ]]; then
    echo "SUCCESS: '$comic' is found."
else
    echo -e "FAIL: '$comic' not found in response:\n'$response'"
fi;

kill %1;