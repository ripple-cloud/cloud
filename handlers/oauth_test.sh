#!/bin/bash
#
# Test script to run against api/oauth/token endpoint.
# Expects 200 http status code.

function test {
    res=`curl -s -I -X POST $1 | grep HTTP/1.1 | awk {'print $2'}`
    if [ $res -eq 200 ]
    then
      echo "PASS: $res"
    elif [ $res -ne 200 ]
    then
        echo "FAIL!"
        echo "Error $res on $1"
        curl -X POST $1 
    else
        echo "UNKNOWN ERROR/NO HTTP RESPONSE: Make sure your server is running."
    fi
}  

test "http://$RIPPLE_HOSTNAME/oauth/token?grant_type=password&username=$RIPPLE_TEST_USERNAME&password=$RIPPLE_TEST_PASSWORD"
