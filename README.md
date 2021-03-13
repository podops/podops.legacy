# PodOps - Programmable Podcasts

[![Join the chat at https://gitter.im/podops/help](https://badges.gitter.im/podops/help.svg)](https://gitter.im/podops/help?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Podops is podcast infrastructure platform that provides functionallity to automate your podcast creation workflow. It allows you to create the podcast feed and delivers media assets like mp3s or images to podcast clients. 

The platform follows an API-first approach and is very light on user-facing frontends. Almost all interactions with Podops happens either using `po` the command line interface or by directly integrating the Rest API.

If you need support or have ideas for improving Podops, please join the [Podops Gitter community](https://gitter.im/podops/) or visit the [GitHub Issues section](https://github.com/podops/podops/issues) of this repo. Please visit the Podops document repository for installation instructions and documentation.

If you find this project interesting, please consider starring the project on GitHub.

## Open-source but not open-contribution?

I am grateful for community involvement, bug reports, & feature requests but I do not expect code contributions. I will gladly review [pull requests](https://github.com/podops/podops/pulls) however.

There is another way to help the project: many of the most valuable contributions are in the forms of [testing, feedback, and documentation](https://github.com/podops/podops/issues). This helps to harden the software and streamlines its usage for other users.

Should you wish to contribute, please review the [contribution guidelines](/docs/contributing.md) first.

## Getting started

To use Pooops, you need an API access key from `https://api.podops.dev`. For this, you have to install `po`, the command line tool first. After installing it, you can register your account and get your API key.

**Note:** While all the Podops code is [here on GitHub](https://github.com/podops), there is no step-by-step guide how install it on your own infrastructure at the moment. This will come later. If you want to deploy Podops on your own infrastructure *TODAY*, join the [community on Gitter](https://gitter.im/podops/) and send me a DM at `@mickuehl` and we will sort it out.

## Documentation
The documentation repository is [here](/docs)

## Examples
Examples on how to use the Command Line Interface or Go Client SDK to create and publish your podcast are [here](/docs/tutorial).

## Development
A description how to build the codebase and how to test locally is [here](/docs/development.md).
