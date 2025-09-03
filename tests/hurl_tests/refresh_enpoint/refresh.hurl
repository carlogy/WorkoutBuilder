# Create User
POST http://{{host}}/api/users
{
    "firstName": "UserRefresh",
    "lastName": "TestToken",
    "email": "RefreshThe@token.com",
    "password": "LiftTheWeight"
}
HTTP 200
[Captures]
email: jsonpath "$['email']"
uid: jsonpath "$['id']"
fname: jsonpath "$['firstName']"
lName: jsonpath "$['lastName']"
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.firstName" == "UserRefresh"
jsonpath "$.lastName" == "TestToken"
jsonpath "$.email" == "RefreshThe@token.com"
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists

# Login as User
POST http://{{host}}/api/login
{
    "email": "{{email}}",
    "password": "LiftTheWeight"
}
HTTP 200
[Captures]
jwt: jsonpath "$['token']"
refresh_token: jsonpath "$['refresh_token']"
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.id" == {{uid}}
jsonpath "$.firstName" == "{{fname}}"
jsonpath "$.lastName" == "{{lName}}"
jsonpath "$.email" == {{email}}
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists
jsonpath "$.token" exists
jsonpath "$.refresh_token" exists

# Refresh the Refresh Token
POST http://{{host}}/api/refresh
Authorization: Bearer {{refresh_token}}
HTTP 200
[Captures]
newRToken: jsonpath "$['token']"
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.token" != {{refresh_token}}

# Revoke the  Refresh Token
POST http://{{host}}/api/revoke
Authorization: Bearer {{newRToken}}
HTTP 204

# Refresh a revoked token
POST http://{{host}}/api/refresh
Authorization: Bearer {{newRToken}}
HTTP 401

# Delete Existing user
DELETE http://{{host}}/api/users/{{uid}}
Authorization: Bearer {{jwt}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.id" == {{uid}}
jsonpath "$.firstName" == "{{fname}}"
jsonpath "$.lastName" == "{{lName}}"
jsonpath "$.email" == "{{email}}"
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists
jsonpath "$.token" exists
