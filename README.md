# doozer
Restful API task execution using echo and machinery in golang

This is a POC that will hopefully expand into something more useful. The goal is to create a restful api that can trigger defined tasks to execute and given to workers. The execution will return a result or issue a task ID if the executing takes longer than 5 seconds. You can then use the task_id to query
the status and retrieve the results from the client.

POC will be based off echo framework and the machinery framework. It is also a learning experience in switching from python to Go.

# Features

* Simple REST API to execute some basic tasks. Using the examples (Add and Multiply from machinery)
* A Simple Client that executes the REST Api Calls.
* Simple Authentication and issuing a JWT.
* Query on task completion/status on task ID.
* Workers to execute tasks and return results (Add and Multiply)
* Server and workers use a broker back end for queueing (Part of machinery library)

# Code Status

* Adding is working 100%
* Multiply is WIP.
