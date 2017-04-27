# poll-service

## Setup

### Step 1: Clone this repo

You know the drill

### Step 2: Set up Go

```
brew update
brew install go
mkdir -p $HOME/go/src
export GOPATH=$HOME/go
```

### Step 3: Download mySQL

From: https://dev.mysql.com/downloads/mysql/

Make sure you select the Max OS X as your operating system.

Remember to save the password given to you at the end of the download.

### Step 4: set up mySQL

After download completes, run `mysql -u root -p` to start the server.\
You'll be prompted to enter the password from step 2.

You could change your password by running `ALTER USER 'root'@'localhost' IDENTIFIED BY 'NEW_PASSWORD';`

You can shut down the server using command `\q`



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

`curl "localhost:8081/polls/get?partId=7890username=namenamename"`

`partId` is required.\
`username` is optional - only include if the user is logged in

This returns a JSON

```
{ 
    "question": "Who is the coolest?",      
    "totalVotes": 1000,
    "userVote": 1,                     //only returned when username is passed in
    "created": "2017-04-27T17:00:35Z",
    "choices": [
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

#### POST request to /polls/vote to upload a vote

```
curl -X POST -d '{
    "storyId": "123456",            //optional but good to have
    "partId": "7890",               //required
    "choiceId": 1,                 //required, it's id of the selected choice
    "username": "namenamename"                  //optional - only include if user is logged in
}' "localhost:8081/polls/vote"
```
