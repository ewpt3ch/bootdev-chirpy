POST http://localhost:8080/admin/reset
HTTP 200

POST http://localhost:8080/api/users
Content-Type: application/json
{
  "email": "eric@ewpt3ch.dev",
  "password": "password123456"
}
HTTP 201

POST http://localhost:8080/api/login
Content-Type: application/json
{
  "email": "eric@ewpt3ch.dev",
  "password": "password123456"
}
HTTP 200
[Captures]
jwt_token: jsonpath "$.token"
refresh_token: jsonpath "$.refresh_token"

POST http://localhost:8080/api/chirps
Authorization: Bearer {{jwt_token}}
Content-Type: application/json
{
  "body": "chirp, this chirp, oh boy"
}
HTTP 201

POST http://localhost:8080/api/refresh
Authorization: Bearer {{refresh_token}}
HTTP 200

POST http://localhost:8080/api/revoke
Authorization: Bearer {{refresh_token}}
HTTP 204

POST http://localhost:8080/api/refresh
Authorization: Bearer {{refresh_token}}
HTTP 401
