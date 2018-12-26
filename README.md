# ccan-archiver

`ccan-archiver` is a command-line program to archive the content of [ccan.de](https://ccan.de) and the [Clonk-Center Archive](https://cc-archive.lwrl.de).
The focus of this program is not archiving entire sites, but only downloadable items (excluding comments, the forum etc.) and some metadata.

The program directly downloads to a zip file called `result.zip`. In this file you can find a `README.md` file that documents the structure.

_Note: The file will be about 6.5GB in size (last checked 25.12.2018)._

### Dependencies

This program mainly depends on the go standard library. The only other package needed is `golang.org/x/net/html`.
You can install it with `go get golang.org/x/net/html` or `go get ./...`.

### Downloading

You can download binaries from the releases section of this repository.

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

The files are in the following schema:

 > `site/username/name.ext`

 - `site`: The site name where the item was hosted. This is either "CCAN" or "Clonk-Center"
 - `username`: username of the user that uploaded the file to CCAN.de
 - `name.ext`: File to download

There are also metadata files according to this schema:

 > `site/username/name.ext.json`

Exception: `README.md` and `failed.json` (this file only exists if a download failed, you can find all metadata there)


#### Metadata

The json file for all "CCAN" items contains the following entry:

```json
{
  "name": Display name of the object,
  "date" Upload date,
  "download_count": Number of downloads,
  "author": Name of the uploader,
  "votes": Number of votes,
  "category": Category,
  "engine": Engine for which this file was created,
  "download_link": download link to ccan.de (Usually a redirect)
}
```

The json file for all "Clonk-Center" items contains the following entry:

```json
{
  "name": Display name of the object,
  "date" Upload date,
  "download_count": Number of downloads,
  "author": Name of the uploader,
  "engine": Engine for which this file was created,
  "download_link": download link to cc-archive.lwrl.de,
  "description": The description that can be found at the page for this item (as Markdown text)
}
```
Read more in the `README.md` at the root of your archive after it has been downloaded.

### License

MIT, see [LICENSE](LICENSE)