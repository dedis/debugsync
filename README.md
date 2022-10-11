[![Go lint](https://github.com/dedis/debugtools/workflows/Go%20lint/badge.svg)](https://github.com/dedis/debugsync/actions?query=workflow%3A%22Go+lint%22)
[![Go test](https://github.com/dedis/debugtools/workflows/Go%20test/badge.svg)](https://github.com/dedis/debugsync/actions?query=workflow%3A%22Go+test%22)
[![Coverage Status](https://coveralls.io/repos/github/dedis/debugtools/badge.svg?branch=main)](https://coveralls.io/github/dedis/debugsync?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/dedis/debugtools)](https://goreportcard.com/report/github.com/dedis/debugsync)

# debugtools
This repository aims to bring useful tools to help debugging applications 
written in GO.

## sync
Package that helps debugging mutexes and workgroups in a distributed system. Any
system deadlock will fire off after a defined timeout.
This is a drop-in replacement for the sync standard library.

## channel
Package that helps debugging locked channels. The created channel will generate
a log if we need to wait more than the timeout before writing or reading a value
to/from the channel.