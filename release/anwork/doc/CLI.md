This documentation is generated from com.marshmallow.anwork.app.cli.GithubReadmeDocumentationGenerator

## _anwork_ ... : ANWORK CLI commands
* `[-d|--debug]` : Turn on debug printing
* `[-c|--context <name:STRING>]` : Set the persistence context
* `[-o|--output <directory:STRING>]` : Set persistence output directory
* `[-n|--no-persist]` : Do not persist any task information
### `anwork reset [-f|--force]`
* Completely delete all ANWORK stuff related to this persistence context
* `[-f|--force]` : Force the persistence context to be deleted, i.e., don't prompt the user for approval
### `anwork summary <days:NUMBER>`
* Show a summary of the past days of work
* `<days:NUMBER>` : Number of days to look back for work

## _anwork journal_ ... : Journal commands...
### `anwork journal show <task-specifier:STRING>`
* Show the entries in the journal for a task
* `<task-specifier:STRING>` : The task-specifier for the task(s) for which to show journal entries
### `anwork journal show-all`
* Show all of the entries in the journal

## _anwork task_ ... : Task commands...
### `anwork task create [-e|--description <description:STRING>] [-p|--priority <priority:NUMBER>] <task-name:STRING>`
* Create a task
* `[-e|--description <description:STRING>]` : The description of the task
* `[-p|--priority <priority:NUMBER>]` : The priority of the task
* `<task-name:STRING>` : The name of the task to create
### `anwork task delete <task-specifier:STRING>`
* Delete a task
* `<task-specifier:STRING>` : The task-specifier for the task(s) to delete
### `anwork task delete-all`
* Delete all tasks
### `anwork task note <task-specifier:STRING> <note:STRING>`
* Add a note to a task
* `<task-specifier:STRING>` : The task-specifier for the task(s) to which to add a note
* `<note:STRING>` : The note to add to a task
### `anwork task set-blocked <task-specifier:STRING>`
* Set a task as blocked
* `<task-specifier:STRING>` : The task-specifier for the task(s) to set blocked
### `anwork task set-finished <task-specifier:STRING>`
* Set a task as finished
* `<task-specifier:STRING>` : The task-specifier for the task(s) to set finished
### `anwork task set-priority <task-specifier:STRING> <priority:NUMBER>`
* Set the priority of a task
* `<task-specifier:STRING>` : The task-specifier for the task(s) on which to set the priority
* `<priority:NUMBER>` : The priority to set on the task
### `anwork task set-running <task-specifier:STRING>`
* Set a task as running
* `<task-specifier:STRING>` : The task-specifier for the task(s) to set running
### `anwork task set-waiting <task-specifier:STRING>`
* Set a task as waiting
* `<task-specifier:STRING>` : The task-specifier for the task(s) to set waiting
### `anwork task show [-s|--short]`
* Show all tasks
* `[-s|--short]` : Show a shorter description of all of the tasks
