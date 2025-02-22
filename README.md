# projectctl

projectctl is a CLI tool to interact with GitHub ProjectV2 boards provided by Giant Swarm.

## Features

- List projects for an owner (user or organization).
- List items for a given project.
- List fields for a given project.

## Prerequisites

- Go installed on your system.
- A valid GitHub token with appropriate scopes. Ensure the environment variable GITHUB_TOKEN is set. The token requires the 'read:project' scope for proper functioning.

## Installation & Build

1. Clone the repository.
2. Set your GitHub token:
   ```bash
   export GITHUB_TOKEN=your_token_here
   ```
3. Build the project:
   ```bash
   go build -o projectctl
   ```

## Usage

projectctl provides several commands:

### List Projects

Lists GitHub ProjectV2 boards for a given owner:

```bash
./projectctl projects --owner <owner> [--owner-type user|organization] [--output default|table|json|yaml]
```

### List Items

Lists items for a specified project:

```bash
./projectctl items --project-id <project_id> [--output default|table|json|yaml]
```

### List Fields

Lists all fields for a specified project:

```bash
./projectctl fields --project-id <project_id> [--output default|table|json|yaml]
```

## License

This project is licensed under the terms of the LICENSE file.
