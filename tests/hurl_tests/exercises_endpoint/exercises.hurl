# Create User
POST http://{{host}}/api/users
{
    "firstName": "UserTestExercises",
    "lastName": "LiftMcgee",
    "email": "CreateThe@exercise.com",
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
jsonpath "$.lastName" == "LiftMcgee"
jsonpath "$.email" == "CreateThe@exercise.com"
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
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.id" == {{uid}}
jsonpath "$.firstName" == "{{fname}}"
jsonpath "$.lastName" == "{{lName}}"
jsonpath "$.email" == {{email}}
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists
jsonpath "$.token" exists


# Create exercise
POST http://{{host}}/api/exercises
Authorization: Bearer {{jwt}}
{
        "name": "Barbell Benchpress",
        "exerciseType" : "WeightedReps",
        "equipment": "Barbell",
        "primaryMuscleGroups": [
        {
            "bodyPart": "Chest",
            "muscleGroup": "Pectorals",
            "muscleName": "Pectoralis Major"
        }, {
            "bodyPart": "Chest",
            "muscleGroup": "Pectorals",
            "muscleName": "Pectoralis Minor"
        }
        ],
        "secondaryMuscleGroups" : [
            {
            "bodyPart": "Shoulders",
            "muscleGroup": "Deltoid",
            "muscleName": "Anterior Deltoid"
        }
        ],
        "description": "Lie faceup on a horizontal bench, with glutes on the bench and feet flat on the ground. Grasp the barbell in an overhand grip with, slightly larger than shoulder width apart. Keeping feet flat, and glutes flexed, core braced bring the barbell down to your chest and then extend your arms lifting the barbell back to starting position."
    }
HTTP 200
Content-Type: application/json
[Captures]
exID0: jsonpath "$['id']"
exDescription: jsonpath "$['description']"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.name" == "Barbell Benchpress"
jsonpath "$.equipment" == "Barbell"
jsonpath "$.exerciseType" == "WeightedReps"
jsonpath "$.primaryMuscleGroups" exists
jsonpath "$.secondaryMuscleGroups" exists
jsonpath "$.description" exists
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists


# Get Exercise by ID
GET http://{{host}}/api/exercises/{{exID0}}
Authorization: Bearer {{jwt}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{exID0}}
jsonpath "$.name" == "Barbell Benchpress"
jsonpath "$.exerciseType" == "WeightedReps"
jsonpath "$.primaryMuscleGroups" exists
jsonpath "$.secondaryMuscleGroups" exists
jsonpath "$.description" == {{exDescription}}
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists


# Create exercise
POST http://{{host}}/api/exercises
Authorization: Bearer {{jwt}}
{
    "name": "Incline Barbell Benchpress",
    "exerciseType": "WeightedReps",
    "equipment": "Barbell",
    "primaryMuscleGroups": [
      {
        "bodyPart": "Chest",
        "muscleGroup": "Pectorals",
        "muscleName": "Pectoralis Major"
      },
      {
        "bodyPart": "Chest",
        "muscleGroup": "Pectorals",
        "muscleName": "Pectoralis Minor"
      }
    ],
    "secondaryMuscleGroups": [
      {
        "bodyPart": "Arm Extensors",
        "muscleGroup": "Triceps",
        "muscleName": "Triceps Brachii"
      },
      {
        "bodyPart": "Deltoids",
        "muscleGroup": "Deltoid",
        "muscleName": "Anterior Deltoid"
      },
      {
        "bodyPart": "Abdominals",
        "muscleGroup": "Scapulae Fixers",
        "muscleName": "Serratus Anterior"
      }
    ],
    "description": "Sit on an incline bench angled between 45 and 60 degrees. With an overhand grip hold the barbell wider than shoulder width. Lower the barbell to the top of your sternum. Then press the bar away from the sternum by extending your arms."
  }
HTTP 200
Content-Type: application/json
[Captures]
exID1: jsonpath "$['id']"
exDescription: jsonpath "$['description']"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.name" == "Incline Barbell Benchpress"
jsonpath "$.exerciseType" == "WeightedReps"
jsonpath "$.equipment" == "Barbell"
jsonpath "$.primaryMuscleGroups" exists
jsonpath "$.secondaryMuscleGroups" exists
jsonpath "$.description" == {{exDescription}}
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists


# Get Exercise by ID
GET http://{{host}}/api/exercises/{{exID1}}
Authorization: Bearer {{jwt}}
HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.id" == {{exID1}}
jsonpath "$.name" == "Incline Barbell Benchpress"
jsonpath "$.exerciseType" == "WeightedReps"
jsonpath "$.equipment" == "Barbell"
jsonpath "$.primaryMuscleGroups" exists
jsonpath "$.secondaryMuscleGroups" exists
jsonpath "$.description" == {{exDescription}}
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists


# Create exercise
POST http://{{host}}/api/exercises
Authorization: Bearer {{jwt}}
{
    "name": "Decline Barbell Benchpress",
    "exerciseType": "WeightedReps",
    "equipment": "Barbell",
    "primaryMuscleGroups": [
      {
        "bodyPart": "Chest",
        "muscleGroup": "Pectorals",
        "muscleName": "Pectoralis Major"
      }
    ],
    "secondaryMuscleGroups": [
      {
        "bodyPart": "Arm Extensors",
        "muscleGroup": "Triceps",
        "muscleName": "Triceps Brachii"
      }
    ],
    "description": "Lie on a decline bench with an angle between 20 and 40 degress. Grasp the barbell in an overhand grip wider or at shoulder width apart. Lower the barbell to the lower part of your chest. Then, extend your arms driving the bar away from your chest."
  }
HTTP 200
Content-Type: application/json
[Captures]
exID2: jsonpath "$['id']"
exDescription: jsonpath "$['description']"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.name" == "Decline Barbell Benchpress"
jsonpath "$.exerciseType" == "WeightedReps"
jsonpath "$.equipment" == "Barbell"
jsonpath "$.primaryMuscleGroups" exists
jsonpath "$.secondaryMuscleGroups" exists
jsonpath "$.description" exists
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists


# Get Exercise by ID
GET http://{{host}}/api/exercises/{{exID2}}
Authorization: Bearer {{jwt}}
HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.id" == {{exID2}}
jsonpath "$.name" == "Decline Barbell Benchpress"
jsonpath "$.exerciseType" == "WeightedReps"
jsonpath "$.equipment" == "Barbell"
jsonpath "$.primaryMuscleGroups" exists
jsonpath "$.secondaryMuscleGroups" exists
jsonpath "$.description" == {{exDescription}}
jsonpath "$.createdAt" exists
jsonpath "$.modifiedAt" exists


# Get all the exercises
GET http://{{host}}/api/exercises
Authorization: Bearer {{jwt}}
HTTP 200
Content-Type: application/json
[Asserts]
# jsonpath "$.total" == 3


# Delete Exercise by ID
DELETE http://{{host}}/api/exercises/{{exID0}}
Authorization: Bearer {{jwt}}
HTTP 200


# Delete Exercise by ID
DELETE http://{{host}}/api/exercises/{{exID1}}
Authorization: Bearer {{jwt}}
HTTP 200


# Delete Exercise by ID
DELETE http://{{host}}/api/exercises/{{exID2}}
Authorization: Bearer {{jwt}}
HTTP 200

# Get Deleted Workout
GET http://{{host}}/api/exercises/{{exID2}}
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
