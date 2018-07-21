# WikiShell

## Description

This is a repository for an application used for browsing wikipedia using your shell written in Go.

## Installing/Building

To install/build the application run 

```shell
$ go build main.go
```

You should have a binary for the application in your `${gopath}/bin` directory

## Usage

You can pass an argument to the application to fetch an article on load. 

```shell
$ ./wikishell "Go (programming language)"
```

Inside the application you can use the following commands

### Load new article 
CTRL + L - Opens a modal window allowing you to navigate to a new article.

### Navigate to section
CTRL + S - Opens a modal window with a list of sections allowing you to jump to a section you are interested in.

## Contribute

PRs accepted.

## License

MIT © Nenad Lukic
