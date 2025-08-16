# Create User
POST http://{{host}}/api/users
{
    "firstName": "UserTestExercises",
    "lastName": "Log exercise",
    "email": "LogThe@exercise.com",
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
jsonpath "$.firstName" == "UserTestExercises"
jsonpath "$.lastName" == "Log exercise"
jsonpath "$.email" == "LogThe@exercise.com"
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists

# Login as User
POST http://{{host}}/api/login
{
    "email": "{{email}}",
    "password": "LiftTheWeight",
    "expires_in_seconds": 150
}
HTTP 200
[Captures]
jwt: jsonpath "$['token']"
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.id" == {{uid}}
jsonpath "$.firstName" == "{{fname}}"
jsonpath "$.lastName" == "{{lName}}"
jsonpath "$.email" == {{email}}
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists
jsonpath "$.token" exists

# Create user exercise record
POST http://{{host}}/api/record-exercise
Authorization: Bearer {{jwt}}
{
   "userID": "{{uid}}",
   "exerciseID": "4b7f40a2-ff00-4e3c-b5e4-01fc338cb36f",
   "sets_weight": {
       "set1": {"12": 135},
        "set2": {"12": 135},
        "set3": {"9": 135}
    },
    "rest": 60,
    "notes": "First time back and extremely exhausting"
}
HTTP 201
[Captures]
rID: jsonpath "$['recordID']"
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.recordID" exists
jsonpath "$.userID" == "{{uid}}"
jsonpath "$.exerciseID" == "4b7f40a2-ff00-4e3c-b5e4-01fc338cb36f"
jsonpath "$.sets_weight" exists
jsonpath "$.rest" == 60
jsonpath "$.notes" == "First time back and extremely exhausting"
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists

# Get the user exercise record
GET http://{{host}}/api/record-exercise/{{rID}}
Authorization: Bearer {{jwt}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.recordID" == {{rID}}
jsonpath "$.userID" == "{{uid}}"
jsonpath "$.exerciseID" == "4b7f40a2-ff00-4e3c-b5e4-01fc338cb36f"
jsonpath "$.sets_weight" exists
jsonpath "$.rest" == 60
jsonpath "$.notes" == "First time back and extremely exhausting"
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists


# Update the user exercise record
POST http://{{host}}/api/record-exercise/{{rID}}
Authorization: Bearer {{jwt}}
{
   "userID": "{{uid}}",
   "exerciseID": "4b7f40a2-ff00-4e3c-b5e4-01fc338cb36f",
   "sets_weight": {
       "set1": {"12": 135},
        "set2": {"12": 135},
        "set3": {"12": 135}
    },
    "rest": 60,
    "notes": "First time back and extremely exhausting, but I did it"
}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.recordID" == {{rID}}
jsonpath "$.userID" == "{{uid}}"
jsonpath "$.exerciseID" == "4b7f40a2-ff00-4e3c-b5e4-01fc338cb36f"
jsonpath "$.sets_weight" exists
jsonpath "$.rest" == 60
jsonpath "$.notes" == "First time back and extremely exhausting, but I did it"
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists


DELETE http://{{host}}/api/record-exercise/{{rID}}
Authorization: Bearer {{jwt}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.recordID" == {{rID}}
jsonpath "$.userID" == "{{uid}}"
jsonpath "$.exerciseID" == "4b7f40a2-ff00-4e3c-b5e4-01fc338cb36f"
jsonpath "$.sets_weight" exists
jsonpath "$.rest" == 60
jsonpath "$.notes" == "First time back and extremely exhausting, but I did it"
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists

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
