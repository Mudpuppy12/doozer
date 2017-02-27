# doozer
Restful API task execution using echo and machinery in golang

This is a POC that will hopefully expand into something more useful. Goal is to create a resetful api that can trigger defined tasks to execute. The execution will issue an ID in which api calls can querry the status of the job/task.

POC will be based off echo framework and the machinery framework. It is also a learning exprience in switching from python to Go.

# Features
* Simple API to some basic tasks. Using the examples (Add and Multiply)
* Querry on task completion/status
*  Simple Authentication and issuing a JWT.
* Workers to execute tasks and return results (Add and Multiply)



