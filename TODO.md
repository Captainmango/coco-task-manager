## Todo list

- [] Set up API
  - [x] Set up middleware
    - [x] Logging
    - [x] RequestID
    - [] ~CORS? (probably not... yet)~
  - [x] Accept requests on routes specified in e2e tests
  - [] DTOs?
  - [] database? (how do I store tasks that the API will reference?)
  - [] Write to crontab file
- [] Set up publisher
  - [] Set up RabbitMQ
    - [] Docker compose?
    - [] env vars?
  - [] exe to push to topic
    - [] Read from db or just encode in command in crontab?