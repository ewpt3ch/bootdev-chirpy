# bootdev-chirpy# Web server written in Go for [Boot.dev: Learn HTTP servers in Go](https://www.boot.dev/courses/learn-http-servers-golang) course.

Demonstrates:
- Understanding and architecture of a RESTful api using go.
- Understanding of web server patterns using the go standard libraries net/http package.
- Usage of a database, postgres, and support tooling, sqlc to generate code from sql queries and goose for managing schema migrations.
- Usage of internal packages to seperate concerns
- Unit testing
- I used [hurl](https://hurl.dev/) as an endpoint testing client
- Usage of JWT and other tokens for access.
- Common endpoints with different responses depending on query information and http method in the request.
- Error handling and responses to the client
- Making sure secrets stay secret.


### Config

the http root will be at static/ in the root of the project dir.
