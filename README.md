> [!IMPORTANT]
> This is a work in progress. The plugin is not yet fully functional and is still being developed.
> It is mainly intended for my own use and testing, but I hope it will be useful to others as well.

# About
Extended Docker Stack CLI plugin

```
Usage:  docker stackx COMMAND

Extended Docker Stack CLI plugin

Commands:
  config      Outputs the final config file, after doing merges and interpolations
  deploy      Deploy a new stack or update an existing stack

Run 'docker stackx COMMAND --help' for more information on a command.
```

## Example

The following environment variables will always available when using the command provided by this plugin:
- `DOCKER_REGISTRY_URL`: The Docker registry to use for the images (default: `docker.io`).
- `DOCKER_STACK_NAMESPACE`: The namespace to use for the stack (default: `default`).
- `RANDOM`: A random number to use for the stack.

```yml
services:
  nginx:
    image: ${DOCKER_REGISTRY_URL}/nginx
    environment:
      - DOCKER_STACK_NAMESPACE=${DOCKER_STACK_NAMESPACE}
      - RANDOM=${RANDOM}
```
