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
    "exerciseBlocks": [
        {
            "ordinal": 0,
            "exercises": [
                {
                    "ordinal": 0,
                    "exerciseID": "1880fef4-61bc-4497-bcbc-edddc73a744b",
                    "sets": [
                        {
                            "ordinal": 0,
                            "weight": 165,
                            "reps": 12
                        },
                        {
                            "ordinal": 1,
                            "weight": 165,
                            "reps": 12
                        },
                        {
                            "ordinal": 2,
                            "weight": 165,
                            "reps": 12
                        }
                    ],
                    "notes": "hold for a beat at check before pressing away"
                }
            ],
            "restAfterBlock": 60
        }
    ]
}
HTTP 200
[Captures]
wid: jsonpath "$['id']"
wname: jsonpath "$['name']"
wdesc: jsonpath "$['description']"
firstBlockID: jsonpath "$['exerciseBlocks'][0].workoutBlockID"
firstExerciseID: jsonpath "$['exerciseBlocks'][0]['exercises'][0].workoutExerciseID"
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.id" exists
jsonpath "$.name" == "TestWorkout"
jsonpath "$.description" == "This is a test workout"
jsonpath "$.exerciseBlocks[0].exercises" exists
jsonpath "$.exerciseBlocks[0].exercises[0].sets" exists
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists


# Get Workout by ID
GET http://{{host}}/api/workouts/{{wid}}
Authorization: Bearer {{jwt}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.id" == "{{wid}}"
jsonpath "$.name" == "{{wname}}"
jsonpath "$.description" == "{{wdesc}}"
jsonpath "$.exerciseBlocks[0].workoutBlockID" == "{{firstBlockID}}"
jsonpath "$.exerciseBlocks[0].exercises[0].workoutExerciseID" == "{{firstExerciseID}}"

# Delete Workout By ID
DELETE http://{{host}}/api/workouts/{{wid}}
Authorization: Bearer {{jwt}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.id" == "{{wid}}"
jsonpath "$.name" == "{{wname}}"
jsonpath "$.description" == "{{wdesc}}"
jsonpath "$.exerciseBlocks[0].workoutBlockID" == "{{firstBlockID}}"
jsonpath "$.exerciseBlocks[0].exercises[0].workoutExerciseID" == "{{firstExerciseID}}"

# Get Deleted Workout
GET http://{{host}}/api/workouts/{{wid}}
Authorization: Bearer {{jwt}}
HTTP 500
Content-Type: text/plain; charset=utf-8

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
