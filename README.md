# Go Interview Challenge

Template interview challenge for the GPS Insight api team

This repo contains a little demo service for pulling data about intraday stock values off of a kafka stream and inserting them into a postgres database. The service also has a rest api for making that data consumable by clients.


## What is the point of this challenge?

- Demonstrate an ability to quickly get comfortable in a new codebase and make contributions
- Demonstrate familiarity with the technologies we use every day
- Show us your approach to solving problems


### Prerequisites

- [Go](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/)


## What will you be implementing?

- Consume kafka stream and write to postgres database
- Implement REST endpoint to get stock info from the database

![go-interview-challenge visualization](/overview.png "go-interview-challenge visualization")

If you are looking for ways to spice up your project here are some ideas...
- Try using gRPC or GraphQL instead of REST
- Show us your take on testing
- Allow for the data to be queried with filters/pagination/etc.


## How long should this challenge take?

You are welcome to put in as much time as you like but you should be able to finish it within an hour or so


## Getting Started

- Download and unpack zip of repository

<img src="/download.png" alt="download" width="400">

<p></p>

- Search the repo for `TODO:` tags (there should be two) and follow the instructions


### Running locally

To start local development environment with sample data and live reload of go-interview-challenge run
```
make run
```


## Generating data

To generate message on the kafka topic, use the command
```
make generate-messages
```
