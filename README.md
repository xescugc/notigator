# Notigator

This library provides a way to aggregate notifications in one sigel view

## How to install

It uses GO +1.11

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

It'll start a server at [`localhost:3000`](localhost:3000) by default, to get more information use

```
$> notigator -h
```

## Configure

You can configure `notigator` with config file. The config file should look like this (by default on `$HOME/.notigator`):

```json
{
  "sources": [
    {
      "name": "Display name",
      "token": "Token",
      "canonical": "github"
    }
  ]
}
```

To have more documentation on the supported Sources and how to configure them read the wiki
