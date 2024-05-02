# jsonnet-playground

Inspired by the [Go playground](https://go.dev/play/) and others.

This project is built with Go and [htmx][htmx].

It is available at [jsonnet.jdocklabs.co.uk][hosted-playground] for broader use.
Keep in mind that snippets are ephemeral and may expire at anytime, this is because
storage is entirely in-memory for cloud billing purposes.

# Development

The `justfile` in this repository provides various commands for local development.

`just run` will run the server and make the relevant static assets available[^1].

[^1]: The assets are named `kodata` alongside the server as the application is built and deployed with `ko`.

[htmx]: https://htmx.org/docs/#introduction
[hosted-playground]: https://jsonnet.jdocklabs.co.uk
