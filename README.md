# RESTful API for skwiz.it
This is a relatively simple API thrown together for the _"exquisite corps"_ app skwiz.it.

## Install and Run
```shell
$ go get github.com/benkauffman/skwiz-it-api

$ cd $GOPATH/src/github.com/benkauffman/skwiz-it-api
$ go get && go build
$ chmod +x ./skwiz-it-api && ./skwiz-it-api
```

## Configure
An example configuration is listed below and is also saved under example.config.json  
The API requires that a config.json file live in the root directory with the correct S3 and DB properties defined
```json
{
    "S3": {
      "key": "my-special-key",
      "secret": "sUp3rS3cre7",
      "bucket": "some-bucket",
      "region":"us-west-2"
    },
    "MySQL": {
      "host": "some.db.net",
      "database": "api_db",
      "user": "username",
      "password": "3245ertfdsa"
    }
  }
```

## Secured API Endpoints
End points with a _private_ base, require the header `X-App-User` to be provided.  
The `X-App-User` header should contain the user json object base64 encoded
`X-App-User` header example : `eyJuYW1lIjoiQmVuIiwgImVtYWlsIjoiYmVuQGtyYXNoaWRidWlsdC5jb20iLCAiaWQiOiAxfQ==`
- http://localhost:3000/api/v1/private/section
    - `GET`: Get the section _type_ that should be drawn.
    - `"top"`, `"middle"`, `"bottom"`
    
- http://localhost:3000/api/v1/private/section/{type}
    - `POST`: Post a base64 encoded PNG that has been drawn for a section
    - `data:image/png;base64,iVkhdfjdAjdfirtn=`
    - returns the drawing associated with the section that was posted

## User
```json
{
    "id": 1,
    "name": "Ben Kauffman",
    "email": "ben@krashidbuilt.com"
}
```

## Drawing
```json
{
    "id": 1,
    "url": null,
    "top": {
        "name": null,
        "email": null,
        "url": null
    },
    "middle": {
        "name": null,
        "email": null,
        "url": null
    },
    "bottom": {
        "name": null,
        "email": null,
        "url": null
    }
}
```
