# Contributing to goSpace
Welcome! Come join the fun and hack away.

## Naming
goSpace is the name used for this project.

`gospace` is used to talk about the source code attached to goSpace project.

## Source control
### Messages
In order to have a clean history of changes through the project and for identifying bugs, messages are important.

The following format is enforced on commit messages:

```
[project-or-submodule-name]: [short-summary-of-change(s)].

[extended-summary-or-crucial-information-to-commit]

[file-or-group-of-files-with-same-extension]: [change]

[extended-summary-or-crucial-information-to-change]

[file-or-group-of-files-with-same-extension]:
- [first-change]
- [second-change]
- [sublist-with-more-precise-changes]:
  * [first-subchange]
  * [second-subchange]
  * ...
- ...
```

Example for a commit message:
```
gospace: A lot of important changes are introduced.

*.go: Small and same changes in all source files.

{example, example_test}.go: Added Z.

Z will improve and refactor certain aspects of gospace/example.
This includes addition of A, B and C.

example.go:
   - Important feature X is added.
       * X is need because Y and improves on Z.

example_test.go:
   - Added test X.
   - Changed test A with Y.
   - Removed test B due to deprecation of feature Y.
```

Referencing to commits that introduced a certain feature, fix or regression is encouraged.

Avoid deviating from the message format as this makes the bug hunting game harder.

Through out the development of this project, expect the message format to change, as the need for better messages could occur.

### Control
It is crucial that changes in the development history are kept short and precise, unless:
- Features or interfaces are being added in which will need explanation.
- Features are being deprecated and removed in favor of better alternatives.
- Layout or structure of the project changes.

The general rules are as following:
- Partial commits are not welcome; branching is made exactly for that purpose.
- Squash commits locally through a rebase, in order too avoid to many small changes, trival fixes and errs in the commits.
- If multiple developers are working on a branch, communicate if rebasing or merging should/will occur and perform that action.
- If a single developer is working on a development branch and wants to merge into a stable branch, perform rebasing directly onto the stable branch.
- Make sure that a commit always has the needed information before adding it to branch.

## Contribution changes
This project is under development and so is this file. Make sure to keep up to date with regarding how to contribute to the project by reading updates to this file. 
