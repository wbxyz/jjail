# jjail

`jjail` is a security-focused utility designed to wrap the [Jujutsu (jj)](https://github.com/martinvonz/jj) CLI for use by AI agents. It provides a sandboxed environment that restricts an agent's ability to modify the repository history outside of a designated subtree.

## Goal

The primary goal of `jjail` is to allow AI agents to interact with a Jujutsu repository while guaranteeing that they cannot modify, rebase, or otherwise interfere with commits outside of their authorized "jail."

## How it Works

`jjail` enforces a sandbox boundary using a specific Jujutsu bookmark (default: `gemini`). Any operation that involves a revision is validated to ensure that the revision falls within the subtree defined by that bookmark.

The validation logic uses the Jujutsu revset expression:
`(<target_rev>) ~ (<bookmark>::)`

If this expression returns any commits, it means the target revision is outside the allowed subtree, and the operation is blocked.

## Allowed Commands

`jjail` only exposes a subset of `jj` commands that are safe and relevant for agent workflows:

- `log` / `list`: View the allowed subtree (`gemini::`).
- `new [base_rev]`: Create a new change on top of `base_rev` (defaults to `@`).
- `describe <rev> <msg>`: Update the description of a change.
- `rebase <src> <dest>`: Rebase a change within the subtree.
- `squash <src> [into_rev]`: Squash changes.
- `split <rev>`: Split a change into two.
- `duplicate <rev>`: Duplicate a change.

## Usage

```bash
jjail <command> [args...]
```

## Security Guarantee

By restricting operations to the `gemini::` subtree, `jjail` ensures that an agent cannot:
- Accidentally rebase internal commits onto public branches.
- Modify historical commits that are not part of its current task.
- Access or manipulate sensitive parts of the repository history outside its assigned scope.
