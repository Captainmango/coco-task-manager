GET {{ host }}/api/v1/livez
Content-Type: application/json
HTTP 200
[Asserts]
jsonpath "$.status" == "OK"

GET {{ host }}/api/v1/tasks/scheduled
Content-Type: application/json
HTTP 200
[Asserts]
jsonpath "$.type" == "task"
jsonpath "$.data" count == 0

POST {{ host }}/api/v1/tasks
Content-Type: application/json
{
    "task_id": 1,
    "scheduled_time": "* * * * *",
    "args": {
        "room_id": 1
    }
}
HTTP 201
""

POST {{ host }}/api/v1/tasks
Content-Type: application/json
{
    "task_id": 1,
    "scheduled_time": "wrong-cron",
    "args": {
        "room_id": 1
    }
}
HTTP 422
{
    "error": "InvalidCronExpression",
    "message": "Cron expression provided is malformed"
}

POST {{ host }}/api/v1/tasks
Content-Type: application/json
{
    "task_id": 9999999,
    "scheduled_time": "1 1 1 1 1",
    "args": {
        "room_id": 1
    }
}
HTTP 404
{
    "error": "EntityNotFound",
    "message": "Entity not found"
}

