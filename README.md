# jjail

`jjail` is a utility designed to provide VCS guardrails by wrapping the [Jujutsu (jj)](https://github.com/jj-vcs/jj) CLI for use by AI agents. It provides a sandboxed environment that restricts an agent's ability to modify the repository history outside of a designated subtree.

## How it Works

To use `jjail`, you place the `agent` bookmark in your jj graph somewhere. `jjail` will then only allow changes within that subtree.

Any operation that involves a revision is validated to ensure that the revision falls within the subtree defined by the `agent` bookmark.

The validation logic uses the Jujutsu revset expression:
`(<target_rev>) ~ (agent::)`

If this expression returns any commits, it means the target revision is outside the allowed subtree, and the operation is blocked.

## Installation and Setup

1. Install `jjail` using `go install`:
   ```bash
   go install github.com/wbxyz/jjail/cmd/jjail@latest
   ```
   Ensure that your `$(go env GOPATH)/bin` directory is in your system's `PATH`.

2. Place the `agent` bookmark in your repository where you want the agent to start working:
   ```bash
   jj bookmark create agent -r @
   ```

3. Update your `AGENTS.md` or `GEMINI.md` file with the [AI Agent Instructions](#ai-agent-instructions-agentsmd--geminimd) below.

4. Run your AI agent.

## Usage

```bash
jjail <command> [args...]
```

## AI Agent Instructions (AGENTS.md / GEMINI.md)

Copy and paste the following block into your repository's `AGENTS.md` or `GEMINI.md` file to instruct AI agents on how to safely interact with Jujutsu using `jjail`:

```markdown
# Repository Management (Jujutsu via jjail)

This repository uses Jujutsu (jj) for version control. However, as an AI agent, you MUST NOT use the `jj` or `git` CLI directly to modify the repository's history. 

Instead, you MUST use the `jjail` CLI tool, which is a sandboxed wrapper around `jj` ensuring you stay within your designated working bounds.

### Rules for using jjail:
1. **Never use `jj` or `git` directly for modifications.** Always use `jjail` to create commits, update descriptions, squash, rebase, or abandon changes.
2. **Your sandbox is the `agent` bookmark.** All your work will be rooted from the commits in the `agent::` subtree. `jjail` enforces this boundary. If `jjail` throws a "Sandbox violation!" error, it means you are trying to operate on a commit outside your allowed scope.
3. **Interactive Commands.** Do not use `jjail split` without specific filesets, as it will attempt to open an interactive terminal and crash. Use `jjail new`, copy partial file contents, and `jjail squash` to manually split changes if needed.

### Allowed Commands:
- `log` / `list`: View the allowed subtree (`agent::`).
- `status` / `st`: Show the working copy status.
- `diff [rev]`: Show changes in a revision (defaults to `@`).
- `show [rev]`: Show commit message and changes in a revision (defaults to `@`).
- `new [base_rev]`: Create a new change on top of `base_rev` (defaults to `@`).
- `edit <rev>`: Set a revision as the working copy.
- `describe <rev> <msg>`: Update the description of a change.
- `rebase <src> <dest>`: Rebase a change within the subtree.
- `squash <src> [into_rev]`: Squash changes.
- `split <rev> [fileset]`: Split a change (Note: providing no fileset triggers interactive mode which will fail).
- `duplicate <rev>`: Duplicate a change.
- `abandon <rev>`: Abandon a change.

### Common Workflow:
- Check your workspace: `jjail log` and `jjail status`
- See your current changes: `jjail diff`
- Finalize the current working copy: `jjail describe @ "Your commit message"`
- Start a new change: `jjail new`
- Drop a bad change: `jjail abandon <rev>`
```

## Security Guarantee

By restricting operations to the `agent::` subtree, `jjail` ensures that an agent cannot:
- Accidentally rebase internal commits onto public branches.
- Modify historical commits that are not part of its current task.
- Access or manipulate sensitive parts of the repository history outside its assigned scope.
