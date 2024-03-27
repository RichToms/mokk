# Mokk

Mokk is a CLI library that allows you to create mock APIs from a config file.
This application was built to improve the process of integrating with third-party APIs during development.

## Installation

You can install this CLI tool if you have Go installed and your GOPATH correctly configured.

```shell
go install github.com/richtoms/mokk@latest
```

## How to run Mokk

### From Docker (Recommended)

If you are using Docker, you can start your Mokk server without even needing to install the CLI.

It is recommended to create your own config file and provide it to the container, otherwise it will start with its own example Users API.

```shell
docker run \
  -p 8080:80 \
  --volume my-app.mokk.yml:/app/mokk.yml \
  richtoms/mokk:v1.1
```

Alternatively, you can use one of the built-in example APIs found in the `./examples` directory.

```shell
docker run \
  -p 8080:80 \
  richtoms/mokk:v1.0 \
  ./mokk start --config=./examples/aws.apigateway-ws.yml
```

### Via CLI

```shell
mokk generate example
mokk start --port=8080 --config=example.mokk.yml
```

These commands should provide an example config file for you to get started.
If you open another terminal you should now be able to reach the server via cURL:

```shell
curl -i -X "GET" http://localhost:8080/
```

## Configuration

Mokk is config-driven by design, using YAML to be developer-friendly. Below is an example of a Mokk config 
file to help get you started:

```yaml
name: Mokk Example Server

options:
  printRequestBody: true
  
routes:
  - path: "users"
    method: "GET"
    statusCode: 200
    response: '{"status":"Success","users":[{"name":"MockZilla"}]}'

  - path: "users/:id"
    method: "GET"
    statusCode: 200
    response: '{"status":"Success","user":{"name":"MockZilla"}}'
    delay: 100
```

### Routes
#### Path

The path field follows [GoFiber's routing patterns](https://docs.gofiber.io/guide/routing#paths) therefore you can utilise wildcards in your paths.

**Note:** it is important when defining a mix of route params and static values to define the static paths first. Due to Fiber's routing system it will struggle to match your static path when a dynamic is defined first.

#### Method

The method field *should* be one of the major HTTP verbs:

- `GET`
- `POST`
- `PUT`
- `PATCH`
- `DELETE` 
- `HEAD`
- `OPTIONS`

#### Response

Mokk currently only supports JSON APIs, therefore this property must contain some form of JSON array/object.

#### Variants

Mokk supports multiple variants for the same path to allow you to have multiple responses to interact with easily. The provided fields
of a variant's parameters must match the named route parameter for a match to be made.

Below is an example config of a route definition with variants:

```yaml
  - path: "users/:user"
    method: "GET"
    statusCode: 200
    response: '{"status":"Success","user":{"id":1,"name":"MockZilla"}}'
    variants:
      # /users/123 will return an alternative success
      - params:
          user: 123
        statusCode: 200
        response: '{"status":"Success","user":{"id":123, "name":"MockZilla Jr."}}'
      # /users/999 will return 404
      - params:
          user: 999
        statusCode: 404
        response: '{"status":"Failure"}'
```

This feature can be used to provide failure states within the API, along with multiple success states to further your test scenarios.

#### Delay

This feature will allow you to add a number of milliseconds delay to Mokk's response in order to simulate actual response times and even timeouts in your API clients.

### Options

#### printRequestBody

If enabled, any request with a JSON body will have the contents printed to the console. This feature
can be useful when debugging what your application is sending to the mocked API.

#### trackRequests

If enabled, any request made to a configured endpoint will be tracked and retrievable from the System API. 

This feature is disabled by default.


## System API

Mokk provides an API that can provide you with logs of any captured requests while the server is running.

The following paths are implemented:

- `/_mokk/requests`
- `/_mokk/requests/:id`
