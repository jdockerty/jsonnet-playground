# jsonnet-playground

Inspired by the [Go playground](https://go.dev/play/) and others.

This project is built with Go, [`templ`][templ], and [htmx][htmx].

It is available at [jsonnet.jdocklabs.co.uk][hosted-playground] for anyone[^1] to use.

## How it works

Pages are directly rendered and served by the Go HTTP server through the `templ`
library.

Jsonnet is taken as input and evaluated through the Jsonnet VM, available through
the `/api/run` path, and htmx's target functionality displays the evaluated input
or error message in the corresponding output box.

Sharing is achieved through the creation of a shareable hash, the Jsonnet must be
valid to be shared.
If a shareable snippet has recently been created, it can be shared with others
by simply sending them the link. This will then be loaded by the server and displayed.

### Features

- Evaluate Jsonnet in browser
- Share snippets with others
- Utilise [`kubecfg`](https://github.com/kubecfg/kubecfg) functionality

## Development

The `justfile` in this repository provides various commands for local development.

`just run` will run the server and make the relevant static assets available[^2].

[^1]: Shared snippets are currently stored in memory, this is to keep the project near-free to run.
[^2]: The assets are named `kodata` alongside the server as the application is built and deployed with `ko`.

[htmx]: https://htmx.org/docs/#introduction
[templ]: https://github.com/a-h/templ
[hosted-playground]: https://jsonnet.jdocklabs.co.uk
