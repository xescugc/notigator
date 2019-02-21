# Notigator

This library provides a way to aggregate notifications in one sigel view

## How to install

It uses GO 1.11>

```
$> go get github.com/xescugc/notigator
$> cd $GOPATH/src/github.com/xescugc/notigator
$> go install .
```

## How to use

Basically run 

```
$> notigator serve
```

It'll start a server at `localhost:3000` by default, to get more information use

```
$> notigator -h
```

## Configure

You can configure `notigator` using flags, envs or a config file, the config file. The config file should look like this (by default on `$HOME/.notigator`):

```json
{
  "github-token": "token",
  "gitlab-token": "token",
  "trello-api-key": "api-key",
  "trello-token":"token"
}
```

Right now the only Sources supported are Trello, Github and Gitlab, you do not need to configure all of them.
