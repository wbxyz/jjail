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

3. Integrate the `jjail` skill into your AI agent's workspace by copying the `.agents` directory from this repository into the root of your target repository:
   ```bash
   cp -r path/to/jjail/.agents path/to/your/repo/
   ```
   This will make the `jjail` skill available to the agent. The canonical instructions for the agent are maintained in [.agents/skills/jjail/SKILL.md](.agents/skills/jjail/SKILL.md).

4. Run your AI agent.

## Usage

```bash
jjail <command> [args...]
```


## Security Guarantee

By restricting operations to the `agent::` subtree, `jjail` ensures that an agent cannot:
- Accidentally rebase internal commits onto public branches.
- Modify historical commits that are not part of its current task.
- Access or manipulate sensitive parts of the repository history outside its assigned scope.
