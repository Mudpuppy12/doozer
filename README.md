# doozer
Restful API task execution using echo and machinery (https://github.com/RichardKnop/machinery) in Go!

This is a POC that will hopefully expand into something more useful. The goal is to create a restful api that can trigger defined tasks to execute and given to workers. The execution will return a result or issue a task ID if the executing takes longer than 5 seconds. You can then use the task_id to query
the status and retrieve the results from the client.

POC will be based off echo framework and the machinery framework. It is a learning experience for me switching from python to Go. It is based off the examples from machinery, using send/task/worker for inspiration.

I am using gb (https://getgb.io/) for a build tool  and atom (https://atom.io/) Go extensions
for an editor.


# Features

* Simple REST API to execute some basic tasks. Using the examples (Add and Multiply from machinery)
* A Simple Client that executes the REST API Calls.
* Simple Authentication and issuing a JWT to a client.
* Query on task completion/status on task ID.
* Workers to execute tasks and return results (Add and Multiply)
* Server and workers use a broker back end (Using Redis) for queueing (Part of machinery library)

# Code Status

* Finished 100%


# Client Help

<pre>
$ ./client help
Simple client to interact with Dozer API service.

Usage:
  client
  client [command]

Available Commands:
  add         Add api call.
  help        Help about any command
  lookup      Lookup a task uuid.
  mul         Multiply api call.
  token       Print a JWT token.
  version     Print the version.

Use "client [command] --help" for more information about a command.
</pre>

# Client adding
<pre>
$ ./client add --i 1,2,3,4,5,6
Result: 21

</pre>

# Client adding when workers down

<pre>
$ ./client add --i 1,2,3
Defered! task_1db8fd1f-aff0-4db6-9c9a-ada3d20cb006

$ ./client lookup --uuid=task_1db8fd1f-aff0-4db6-9c9a-ada3d20cb006
Status : PENDING

</pre>

# When Workers come back online...

<pre>
./client lookup --uuid=task_1db8fd1f-aff0-4db6-9c9a-ada3d20cb006
Status : SUCCESS
Result : 6
</pre>
