POST http://localhost:8080/api/users
{
    "firstName": "Arnold",
    "lastName": "Schwarzenegger",
    "email": "Arnold.Schwarzenegger@lifts.com",
    "password": "Lift Heavy"
}
HTTP 500
# [Asserts]
# header "Content-Type" == "application/json"
# jsonpath "$.firstName" == "Arnold"
# jsonpath "$.lastName" == "Schwarzenegger"
# jsonpath "$.email" == "Arnold.Schwarzenegger@lifts.com"
# jsonpath "$.createdAt" exists
# jsonpath "$.modifiedAt" exists


POST http://localhost:8080/api/login
{
    "email": "Arnold.Schwarzenegger@lifts.com",
    "password": "Lift Heavy",
    "expires_in_seconds": 1800
}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.id" exists
jsonpath "$.firstName" == "Arnold"
jsonpath "$.lastName" == "Schwarzenegger"
jsonpath "$.email" == "Arnold.Schwarzenegger@lifts.com"
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists
jsonpath "$.token" exists
