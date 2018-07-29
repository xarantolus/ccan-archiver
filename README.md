# ccan-archiver

`ccan-archiver` is a command-line program to archive the content of [ccan.de](https://ccan.de).
The focus of this program is not archiving the entire site, but only downloadable items (excluding comments, the forum etc.) and some metadata.

The program directly downloads to a zip file called `result.zip`. In this file you can find a `README.md` file that documents the structure.


### Dependencies

This program mainly depends on the go standard library. The only other package needed is `golang.org/x/net/html`.
You can install it with `go get golang.org/x/net/html` or `go get ./...`.


### Copiling

To compile this program, you need to have [`go`](http://golang.org/) and [`go-bindata`](https://github.com/jteeuwen/go-bindata) installed.

_Note: This was tested with `go version go1.10.3 windows/amd64` and `go-bindata 3.1.0`, but it should work on any platform supported by go._

First, you have to generate some assets using go generate:
```
go generate zipfactory/*.go
```

After this, you can build the executable:
```
go build
```


### Structure

This is a brief documentation of the content of the resulting file called `result.zip`, you can find more in the file `README.md` in the archive.

The files in it are in the following schema:
```
username/name.ext
```

 - `username`: username of the user that uploaded the file to CCAN.de
 - `name.ext`: File to download

There are also metadata files according to this schema:

```
username/name.ext.json
```

Exception: `README.md`


The metadata file contains the following entry:

```
{
  "name": Display name of the scenario,
  "date." Upload date,
  "download_count": Number of downloads,
  "author": Name of the uploader,
  "votes": Number of votes,
  "category": Category_,
  "engine": Engine for which this file was created_,
  "download_link": _Downloadlink to ccan.de (Forwarding-Link)_,
  "direct_link": Direct download link, usually on another server
}
```

### License

MIT