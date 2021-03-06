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
The API requires that a config.json file live in the root directory with the correct S3, DB, MailGun and App properties defined
```json
{
  "S3": {
    "Key": "my-special-key",
    "Secret": "sUp3rS3cre7",
    "Bucket": "some-bucket",
    "Region": "us-west-2"
  },
  "MySQL": {
    "Host": "some.db.net",
    "Database": "api_db",
    "User": "username",
    "Password": "3245ertfdsa"
  },
  "MailGun": {
    "Domain": "YourDomain",
    "ApiKey": "YourApiKey",
    "PublicApiKey": "YourPublicApiKey"
  },
  "App": {
    "Domain": "http://localhost:3000"
  }
}
```

## Public API Endpoints
- http://localhost:3000/api/v1/public/register
    - `POST`: Register the user with the API
    - `{"email":"some@email.com","name":"User X"}`
    - returns the user with an id after it is created (should be used for auth)
    
- http://localhost:3000/api/v1/public/drawing/{id}
    - `GET`: Get a specific drawing
    - returns the associated drawing

- http://localhost:3000/api/v1/public/drawings
    - `GET`: Get all of the drawings
    - returns an array of all the drawings (could be used as a gallery)


## Secured API Endpoints
End points with a _private_ base, require the header `X-App-User` to be provided.  
The `X-App-User` header should contain the user json object base64 encoded  
`X-App-User` header example : `eyJuYW1lIjoiQmVuIiwgImVtYWlsIjoiYmVuQGtyYXNoaWRidWlsdC5jb20iLCAiaWQiOiAxfQ==`
- http://localhost:3000/api/v1/private/section/type
    - `GET`: Get the section _type_ that should be drawn.
    - `"top"`, `"middle"`, `"bottom"`
    
- http://localhost:3000/api/v1/private/section/{type}
    - `POST`: Post a base64 encoded PNG that has been drawn for a section
    - `data:image/png;base64,iVkhdfjdAjdfirtn=`
    - returns the drawing associated with the section that was posted

- http://localhost:3000/api/v1/private/drawings
    - `GET`: Get all of the _drawings_ that the user has participated in.
    - `[{"id": 1, "url": null, "top": { "name": null, "email": null, "url": null }, "middle": { "name": null, "email": null, "url": null }, "bottom": { "name": null, "email": null, "url": null }}]`

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
