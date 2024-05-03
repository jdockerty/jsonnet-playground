# jsonnet-playground

Inspired by the [Go playground](https://go.dev/play/) and others.

This project is built with Go, [`templ`][templ], and [htmx][htmx].

It is available at [jsonnet.jdocklabs.co.uk][hosted-playground] for public use[^1].

# Development

The `justfile` in this repository provides various commands for local development.

`just run` will run the server and make the relevant static assets available[^2].

[^1]: Shared snippets are currently stored in memory, this is to keep the project near-free to run.
[^2]: The assets are named `kodata` alongside the server as the application is built and deployed with `ko`.

[htmx]: https://htmx.org/docs/#introduction
[templ]: https://github.com/a-h/templ
[hosted-playground]: https://jsonnet.jdocklabs.co.uk
