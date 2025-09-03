# Create User
POST http://{{host}}/api/users
{
    "firstName": "TestUserLift",
    "lastName": "LiftNotSoMuch",
    "email": "I.dont@lifts.com",
    "password": "Weakling"
}
HTTP 200
[Captures]
email: jsonpath "$['email']"
uid: jsonpath "$['id']"
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.firstName" == "TestUserLift"
jsonpath "$.lastName" == "LiftNotSoMuch"
jsonpath "$.email" == "I.dont@lifts.com"
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists

# Login as User
POST http://{{host}}/api/login
{
    "email": "{{email}}",
    "password": "Weakling"
}
HTTP 200
[Captures]
jwt: jsonpath "$['token']"
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.id" == {{uid}}
jsonpath "$.firstName" == "TestUserLift"
jsonpath "$.lastName" == "LiftNotSoMuch"
jsonpath "$.email" == {{email}}
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists
jsonpath "$.token" exists
jsonpath "$.refresh_token" exists

# Incorrect PW at Login
POST http://{{host}}/api/login
{
    "email": "{{email}}",
    "password": "Wrong password"
}
HTTP 401
Content-Type: text/plain; charset=utf-8


#Update Existing user
POST http://{{host}}/api/users/{{uid}}
Authorization: Bearer {{jwt}}
{
    "firstName": "LiftTestUser",
    "lastName": "NeverLifts",
    "email": "I.never@lifts.com",
    "password": "Weakling123"
}
HTTP 200
[Captures]
fname: jsonpath "$['firstName']"
lname: jsonpath "$['lastName']"
nemail: jsonpath "$['email']"
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.id" == {{uid}}
jsonpath "$.firstName" != "TestUserLift"
jsonpath "$.lastName" != "LiftNotSoMuch"
jsonpath "$.email" != "{{email}}"
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists
jsonpath "$.token" exists

# Incorrect token Update Existing user
POST http://{{host}}/api/users/{{uid}}
Authorization: Bearer 123456
{
    "firstName": "LiftTestUser",
    "lastName": "NeverLifts",
    "email": "I.never@lifts.com",
    "password": "Weakling123"
}
HTTP 401
Content-Type: text/plain; charset=utf-8

# Incorrect Token Delete Existing user
DELETE http://{{host}}/api/users/{{uid}}
Authorization: Bearer wrongToken1
HTTP 401
Content-Type: text/plain; charset=utf-8

# Delete Existing user
DELETE http://{{host}}/api/users/{{uid}}
Authorization: Bearer {{jwt}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.id" == {{uid}}
jsonpath "$.firstName" == "{{fname}}"
jsonpath "$.lastName" == "{{lname}}"
jsonpath "$.email" == "{{nemail}}"
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists
jsonpath "$.token" exists
