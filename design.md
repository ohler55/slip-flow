# SLIP-Flow Design

SLIP-Flow is a process flow package for use with SLIP which is SLIce
Processing for golang, a mostly Common LISP implementation.

## Concepts

At the highest level, a Flow is a collection of Tasks that for a
processing unit. Data enters a Flow and transitions from one Task to
another until processing is complete.

SLIP-Flow is implemented primarlify in golang with an API that is
primarily Flavors based but with corresponding functions to to also be
a CLOS API.

The classes or flavors in the package are:

 * flow-manager
 * flow
 * task
 * actor
 * box
 * track

## Classes (Flavors)

Elements of the packages are implemented as classes or more
specifically as Flavors with additional CLOS style method functions in
addition to the Flavors methods.

### flow-manager

An instance of the **flow-manager-flavor** is used to load or create
instances of the **flow-flavor** . The flows (**flow-flavor**
instances) are typically initialized with a configuration what is
loaded from either a JSON or LISP file although a flow can be build
fully using LISP code by creating tasks explicitly. After loading
flows managed by the flow manager can be accessed by flow name.

The flow manager includes a logger that is shared with all flows and
tasks to make monitoring more consistent. Another shared resource in
the flow manager is an error handler that acts as the backstop for
errors not handled by the flows themselves.

### flow

Every flow must have a unique name in the flow manager. Flows are
primarily a container for tasks although there are some elements of
the flow that are shared across tasks. One shared element is the log
level and an error handler if different than flow manager. If set the
flow error handle overrides the flow manager error handler for that
specific flow.

A flow also has an option entry task. If data can be submitted
directly to a flow then the entry task must be set. If the entry task
is not set a trigger task such as an HTTP server task or periodic time
task is needed.

Flows are typically build from a configuration JSON that adheres to a
specific format. A flow can also be created directly from code by
adding tasks to a flow and then linking up the tasks.

### task

Tasks represent steps in a process. They are a container for an actor
which does the actuall processing for the task. Tasks can be either
synchronous or make use of one or more workers that pull data from a
work queue. The task is handles transitions through links to other
tasks based on the response from an actor when the actor is asked to
perform.

The configuration for a task includes what kind of actor to create. An
actor can be either a single LISP function or lambda or else a flavor
name. If a function then on receipt of a box of data it is either
passed to the function or placed on a queue and a worker routine calls
the function. If the actor is a flavor name then an instance of that
flavor is created with a configuration from the task configuration and
the `:perform` method is called on the actor when a box is received or
popped off a work queue. The response from the actor must include the
transition (link name) to follow and a new box of data. The task then
sends the new box on the link identifed by the transition specified.

All calls to the actor are wrapped in a recover so any panic is
captured and sent on an "error" link if it exists or to the error
handler of the flow.

Shared resources other then an instance of the **flow-box-flavor** for
a task with more than one worker must be concurrent safe.

### actor

A actor can be a function or an instance of a Flavor.

If a function then it must expect exactly one argument and return a
list of transition name and a box.

If an instance then it is expected to have at least the following methods:

 - _:perform_ (box) that return a transition name and new box.
 - _:init_ with a keyword of _:configuration_.
 - _:set-task_ to let the instance know what task it is contained in.
 - _:shutdown_ to cleanup any open resources and to stop processing.

Actors should use the task for logging but generally don't need to
access the task.

If an error occurs then a panic is expected since it will be caught by
the task which will pick the appropriate transition.

### link (not a Flavor)

Links exist but are not visible outside the internals of the
package. They include a name and a target task.

### box

An instance of the **flow-box-flavor** is referred to as simply a
box. A box contains tracking information and data content. The
**flow-box-flavor** is similar to the SLIP **bag-flavor** except then
content is immutable and any call to modify the content make a copy of
the data first.

### track

The tracking information in a box is also a copy on set attribute of
the box. It contains a tracking identifier and a history of the tasks
and times the box has transitioned through.

## Reasons

It may not be obvious why the design follows the patterns described
without thinking through the various scenarios and use cases that
could be encountered. This section covers of a few of the less obvious
use cases and how they influenced the design.

### Flow Manager

By using a top level manager for all flows a shared logging and error
handler can be employed. It also allows for flows to call transitions
to other flows by using a link than includes not only a task name but
also a flow name.

A flow manager also allows flow validation that across all flows if
nested flows are used.

### Triggers

To allow flows to be used for both manual invocation as well as
triggered invocation both are supported. A manual or nested call makes
used of the flow entry task while others can uses a trigger task to
start a flow. Triggered task examples are an HTTP server, and periodic
timer, or a channel listener. All of those type of tasks are
implemented by using built in actors in the package.

### Concurrency

It is often the case that flows include a task that takes time to
complete. This might be a database call or an HTTP request. By
supporting concurrency in the flow multiple jobs can be processed in
parallel. This package supports concurrency by allowing separate go
routines for each task as well as multiple go routines within a task.

### Immutable Data

By supporting concurrency in processing it is important to keep data
from being modified concurrently. That can be done with the use of a
mutex but this package instead uses immutable data containers referred
to as boxes. Data used in one task is made immutable on exiting the
task and can only be modified by first making a copy of that
data. This approach allows the actors in the tasks to not have to
block when modifying data.

### Actors

The first thought that might come to mind is that tasks should be
subclassed to support different behavior. That presents an issue then
trying to use multiple workers against a queue. It also would mean the
implementor would have to deal with concurrency and transition
selection. Using an actor isolates the implementor in almost all cases
which reduces the chance for the implementor to shoot themselves in
the foot.

### Tracking

Tracking data as it passes through a flow is key to monitoring and
gathering metrics. A unique tracking identifier also allow for a
processing path to split and rejoin after each processing path has
completed. In addition to having a unique tracking identifier each
track includes the history of the tasks the job passed through along
with the time it entered a task. One issue with this flattened
approach is that parallel tracks are merged and flattened. This can be
rectified in the analysis by knowing how tasks are linked so if it
becomes important the history can be overlaid on the flows.

### Nested flows

As a processing flow gets more complicated it is advantageous to be
able to break the flow into sub-flows or nested flows. This allows for
more understandable flows and for reuse of common flows. Nested flows
are supported by the use of a flow manager and a sub-flow actor.

### Parallel Paths

Being able to process multiple longer running processing paths
concurrently speeds the processing where there are paths that can be
executed in parallel. The issue with many process flow systems is how
to merge the results of the parallel paths. SLIP-Flow makes use of the
tracking identifier to match up merging paths. A data merge function
is then employed to merge the box data into a single box. A timeout is
employed to deal with failed processing paths.

### Workers and Queues

When needed a task can set up one or more worker routines that pull
from a common queue (golang chan). By using instances to perform the
processing worker actors can be cleaned up using the `:shutdown`
method when the flow is shutdown. It also allows the processing to be
distributed across multiple go routines.

### Errors

With a common configurable error handler errors can all be handled in
the same way. Any error is also tagged with the flow and task name
that raised the error. Actors need not be concerned with how errors
will be handled. They simply panic and let the flow package take care
of the rest. Trigger tasks are the exception in that they need to
handle panics themselves and log accordingly.

### Logging

Logging makes use of the SLIP gi:Logger which is an async logger that
supports concurrent logging. All components of the package use the
same logger and actors are expected to also use the same logger by
calling the logging methods of the Task that contains them.

Logging is hierarchical in that each log entry includes the flow and
task name and can be controlled by the logging level in the flow
manager, flow, and task. Log levels are error, warn, info, and debug.
