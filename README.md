# Minecraft Overseer

An application easing configuration and operation of Minecraft in a modern environment.

# Deployment Scenarios

## Docker
`docker run --rm -p25565:25565 -v $PWD/config:/mc/config:ro meschbach/minecraft-overseer:latest`

This will bootstrap and start Minecraft.  The file `config/manifest.json` must exist.  Please see `test/config/manifest.json`
for an example for running Minecraft 1.17.1.  Use the option `-v $PWD/data:/mc/instance` to mount a place to persist
world data.

## Kubernetes

There is a helm chart in `deployment/minecraft-overseer`.  This will create the manifest for your instance and by
default mount an `emptyDir.medium: "Memory"` instance for the game.  See `deployment/minecraft-overseer/values.yaml` to
salt and pepper to taste.
