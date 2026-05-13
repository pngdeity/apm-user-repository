---
name: project-management
description: Integrate with project management software when completing tasks. Use when creating issues, updating task status, or interacting with project trackers like GitHub Projects or TODO.md.
allowed-tools: bash, gh
metadata:
  tags: "project-management issues tracking github-projects todo"
compatibility: Generic — requires gh CLI for GitHub Projects integration. TODO.md works with any environment.
---
When completing tasks, integrate with existing project management software.

Prompt the user if any details about the project management process is unclear.

This skill should eventually be provider agnostic. Or another way to phrase it, this is the interface specification for project management processes, not how to interface with project management software. Since this is a less than alpha draft, there are implementation details for now.

Some possible implementation targets:

TODO.md 
Already widely used with a proportional amount of variance in formatting, content, and use. Its implementation and interaction is almost entirely obvious.

GitHub Projects
This is available through a REST API or the gh command-line tool
https://cli.github.com/manual/gh_project

Mozilla bugzilla
Not sure if this can be interfaced with outside of the web browser, but it would be nice to get blockers in mozilla repos without having to leave the terminal.
