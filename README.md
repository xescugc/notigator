# Notigator

This library provides a way to aggregate notifications in one sigle view

![notigator](https://github.com/xescugc/notigator/blob/master/docs/screenshot.png)

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

## Docker

If you want to use it from docker this is the "default" example:

```
$> docker run --rm -p 3000:3000 -v "${HOME}/.notigator.json":/app/.notigator.json xescugc/notigator serve --config /app/.notigator.json
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

**Note:** All the flags can also be passed via config file or envs.

To have more documentation on the supported Sources and how to configure them read the [wiki](https://github.com/xescugc/notigator/wiki/Config-File)
