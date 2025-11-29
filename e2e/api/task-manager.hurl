GET {{ host }}/api/v1/healthz
Content-Type: application/json
HTTP 200

GET {{ host }}/api/v1/tasks
Content-Type: application/json
HTTP 200
[Asserts]
jsonpath "$.type" == "task"
jsonpath "$.data" count == 1

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
