#!/bin/bash

if ! command -v curl &> /dev/null
then
    echo "Could not run script cause cURL command not found"
    exit 1
fi;

if ! command -v jq &> /dev/null
then
    echo "Could not run script cause jq command not found"
    exit 1
fi;

# build project
make;

# start server
./xkcd-server -p 8080 &> /dev/null &

sleep 5;
# login
echo "token requested";
token=$(curl -s --request POST \
      --data '{"name": "admin", "password": "admin"}'\
      localhost:8080/login \
      | jq -r '.token');
echo "got token: '$token'";

# update database
echo "update requested";
curl -s --request POST \
      -H "Authorization: $token" \
      localhost:8080/update;

# search pics
echo "pics requested";
response=$(curl -s --request GET \
           -H "Authorization: $token" \
           localhost:8080/pics?search="apple,doctor");

# check comic presence
comic="https://imgs.xkcd.com/comics/an_apple_a_day.png";
if [[ $response == *"comic"* ]]; then
    echo "'$comic' found"
else
    echo "'$comic' not found in response:"
    echo "$response"
fi;

kill %1;