## ADDED Requirements

### Requirement: PreToolUse hook registration
The repository SHALL include a `.claude/settings.json` that registers a `PreToolUse` hook on the Bash tool.

#### Scenario: Hook registered for Bash tool
- **WHEN** Claude Code loads the project
- **THEN** it SHALL read `.claude/settings.json` and register the hook script as a `PreToolUse` handler for the `Bash` tool

#### Scenario: Hook script is executable
- **WHEN** the repository is cloned
- **THEN** `.claude/hooks/guard-risky-cmds.sh` SHALL have execute permissions (`chmod +x`)

### Requirement: Hard-block dangerous commands
The hook script SHALL hard-block dangerous command patterns by exiting with code 2, with no override mechanism.

#### Scenario: Credential exposure blocked
- **WHEN** Claude attempts to run a command reading SSH private keys (e.g., `cat ~/.ssh/id_*`), exporting GPG secret keys, or dumping all environment variables
- **THEN** the hook SHALL exit 2 with a message explaining what was blocked and why

#### Scenario: Pipe-to-shell blocked
- **WHEN** Claude attempts to run a command piping remote content directly to a shell (e.g., `curl | bash`, `wget | sh`)
- **THEN** the hook SHALL exit 2, as this is an unconditional supply-chain risk

#### Scenario: Sudo blocked
- **WHEN** Claude attempts any command prefixed with `sudo`
- **THEN** the hook SHALL exit 2; the user SHALL run sudo commands manually

#### Scenario: System file writes blocked
- **WHEN** Claude attempts to write to `/etc/` (e.g., `/etc/hosts`, `/etc/passwd`)
- **THEN** the hook SHALL exit 2

#### Scenario: Destructive git commands blocked
- **WHEN** Claude attempts `git push --force` or `git push -f` targeting `main` or `master`, `git reset --hard`, or `git clean -f`
- **THEN** the hook SHALL exit 2

#### Scenario: Destructive Docker commands blocked
- **WHEN** Claude attempts `docker system prune`, `docker volume rm`, or `docker rm -f`
- **THEN** the hook SHALL exit 2

#### Scenario: SSH tunnels and netcat listening blocked
- **WHEN** Claude attempts `ssh -L`, `ssh -R`, `ssh -D` (port forwarding), or `nc -l` (netcat listener)
- **THEN** the hook SHALL exit 2

### Requirement: Clear block messages
Every blocked command SHALL produce a human-readable message explaining what was blocked and how to proceed.

#### Scenario: Block message is informative
- **WHEN** the hook blocks a command
- **THEN** the stderr message SHALL name the matched pattern, state why it is dangerous, and instruct the user to run the command manually if intended

### Requirement: CLAUDE.md documents the policy
The repository SHALL include a `CLAUDE.md` at the root documenting the hook policy for developers using Claude Code.

#### Scenario: CLAUDE.md describes blocked categories
- **WHEN** a developer reads `CLAUDE.md`
- **THEN** they SHALL find a list of all blocked command categories with rationale and instructions for running blocked commands manually
