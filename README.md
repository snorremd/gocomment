# gocomment

A simple comment application for the modern age. Runs a RESTish comment server
built on [go](https://golang.org/) and an embeddable client side application
built on the [re-frame](https://github.com/Day8/re-frame) framework on top of
[ClojureScript](https://clojurescript.org/).

Still very much work in progress!

## Development Mode

Make sure you install both [go](https://golang.org/doc/install) and
[leiningen](https://leiningen.org/). Go >= 1.9 and Leiningen >= 2.8.1 on Java 9
is the current development environment used by myself. Older versions might not
work as expected.


### Run application:

Start the gocomment go http server. It will create the `comments.db` sqlite
database if it does not already exist and start serving the API at
[localhost:8080](http://localhost:8080).

```bash
DB=comments.db  HOST=localhost:8080 go run api/main.go
```

Start a nrepl and figwheel repl to compile the Clojurescript code and serve a
client side code with hot code reload at [localhost:3449](http://localhost:3449).
You can run `(js/alert "Hello browser!")` to check if the repl works as you
would expect.

```bash
lein dev
```

## Production Build

TODO
