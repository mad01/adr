[![Build Status](https://travis-ci.com/marouni/adr.svg?branch=master)](https://travis-ci.com/marouni/adr)

# ADR Go
A minimalist command line tool written in Go to work with [Architecture Decision Records](http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions) (ADRs).

Greatly inspired by the [adr-tools](https://github.com/npryce/adr-tools) with all the added benefits of using the Go instead of Bash.

# Quick start
## Installing adr
Go to the [releases page](https://github.com/marouni/adr/releases) and grab one of the binaries that corresponds to your platform.

Alternatively, if you have a Go development environment setup you can install it directly using :
```bash
go install github.com/mad01/adr@latest
```


## Initializing adr
```bash
adr init --readme README.md
```

## Creating a new ADR

As simple as :
```bash
adr new --title "google managed prometheus"
```
this will create a new numbered ADR in your ADR folder :
`xxx-my-new-awesome-proposition.md`.
Next, just open the file in your preferred markdown editor and starting writing your ADR.
