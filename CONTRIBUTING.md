# Contributing

Contributions are welcome. All contributed code will be covered by the Apache License v2 of this project.

Below are the instructions. We understand that there is much left to be desired, and if you see any room for improvement, please let us know. Thank you.

## Linting

tron-docker CI uses [pre-commit](https://pre-commit.com/) to lint all code within the repo. Add it to your local
copy following the [installation](https://pre-commit.com/#installation).

This repo uses a squash-and-merge workflow to avoid extra merge commits. After forking it, create an `upstream` remote
with `git remote add upstream git@github.com:tronprotocol/tron-docker.git`, and create a git alias with
`git config --global alias.push-clean '!git fetch upstream main && git rebase upstream/main && git push -f'`. You can
then `git push-clean` to your fork before opening a PR.
