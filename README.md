# REST API photo gallery on golang

## Requirements

    go 1.14+
    sqlite 3.28+

## Run the app

    go run .

# REST API

## Get list of photos

### Request

`GET /photos`

    curl -i -H 'Accept: application/json' http://localhost:8081/photos

### Response

    HTTP/1.1 200 OK
    Date: Sat, 24 Oct 2020 21:15:01 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 2

    [
        {
            "PhotoUrl": "localhost:8081/photos/c42eacfb9b85134387aa7870",
            "PhotoPreviewUrl": "localhost:8081/photos/c42eacfb9b85134387aa7870/preview"
        },
        {
            "PhotoUrl": "localhost:8081/photos/5c679707325d21ed3d5525fb",
            "PhotoPreviewUrl": "localhost:8081/photos/5c679707325d21ed3d5525fb/preview"
        },
        {
            "PhotoUrl": "localhost:8081/photos/0cf279bf497527df992ee782",
            "PhotoPreviewUrl": "localhost:8081/photos/0cf279bf497527df992ee782/preview"
        }
    ]

## Upload a new photo

### Request

`POST /photos/upload/`

    curl --location --request POST 'http://localhost:8081/photos/upload' --header 'Content-Type: application/json' --form 'uploadFile=@/home/user/downloads/file.jpg'

### Response

    HTTP/1.1 201 Created
    Date: Sat, 24 Oct 2020 21:23:50 GMT
    Status: 201 Created
    Connection: close
    Content-Type: application/json
    Content-Length: 136

    {"PhotoUrl":"localhost:8081/photos/9518b314a6ee537212683c88","PhotoPreviewUrl":"localhost:8081/photos/9518b314a6ee537212683c88/preview"}

## Delete photo

### Request

`DELETE /photos/:hash`

    curl -i -H 'Accept: application/json' http://localhost:8081/photos/377d1b42e37ee4c96c472a22

### Response

    HTTP/1.1 200 OK
    Date: Sat, 24 Oct 2020 21:27:02 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 36

## Get photo

### Request

`GET /photos/:hash`

    curl -i -H 'Accept: application/json' http://localhost:8081/photos/afd0dba57a9dc4003d970c90

### Response

    HTTP/1.1 200 OK
    Date: Sat, 24 Oct 2020 21:28:01 GMT
    Status: 200 OK
    Connection: close
    Content-Type: image/jpeg

## Get photo preview

### Request

`GET /photos/:hash`

    curl -i -H 'Accept: application/json' http://localhost:8081/photos/afd0dba57a9dc4003d970c90/preview

### Response

    HTTP/1.1 200 OK
    Date: Sat, 24 Oct 2020 21:28:01 GMT
    Status: 200 OK
    Connection: close
    Content-Type: image/jpeg
