Keep you habits, like an astronaut! 🧑‍🚀

![ezgif-1-7eb8bed3e9](https://user-images.githubusercontent.com/9938253/193378739-b96de1c2-3106-41ff-aaf2-f02b594bf22f.gif)

backend at [joaofnds/gastro](https://github.com/joaofnds/gastro)

# Install

```sh
# install
brew install joaofnds/tap/astro
# start the app
astro
```

## Token
The token identifies you as a user, and your habits are tied to your token. If you wish to have your habits shared across
multiple machines, just copy `~/.config/astro/token` to each machine.

Astro will create a token for you on launch, but if you wish do it manually, here's how:

```sh
curl -X POST https://astro.joaofnds.com/token > ~/.config/astro/token
```