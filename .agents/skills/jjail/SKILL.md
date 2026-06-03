---
name: jjail
description: Manage the Jujutsu (jj) repository history using the sandboxed jjail CLI tool. Use when you need to commit, rebase, squash, or manage changes.
---

# Jujutsu Repository Management via jjail

This skill provides instructions for using the `jjail` CLI tool to manage version control in this repository.

## When to Use

Use this skill whenever you need to perform version control operations, such as:
- Checking repository status (`jjail status`)
- Viewing history (`jjail log`)
- Creating new changes (`jjail new`)
- Committing or updating descriptions (`jjail describe`)
- Rebasing, squashing, or splitting changes.

## Rules and Constraints

1. **Never use `jj` or `git` directly for modifications.** Always use `jjail` to modify the repository's history.
2. **Your sandbox is the `agent` bookmark.** All modifications must be rooted from the commits in the `agent::` subtree. `jjail` enforces this boundary for mutating commands. If you get a "Sandbox violation!" error, you are attempting to modify a commit outside your allowed scope. Note: Read-only commands (`log`, `status`, `diff`, `show`) are allowed to run outside this sandbox.
3. **No Interactive Commands.** Do not use `jjail split` without specific filesets, as it will attempt to open an interactive terminal and crash. Use `jjail new`, copy partial file contents, and `jjail squash` to manually split changes if needed.

## Allowed Commands

The `jjail` tool supports a limited set of safe commands:
- `log` / `list`: View the repository log (can be run on any revisions, even outside the jail).
- `status` / `st`: Show the working copy status (can be run on any working copy).
- `diff [rev] [files...]`: Show changes in a revision (can be run on any revision, even outside the jail).
- `show [rev] [files...]`: Show commit message and changes in a revision (can be run on any revision, even outside the jail).
- `new [base_rev]`: Create a new change on top of `base_rev` (defaults to `@`).
- `edit <rev>`: Set a revision as the working copy.
- `describe <rev> <msg>`: Update the description of a change. Keep msg brief and don't use conventional commit message prefixes.
- `rebase <src> <dest>`: Rebase a change within the subtree.
- `squash <src> [into_rev]`: Squash changes.
- `split <rev> [fileset]`: Split a change.
- `duplicate <rev>`: Duplicate a change.
- `abandon <rev>`: Abandon a change.

## Common Workflows

### Checking Status
```bash
jjail log
jjail status
```

### Committing Changes
```bash
jjail describe @ "Added my new feature"
```

### Starting a New Change
```bash
jjail new
```

### Abandoning a Change
```bash
jjail abandon <rev>
```
