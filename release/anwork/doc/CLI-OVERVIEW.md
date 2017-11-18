# ANWORK CLI OVERVIEW

The ANWORK CLI API is fully documented in [CLI.md](CLI.md). Below are some tips and tricks.

## Creating a task

To create a task named _weigh-tuna_, here is the command.
```
$ anwork task create weigh-tuna
```

To create a task named count-fish and give it a description and a priority, here is the command.
```
$ anwork task create count-fish -d 'This task is to count the fish that have been caught today.' -p 15
```

## Setting a task's state

There is a command for setting a task to each of 4 possible states: _waiting_, _running_, _blocked_,
and _finished_. By default, tasks start out in the _waiting_ state.
```
$ anwork task set-waiting weigh-tuna
$ anwork task set-running weigh-tuna
$ anwork task set-blocked weigh-tuna
$ anwork task set-finished weigh-tuna
```

## Adding a note to a task

To add a note to a task, here is the command.
```
$ anwork task note weigh-tuna 'I tried weighing the tuna, but the scale is currently broken'
```

## Showing the status of all tasks

To print out the status of all current tasks, here is the command.
```
$ anwork task show
```

To print out a shortened status of all current tasks, here is the command.
```
$ anwork task show -s
```

## Showing the journal of events

To show the complete journal of events, here is the command.
```
$ anwork journal show-all
```

To show a journal of events for a single task, here is the command.
```
$ anwork journal show
```

## Deleting a task

To delete a task, here is the command.
```
$ anwork task delete weigh-tuna
```

## Using task specifiers

A task specifier starts with the '@' symbol and can refer to one or more tasks. It can be passed to
the CLI commands that take an argument named _task-specifier_. Here are some examples.
```
$ anwork task set-waiting @1 # set the task with ID 1 to the waiting state
$ anwork task note @42 'Here is a note' # add a note to the task with ID 42
``` 

## Setting a persistence context

A persistence context is simple an ID used to specify a single instance of ANWORK tasks. For
example, users may choose to have one context for a bunch of tasks that need to be done at home, and
a separate context for a bunch of tasks that need to be done at work. The persistence context can be
set with a combination of the _-c_ flag and the _-o_ flag. The _-o_ flag specifies the persistence
directory, while the _-c_ flag specifies the context ID. Here is an example
```
$ anwork -p home-context -o ~/.anwork create wash-dishes
$ anwork -p work-context -o ~/.anwork create put-new-cover-sheet-on-tps-reports
``` 