
# Dotfiles Sync

A cli tool for syncing a predefined list of dotfiles via a configuration
to a github repository




[![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](https://github.com/tterb/atomic-design-ui/blob/master/LICENSEs)
[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](https://opensource.org/licenses/)
[![AGPL License](https://img.shields.io/badge/license-AGPL-blue.svg)](http://www.gnu.org/licenses/agpl-3.0)


## Installation

Add the following to your `~/.zshrc` file
```
export PATH=${PATH}:`go env GOPATH`/bin
```
Install the binary with golang
```
go install github.com/cmdctl/dotfiles@latest
```
Or you can download the executables from the [releases page](https://github.com/cmdctl/dotfiles/releases/tag/v0.1.2)

## Usage
Create the following `.dotfiles.yml` at your `HOME` directory
```
touch ~/.dotfiles.yml
```
Add a list of dotfiles to sync with a reposiotory
```
version: "1.0"

include:
  - .vimrc
  - .zshrc
  - .dotfiles.yml
```
Then run in your terminal
```
dotfiles
```
This will create a new repository in **~/.dotfiles** directory and copy the specified files.

If you add a remote to the repository created by the `dotfiles` executable it will automatically push the changes on
each run.

> You can also add the binary to your `~/.zshrc` profile so that a sync is done on every new terminal session.
> Here is an example of my [~/.zshrc config](https://github.com/cmdctl/.dotfiles/blob/37d2f321c0be279d4d0fe79db8f4d5b4fc634120/.zshrc#L16) where the binary is run on startup.



## Contributing

Contributions are always welcome!




## Authors

- [@cmdctl](https://www.github.com/cmdctl)

