# Create User
POST http://{{host}}/api/users
{
    "firstName": "UserTestWorkouts",
    "lastName": "WorkoutMcgee",
    "email": "LogThe@workout.com",
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
jsonpath "$.firstName" == "UserTestWorkouts"
jsonpath "$.lastName" == "WorkoutMcgee"
jsonpath "$.email" == "LogThe@workout.com"
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

# Create Workout
POST http://{{host}}/api/workouts
Authorization: Bearer {{jwt}}
{
    "name": "TestWorkout",
    "description": "This is a test workout",
    "exercises": [
        {
            "exerciseID":  "4b7f40a2-ff00-4e3c-b5e4-01fc338cb36f",
            "sets": [
                {"weight": 165, "reps": 12},
                {"weight": 165, "reps": 12},
                {"weight": 165, "reps": 12}
            ]
        }
    ]
}
HTTP 200
[Captures]
wid: jsonpath "$['id']"
wname: jsonpath "$['name']"
wdesc: jsonpath "$['description']"
wercersises: jsonpath "$['exercises']"
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.id" exists
jsonpath "$.name" == "TestWorkout"
jsonpath "$.description" == "This is a test workout"
jsonpath "$.exercises" exists
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists

#Get workout by ID
GET http://{{host}}/api/workouts/{{wid}}
Authorization: Bearer {{jwt}}
HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.id" == "{{wid}}"
jsonpath "$.name" == "{{wname}}"
jsonpath "$.description" == "{{wdesc}}"
jsonpath "$.exercises" == {{wercersises}}
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists

# Get all workouts
GET http://{{host}}/api/workouts
Authorization: Bearer {{jwt}}
HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.workouts" exists


# Delete Workout by Id
DELETE http://{{host}}/api/workouts/{{wid}}
Authorization: Bearer {{jwt}}
HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.id" == "{{wid}}"
jsonpath "$.name" == "{{wname}}"
jsonpath "$.description" == "{{wdesc}}"
jsonpath "$.exercises" == {{wercersises}}
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
