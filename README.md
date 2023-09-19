# Sheduler-booking demo backend

## How to use

```
docker-compose up --build
```

# API

### GET /units

Returns all neccessary information to build booking dataset. Using `slots + usedslots` approach

#### Response example

```json
{
  "id": 2,
  "title": "Larissa Tillery",
  "category": "Allergist",
  "subtitle": "",
  "details": "",
  "preview": "",
  "price": 120,
  "slots": [
    {
      "from": 9, // 9:00
      "to": 14, // 14:00
      "size": 45,
      "gap": 5,
      "days": [1, 3, 5] // reccuring events
    },
    {
      "from": 15.5, // 15:30
      "to": 20, // 20:00
      "size": 45,
      "gap": 5,
      "dates": [1695254400000] // Thu Sep 21 2023
    }
  ]
  ...
}
```

### GET /doctors

Returns a list of doctors

#### Response example

```json
[
  {
    "id": 1,
    "name": "Conrad Hubbard",
    "subtitle": "2 years of experience",
    "details": "Desert Springs Hospital (Schroeders Avenue 90, Fannett, Ethiopia)",
    "category": "Psychiatrist",
    "price": 45,
    "gap": 20,
    "slot_size": 20,
    "image_url": "https://files.webix.com/30d/d34de82e0a8e3b561988a46ce1e86743/stock-photo-doc.jpg"
  },
  ...
]
```

### GET /doctors/worktime

Returns a list of doctor's schedule with concrete dates (excluding recurring slots and expired dates).
You can show this data on Doctors view in Booking-Scheduler Demo

#### Response exapmle

```json
[
  {
    "id": 2,
    "doctor_id": 1,
    "start_date": "2023-09-20T08:00:00+03:00",
    "end_date": "2023-09-20T16:00:00+03:00"
  },
  {
    "id": 3,
    "doctor_id": 1,
    "start_date": "2023-09-21T08:00:00+03:00",
    "end_date": "2023-09-21T16:00:00+03:00"
  },
  {
    "id": 4,
    "doctor_id": 1,
    "start_date": "2023-09-22T08:00:00+03:00",
    "end_date": "2023-09-22T16:00:00+03:00"
  }
  ...
]
```

### POST /doctors/worktime

Creates a new doctor's schedule with concrete date (Doctors view)

#### Body

```json
{
  "doctor_id": 1,
  "end_date": "2023-09-24 10:30",
  "start_date": "2023-09-24 14:30"
}
```

### Response example

Returns an id of created schedule (Doctors view)

```json
{
  "id": 10
}
```

### PUT /doctors/worktime/{id}

Updates doctor's schedule

#### Body

```json
{
  "doctor_id": 1,
  "end_date": "2023-09-23 12:20",
  "start_date": "2023-09-23 16:55"
}
```

### Response example

Returns an id of updated schedule (Doctors view)

```json
{
  "id": 10
}
```

#### URL Params:

- id [required] - id of the schedule to be updated

### DELETE /doctors/worktime/{id}

Deletes doctor's schedule (Doctors view)

#### URL Params:

- id [required] - id of the schedule to be deleted

### GET /doctors/reservations

Returns all occupied slots (Clients view)

#### Response example

```json
[
    {
        "id": 1,
        "doctor_id": 2,
        "date": 1695367800000,
        "client_name": "Alan",
        "client_email": "alan@gmail.com",
        "client_details": ""
    },
    {
        "id": 2,
        "doctor_id": 3,
        "date": 1695275400000,
        "client_name": "Viron",
        "client_email": "viron@hr.com",
        "client_details": ""
    }
    ...
]
```

### POST /doctors/reservations

Creates reservation (Booking view)

#### Body

```json
{
  "doctor": 2,
  "date": 1695278400000,
  "form": {
    "name": "Alan",
    "email": "alan@gmail.com",
    "details": ""
  }
}
```

### DELETE /doctors/reservations/{id}

Deletes reservation

#### URL Params:

- id [required] - id of the reservation to be deleted

# Config

```yaml
db:
  path: db.sqlite # path to the database
  resetonstart: true # reset data on server restart
server:
  url: "http://localhost:3000"
  port: ":3000"
  cors:
    - "*"
  resetFrequence: 120 # every 2 hours restart data (value in minutes)
```
