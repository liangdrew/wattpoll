# poll-service

A Go microservice which supports the integration of polls in story parts within the Wattpad app. Created for Wattpad's April 2017 Hackathon.

## Setup

### Step 1: Set up Go

```
$ brew update
$ brew install go
$ mkdir -p $HOME/go/src
$ export GOPATH=$HOME/go
```

### Step 2: Clone this repo
```
$ cd $GOPATH
$ git clone https://github.com/liangdrew/poll-service
```

### Step 3: Download MySQL

From: https://dev.mysql.com/downloads/mysql/

Remember to save the temporary password given to you at the end of the download.

### Step 4: Set up MySQL

After download completes, run `mysql -u root -p` to connect to the server.\
You'll be prompted to enter the temporary password from step 3.

Change your password by running `ALTER USER 'root'@'localhost' IDENTIFIED BY 'NEW_PASSWORD';`

Create a database named `poll_service`.

```
mysql> CREATE DATABASE poll_service;
mysql> USE poll_service;
```

Run [setup.sql](https://github.com/liangdrew/poll-service/blob/master/db/sql/setup.sql) in the MySQL shell to set up your tables.

### Step 5: Run the service locally

```
$ cd $GOPATH/poll-service
$ go run main.go
```

Now you're ready to send requests to the service!

## Sample requests

#### POST request to /polls/create to create a new poll

```
curl -X POST -d '{
    "question": "Who would you like to see in the next story part?",
    "storyId": "107474356",
    "partId": "404628388",
    "choices": [
        {
            "choice": "Harry Styles"
        },
        {
            "choice": "Rich Poirier"
        },
        {
            "choice": "Zayn Malik"
        },
        {
            "choice": "Justin Bieber"
        }
    ]
}' "localhost:8081/polls/create"
```

#### GET request to /polls/get to retrieve a poll

`curl "localhost:8081/polls/get?partId=404628388&username=clover"`

`partId` is required.\
`username` is optional - only include if the user is logged in

This returns JSON

```
{ 
    "question": "Who would you like to see in the next story part?",      
    "totalVotes": 3,
    "userVote": 2,                      // If the user is logged in and has voted, ID of their voted choice is returned
                                        // Otherwise, 0 is returned
    "created": "2017-04-27T17:00:35Z",
    "choices": [                        // Array with 2-4 elements
        {
            "id": 1,
            "choice": "Harry Styles",
            "votes": 0
        },
        {
            "id": 2,
            "choice": "Rich Poirier",
            "votes": 2
        },
        {
            "id": 3,
            "choice": "Zayn Malik",
            "votes": 1,
        },
        {
            "id": 4,
            "choice": "Justin Bieber",
            "votes": 0
        }
    ]
}

```

#### POST request to /polls/vote to a vote on a poll

```
curl -X POST -d '{
    "storyId": "107474356",              // Optional, but good to have
    "partId": "404628388",               // Required
    "choice": "Rich Poirier",            // Required
    "choiceId": 2,                       // Required, it is the ID of the selected choice
    "username": "clover"                 // Optional - only include if user is logged in
}' "localhost:8081/polls/vote"
```
