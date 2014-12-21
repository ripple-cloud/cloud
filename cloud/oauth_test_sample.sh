#!/bin/bash
#
# Test script to run against api/oauth/token endpoint.
# Expects 200 http status code.

function test {
    res=`curl -s -I -X POST $1 | grep HTTP/1.1 | awk {'print $2'}`
    if [ $res -ne 200 ]
    then
        echo "FAIL!"
        echo "Error $res on $1"
        curl -X POST $1 
    else
        echo "PASS: $res"
    fi
}  

test 'http://[host]/api/oauth/token?grant_type=password&username=[testuser]&password=[testpassword]'
