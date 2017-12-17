# Dementor

## Overview

Dementor is the CLI tool for [Azkaban](https://github.com/azkaban/azkaban).

- Has wrapper commands of [official APIs](http://azkaban.github.io/azkaban/docs/latest/#ajax-api)
- Supports a YAML recipe file for scheduling projects and flows

## Usage

### Basic wrapper commands

Dementor has wrapper commands of official APIs such as __Create a Project__, __Delete a Project__... (not all APIs at the moment)

You can see all commands with **--help** option.

**Note**: You don't need to call **authenticate** command to fetch a session ID. It is called internally in each command.

```
$ dementor --help
Usage: dementor [--version] [--help] <command> [<args>]

Available commands are:
    authenticate         Fetch a session id by username and password.
    createProject        Create a new project.
    deleteProject        Delete the project.
    fetchFlowsProject    Fetch flows of the project
    fetchJobsFlow        Fetch jobs of the flow.
    fetchSchedule        Fetch the schedule of the flow.
    kiss                 Process a recipe file
    scheduleFlow         Schedule the flow.
    unscheduleFlow       Unschedule the flow.
    uploadProjectZip     Upload the project zip file.

```

### Kiss command (processing a recipe written in YAML)

You can write a recipe which defines projects, flows, and schedules in YAML file.

Dementor has the kiss command for processing it.

```
$ dementor kiss -f path_to_recipe_file.yml
```

A recipe file looks as below:

```
url: "http://localhost:8081/"
username: azkaban
password: azkaban
projects:
  - name: RecipeTest
    description: "This is a recipe test project."
    filepath: "./assets/test.zip"
    flows:
      - name: test
        cron: "0 23/30 5,7-10 ? * 6#3"
```

**Note**: **cron** must be a [Quartz Cron Format](http://www.quartz-scheduler.org/documentation/quartz-2.x/tutorials/crontrigger.html).

If the same project is already registered, Dementor removes it first and then processes a recipe. In other words, **kiss** command has idempotency.

## License

[MIT](https://opensource.org/licenses/MIT)