# PodOps - Programmable Podcasts

[![Join the chat at https://gitter.im/podops/help](https://badges.gitter.im/podops/help.svg)](https://gitter.im/podops/help?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Podops is a podcast infrastructure platform that provides functionallity to automate your podcast creation workflow. It allows you to create the podcast feed and delivers media assets like mp3s or images to podcast clients. 

The platform follows an API-first approach and is very light on user-facing frontends. Almost all interactions with Podops happens either using `po` the command line interface or by directly calling the Rest API.

If you need support or have ideas for improving Podops, please join the [Podops Gitter community](https://gitter.im/podops/) or visit the [GitHub Issues section](https://github.com/podops/podops/issues) of this repo. Please visit the Podops [document repository](/docs) for installation instructions and documentation.

If you find this project interesting, please consider starring it here on GitHub!

## Getting started

To use Podops, you need an API access key for `https://api.podops.dev`. To get your key, you have to install `po`, the command line tool first. After installing it, you can register your account and request a API key.

**Note:** While all the Podops code is [here on GitHub](https://github.com/podops), there is no step-by-step guide how install it on your own infrastructure at the moment. This will come later. If you want to deploy Podops on your own infrastructure *TODAY*, join the [community on Gitter](https://gitter.im/podops/) and send me a DM at `@mickuehl` and we will sort it out.

### Download the binary

Download the latest binary for your system:

* [Linux/amd64](https://storage.googleapis.com/cdn.podops.dev/downloads/cli-0.9.7/po-linux-0.9.7.gz)
* [MacOS/amd64](https://storage.googleapis.com/cdn.podops.dev/downloads/cli-0.9.7/po-mac-0.9.7.gz)
* [Windows](https://storage.googleapis.com/cdn.podops.dev/downloads/cli-0.9.7/po-win-0.9.7.zip)

Unpack the archive and place the `po` binary in a directory that is on your `$PATH`. Verify that the command line tool is accessible:

```shell
$ po help
```

### Register and get the API key

```shell
$ po login <your_email_address>
```

Podops will send you an email with a confirmation code. Use the confirmation code to exchange it for your API access key. The code is `valid for 15min` and can only be used once. In case you missed this time-window, start over with the login command.

```shell
$ po auth <access_code>
```

This will retrieve your current API access code and place it in its default location `$HOME/.po/config`. You can verify that everything is setup correctly by issuing a command that requires authentication e.g.

```shell
$ po list
```

**Note:** If you login for the first time, Podops will send you an email to verify your Email-Address first. Confirm by following the link in the Email. The link is `valid for 15min`.

## Examples
Examples on how to use the Command Line Interface or Go Client SDK to create and publish your podcast are [here](/docs/tutorial).

## Documentation
The documentation repository is [here](/docs)

## Development
A description how to build the codebase and how to test locally is [here](/docs/development.md).

## Open-source but not open-contribution?

I am grateful for community involvement, bug reports, & feature requests but I do not expect code contributions at this point in time as there is not really a substantial user base or community in general. In case someone wants to contribute, I will gladly review [pull requests](https://github.com/podops/podops/pulls).

There is another way to help the project: many of the most valuable contributions are in the forms of [testing, feedback, and documentation](https://github.com/podops/podops/issues). This helps to harden the software and streamlines its usage for other users.

Should you wish to contribute, please review the [contribution guidelines](/docs/contributing.md) first.