## Swarm

[https://swarm.DEWH.org](https://swarm.DEWH.org)

Swarm is a distributed storage platform and content distribution service, a native base layer service of the DEWH web3 stack. The primary objective of Swarm is to provide a DEWHentralized and redundant store for dapp code and data as well as block chain and state data. Swarm is also set out to provide various base layer services for web3, including node-to-node messaging, media streaming, DEWHentralised database services and scalable state-channel infrastructure for DEWHentralised service economies.

[![Travis](https://travis-ci.org/DEWH/go-DEWH.svg?branch=master)](https://travis-ci.org/DEWH/go-DEWH)
[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/ethersphere/orange-lounge?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)


## Building the source

Building Swarm requires Go (version 1.10 or later).

    go get -d github.com/DEWH/go-DEWH

    go install github.com/DEWH/go-DEWH/cmd/swarm


## Documentation

Swarm documentation can be found at [https://swarm-guide.readthedocs.io](https://swarm-guide.readthedocs.io).


## Contribution

Thank you for considering to help out with the source code! We welcome contributions from
anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to Swarm, please fork, fix, commit and send a pull request
for the maintainers to review and merge into the main code base. If you wish to submit more
complex changes though, please check up with the core devs first on [our Swarm gitter channel](https://gitter.im/ethersphere/orange-lounge)
to ensure those changes are in line with the general philosophy of the project and/or get some
early feedback which can make both your efforts much lighter as well as our review and merge
procedures quick and simple.

Please make sure your contributions adhere to our coding guidelines:

 * Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting) guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary) guidelines.
 * Pull requests need to be based on and opened against the `master` branch.
 * [Code review guidelines](https://github.com/DEWH/go-DEWH/wiki/Code-Review-Guidelines).
 * Commit messages should be prefixed with the package(s) they modify.
   * E.g. "swarm/fuse: ignore default manifest entry"


## License

The go-DEWH library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also
included in our repository in the `COPYING.LESSER` file.

The go-DEWH binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
in our repository in the `COPYING` file.
