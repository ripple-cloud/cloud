To sign up: http://[host]:[port]/signup?username=(your username)&password=(your password)&email=(your email)

To request a token: http://[host]:[port]/api/oauth/token?username(your username)&password=(your password)&grant_type=password

To get the `access_token`, make a `POST` request to http://[host]/api/oauth/token with the following params:

`grant_type`
REQUIRED. Must be `password`.

`username`
REQUIRED. The username you registered at signup.

`password`
REQUIRED.

An example request in curl:

```
curl -X POST 'http://[host]/api/oauth/token?grant_type=password&username=(username)&password=(password)'
```

If the response is successful you will receive a `200` response with a JSON body like:

```
{
  "access_token": "i47vsk9cx3zdrpur1qjditf3to8m3eerwdruagtqwmjg2nhjni7rqcr57p2v23wv",
  "token_type": "bearer",
  "expires_in": 2592000
}
```

If the response is not successful, you will receive a `400` response with the following JSON like:

```
{
  "error": "Invalid_request",
  "error_description": "Missing parameter. 'grant_type' is required"
}
```

