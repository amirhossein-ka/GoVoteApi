# GoVoteApi

This is a api to get votes from people.

## **routes**

Base path: `/api/v1/`.

Please note that all endpoints need a trailing slash `/`.
- __USER_ENDPOINTS__:
  - **Register**
    - Path: `/user/register/`
    - Method: `POST`
    - Request body:
    ```json
    {
    "username":"mrBanana",
    "fullname":"mooze Mooz",
    "email":"MoziEmail@gmail.com",
    "password":"bananaPass"
    }
    ```
    - Response:
    ```json
    {
    "status": "created",
    "id": 1,
    "token": "jwt_token"
    }
    ```
  - **Login**
    - Path: `/user/login/`
    - Method: `POST`
    - Request body:
    ```json
    {
    "username":"mrBanana",
    "password":"bananaPass"
    }
    ```
    - Response:
    ```json
    {
    "status": "found",
    "id": 1,
    "token": "jwt_token"
    }
    ```
  - **Info**
    - Authorized
    - Path: `/user/`
    - Method: `GET`
    - Response:
    ```json
    {
	"status": "found",
	"data": {
		"id": 1,
		"fullname": "mooze Mooz",
		"username": "mrBanana",
		"email": "MoziEmail@gmail.com",
		"role": 1
	  }
    }
    ```
  - **Delete**
    - Authorized 
    - Path: `/user/`
    - Method: `DELETE`
    - Response:
    ```json
    {
	  "status": "deleted",
	  "data": "user deleted"
    }
    ```
  - **UPDATE TOKEN**
    - Not Implemented yet
    
- __VOTE ENDPOINTS__

  Note: All of this section endpoints must be authorized
  
  - **Create a New Vote**
  - **Get All Votes**
  - **Get Vote by id**
  - **Get Vote with slug**
  - **Choose option(s) on a open Vote With ID**
  - **Choose option(s) on a open Vote With slug**
  

