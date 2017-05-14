# Dementor

## Overview

Dementor is the CLI tool of [Azkaban](https://github.com/azkaban/azkaban).

- Has Wrapper commands of official APIs
- Supports a YAML recipe file for registering and scheduling projects/flows

## Usage

### Basic wrapper commands

Dementor has wrapper commands of official APIs such as __Create a Project__, __Delete a Project__... (Not all at the moment)

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

Dementor has the kiss command which process a recipe written in YAML.

A recipe looks like below:

```
url: "http://localhost:8081/"
insecure: false
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

## License

[MIT](https://opensource.org/licenses/MIT)