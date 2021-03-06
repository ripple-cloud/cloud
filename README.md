# Ripple Cloud

Ripple Cloud is a web service that provides endpoints to manage hubs and apps that runs on them.

Any device that can run [Ripple Hub](https://github.com/ripple-cloud/hub) could be a hub (eg: Raspberry PI).

## Endpoints

### Sign Up (/signup)

To sign up, make a `POST` request to http://[host]:[port]/signup?username=(your username)&password=(your password)&email=(your email)

### Get Authentication Token (/oauth/token)

To get the `access_token`, make a `POST` request to http://[host]/api/oauth/token with the following params:

```
'grant_type'
REQUIRED. Must be 'password'.

'username'
REQUIRED. The username you registered at signup.

'password'
REQUIRED.
```

An example request in curl:

```
curl -X POST 'http://[host]/api/oauth/token?grant_type=password&username=(username)&password=(password)'
```

If the request was successful, you will receive a response with status code `200` and JSON body like:

```
{
  "access_token": "i47vsk9cx3zdrpur1qjditf3to8m3eerwdruagtqwmjg2nhjni7rqcr57p2v23wv",
  "token_type": "bearer",
  "expires_in": 2592000
}
```

If the request was not successful, you will receive a response with status code `400` and JSON body like:

```
{
  "error": "Invalid_request",
  "error_description": "Missing parameter. 'grant_type' is required"
}
```

### Hub

* Register a hub (`POST /api/v1/hub`)
* Retrieve an existing hub (`GET /api/v1/hub/:id`)
* Delete a hub (`DELETE /api/v1/hub`)

### App

* Register an app (`POST /api/v1/app/:slug`)
* Send a request to an app (`POST /api/v1/app/:slug/job`)
* List all datapoints collected from an app (`GET /api/v1/app/:slug/job/id`)
* Delete an app (`DELETE /api/v1/app/:slug`)

## Development

* Install `go get github.com/mattes/migrate`
* Copy `.env-example` to `.env`
  - Set your postgres DB URL
* Export environment: `source .env`
* To run migrations: `make migrate`
