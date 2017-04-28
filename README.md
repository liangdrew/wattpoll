# poll-service

A Go microservice which supports the integration of polls in story parts within the Wattpad app. Created for Wattpad's April 2017 Hackathon.

## Setup

### Step 1: Set up Go

```bash
$ brew update
$ brew install go
$ mkdir -p $HOME/go/src
$ export GOPATH=$HOME/go
```

### Step 2: Clone this repo
```bash
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

```bash
mysql> CREATE DATABASE poll_service;
mysql> USE poll_service;
```

Run [setup.sql](https://github.com/liangdrew/poll-service/blob/master/db/sql/setup.sql) in the MySQL shell to set up your tables.

### Step 5: Run the service locally

```bash
$ cd $GOPATH/poll-service
$ go run main.go
```

Now you're ready to send requests to the service!

## Sample requests

#### POST request to /polls/create to create a new poll

```bash
curl -X POST -d '{
    "question": "Who would you like to see in the next story part?",
    "storyId": "107474356",
    "partId": "404628388",
    "durationDays": 2,
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
   "question":"Who would you like to see in the next story part?",
   "totalVotes":1,
   "userVote":2,
   "created":"2017-04-28T15:14:51Z",
   "durationDays":2,
   "pollClosed":false,
   "choices":[  
      {  
         "id":1,
         "choice":"Harry Styles",
         "votes":0
      },
      {  
         "id":2,
         "choice":"Rich Poirier",
         "votes":0
      },
      {  
         "id":3,
         "choice":"Zayn Malik",
         "votes":0
      },
      {  
         "id":4,
         "choice":"Justin Bieber",
         "votes":0
      }
   ]
}

```

#### POST request to /polls/vote to a vote on a poll

```bash
curl -X POST -d '{
    "storyId": "107474356",            
    "partId": "404628388",              
    "choice": "Rich Poirier",          
    "choiceIndex": 3,                       
    "username": "clover"                 
}' "localhost:8081/polls/vote"
```
