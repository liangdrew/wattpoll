# poll-service

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
mysql> USE DATABASE poll_service;
```

Run `db/sql/setup.sql` to set up your tables.

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
    "question": "Who is the coolest?",      //required
    "storyId": "123456",                    //optional but good to have
    "partId": "7890",                       //required
    "choices": [                            //required, must have between 2-4 choices
        {
            "choice": "a"                   
        },
        {
            "choice": "b"
        },
        {
            "choice": "c"
        },
        {
            "choice": "d"
        }
    ]
}' "localhost:8081/polls/create"
```

#### GET request to /polls/get to retrieve a poll

`curl "localhost:8081/polls/get?partId=7890username=clover"`

`partId` is required.\
`username` is optional - only include if the user is logged in

This returns JSON

```
{ 
    "question": "Who is the coolest?",      
    "totalVotes": 1000,
    "userVote": 1,                      // If the user is logged in and has voted, ID of their voted choice is returned
                                        // Otherwise, 0 is returned
    "created": "2017-04-27T17:00:35Z",
    "choices": [                        // Array with 2-4 elements
        {
            "id": 1,
            "choice": "a",
            "votes": 100
        },
        {
            "id": 2,
            "choice": "b",
            "votes": 200
        },
        {
            "id": 3,
            "choice": "c",
            "votes": 300
        },
        {
            "id": 4,
            "choice": "d",
            "votes": 400
        }
    ]
}

```

#### POST request to /polls/vote to a vote on a poll

```
curl -X POST -d '{
    "storyId": "123456",            // Optional but good to have
    "partId": "7890",               // Required
    "choiceId": 1,                  // Required, it's id of the selected choice
    "username": "namenamename"      // Optional - only include if user is logged in
}' "localhost:8081/polls/vote"
```
