To sign up: `http://[host]:[port]/signup?username=(your username)&password=(your password)&email=(your email)`
Returns: successful or unsuccessful signup.

To create a token: `http://[host]:[port]/tokencreate?username(your username)&password=(your password)`
Calling this endpoint each time returns a new token which will replace the previous one.
Returns: {username: (your username), token: (new token)}

To request your existing token: `http://[host]:[port]/tokenrequest?username(your username)&password(your password)`
Returns: {username: (your username), token: (existing token)}

TODO: add granttype and refresh token params (oauth doc)?
add UTF-8 character encoding?
