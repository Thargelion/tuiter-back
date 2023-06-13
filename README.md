# tuiter-back

[![Go Coverage](https://github.com/thargelion/tuiter-back/wiki/coverage.svg)](https://raw.githack.com/wiki/thargelion/tuiter-back/coverage.html)

## Description

Tuiter Back is the Backend for the Twitter clone made in Android.

## How to run

There is a makefile that will help you to run the project. You can run the project with the following command:

```shell
make local.up
```

This will run the required mysql database.

To execute the server, run:

```shell
go build -o out/tuiter-back ./cmd/tuiter/main.go && ./out/tuiter-back
```

## Scaffolding

The project is built loosely based on Package Oriented Design. As such, the API domains are featured as packages at the
root of the project. It also separates inbound and outbound ports: the inbound ports are the API endpoints, and the
outbound ports are located at MYSQL package.

There is a KIT package that handles common use interfaces and utilities.

Finally, CMD package handles the entrypoint of the application.

## Domain Packages

### User

User defines the out interface and handles the HTTP endpoints. It contains the User Models.

### Post

Post defines the out interface and handles the HTTP endpoints. It contains the Post Models.

## Technologies

The project used CHI as HTTP Server, and GORM as ORM. It also uses a MySQL database.

## Next Steps

There are some improvements to make to this API like including proper middlewares to handle authentication, errors, and
metrics.

If needed, the User and Posts Router could be separated into their own packages.

Also, the IN implementation of both User and Post can be moved into the API package. But, the MYSQL implementation may
be moved inside both Domains packages. Both scaffoldings have their pros and cons. For example: having both IN and OUT
interface implementation inside the package simplifies the domain isolation. But, it also adds complexity to the
package. Otherwise, having the IN implementation inside the API package simplifies the package, but it also makes more
difficult to debug and evolve the domain packages. For example, if there is a problem regarding to Posts, it will be
required to look after multiple packages.

Finally, the overall coverage must be improved.

## SRE

The project lacks of proper load and stress tests. It should be achieved with a tool like Artillery or JMeter.

