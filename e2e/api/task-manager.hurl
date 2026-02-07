GET {{ host }}/api/v1/livez
Content-Type: application/json
HTTP 200
[Asserts]
jsonpath "$.status" == "OK"

GET {{ host }}/api/v1/tasks/scheduled
Content-Type: application/json
HTTP 200
[Asserts]
jsonpath "$.type" == "scheduled_task"
jsonpath "$.data" count == 1

POST {{ host }}/api/v1/tasks
Content-Type: application/json
{
    "task_id": "schedule-task",
    "scheduled_time": "* * * * *",
    "args": {
        "room_id": "1"
    }
}
HTTP 202
[Asserts]
jsonpath "$.type" == "scheduled_task"
jsonpath "$.data" == "" 

DELETE {{ host }}/api/v1/tasks/019c3909-8ac5-760b-a99b-45cdfb16ce50
Content-Type: application/json
HTTP 204

POST {{ host }}/api/v1/tasks
Content-Type: application/json
{
    "task_id": "schedule-task",
    "scheduled_time": "wrong-cron",
    "args": {
        "room_id": "1"
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

