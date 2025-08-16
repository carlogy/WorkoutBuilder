POST http://{{host}}/api/login
{
    "email": "Arnold.Schwarzenegger@lifts.com",
    "password": "Lift Heavy",
    "expires_in_seconds": 1800
}
HTTP 200
[Captures]
jwt: jsonpath "$['token']"

POST http://{{host}}/api/exercises
Authorization: Bearer {{jwt}}
{
        "name": "Barbell Benchpress",
        "exerciseType" : 1,
        "equipment": "Barbell",
        "primaryMuscleGroups": {
            "Pectoralis Major":"Chest",
            "Pectoralis Minor":"Chest"
            },
        "secondaryMuscleGroups" : {
                "Anterior Deltoid": "Deltoids",
                "Triceps Brachii": "Arm Extensors"
                },
        "description": "Lie faceup on a horizontal bench, with glutes on the bench and feet flat on the ground. Grasp the barbell in an overhand grip with, slightly larger than shoulder width apart. Keeping feet flat, and glutes flexed, core braced bring the barbell down to your chest and then extend your arms lifting the barbell back to starting position."
    }
HTTP 201
Content-Type: application/json

GET http://localhost:8080/api/exercises
Authorization: Bearer {{jwt}}
HTTP 200
Content-Type: application/json
