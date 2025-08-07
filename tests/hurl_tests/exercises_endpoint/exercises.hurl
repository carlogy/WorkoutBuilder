POST http://localhost:8080/api/login
{
    "email": "Arnold.Schwarzenegger@lifts.com",
    "password": "Lift Heavy",
    "expires_in_seconds": 1800
}
HTTP 200
[Captures]
jwt: jsonpath "$['token']"


GET http://localhost:8080/api/exercises
Authorization: Bearer {{jwt}}
HTTP 200
Content-Type: application/json
