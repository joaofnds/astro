Keep your habits, like an astronaut! 🧑‍🚀

![Screen-Recording-2022-11-05-at-6](https://user-images.githubusercontent.com/9938253/200142489-a1cb6bfb-6d68-4f48-9366-46b48cef26e1.gif)

backend at [joaofnds/gastro](https://github.com/joaofnds/gastro)

# Install

```sh
# install
brew install joaofnds/tap/astro
# start the app
astro
```

Astro is also on aur as [astro-bin](https://aur.archlinux.org/packages/astro-bin)

## Token

The token identifies you as a user, and your habits are tied to your token. If you wish to have your habits shared across
multiple machines, just copy `~/.config/astro/token` to each machine.

Astro will create a token for you on launch, but if you wish do it manually, here's how:

```sh
curl -X POST https://astro.joaofnds.com/token > ~/.config/astro/token
```
