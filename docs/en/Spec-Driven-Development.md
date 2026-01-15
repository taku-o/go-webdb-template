**[日本語](../ja/Spec-Driven-Development.md) | [English]**

# cc-sdd

## Overview
cc-sdd is a tool for spec-driven development.
It works on Cursor and Claude Code.

## Commands
```
## Project Analysis
## (Run at the beginning or when the project structure changes.)
/kiro:steering

## Create project requirements, design, and tasks (batch)
/kiro/spec-init "{project requirements}"

## Create project requirements
/kiro:spec-requirements "{project requirements}"

## Create design
/kiro:spec-design

## Create task list
/kiro:spec-tasks

## Execute specified task
/kiro:spec-impl "{task}"
```

## Usage Instructions

### Working on Cursor (for higher document quality)
```
## Run once when project is created
/kiro:steering

## Create project requirements
/kiro:spec-requirements "Create requirements for https://github.com/taku-o/go-webdb-template/issues/3. GitHub CLI is installed."
think.

## If a requirements document creation plan is displayed
Please create the requirements document.

## After reviewing the requirements document content
Approve the requirements document.

## Create design
/kiro:spec-design

## After reviewing the design document content
Approve the design document.

## Create task list
/kiro:spec-tasks

## After reviewing the task list content
Approve the task list.

Please create a git branch for this requirement.
Please commit the work so far.
```

### Working on Claude Code (for faster work speed)

```
## Start implementation
/kiro:spec-impl 0003-gorm-introduction

## Create Pull Request after implementation is done
Please commit the changes so far.
Then, create a pull request for https://github.com/taku-o/go-webdb-template/issues/57

## Review the created Pull Request
/review 61
```
