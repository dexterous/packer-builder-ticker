[![baby-gopher](https://raw.github.com/drnic/babygopher-site/gh-pages/images/babygopher-badge.png)](http://www.babygopher.org)
# Packer Ticker Builder
This is a [packer](packer.io) `builder` that I threw together to try out the plugin API.
It doesn't build any artifacts. It just spits out periodic, ficticious progress updates for a given duration.

## Installation
No binaries are currently published. Simply `go get` the package source and build/install it.

## Configuration
As described on the [packer plugin installation page](http://www.packer.io/docs/extend/plugins.html#toc_2), add the
plugin to your `~/.packerconfig`.

```
{
  "builders": {
    "ticker": "packer-builder-ticker"
  }
}
```

## Template definition
### Basic Example
In it's simplest form, all you need to do is declare a builder of type `ticker`.

```
{
  "builders": [
    { "type": "ticker" }
  ]
}
```

### Configuration Reference
The builder has two configuration options, both optional.

1. `period` (int): the period of one notification; i.e., the builder will publish one notification every `period`
    seconds [default: 1]
1. `duration` (int): the duration for which the builder will run in seconds [default: 5]

### Example with configuration options
A configuration as below...

```
{
  "builders": [
    { "type": "ticker", "period": 2 },
    { "type": "ticker", "name": "other", "duration": 10 }
  ]
}
```

... will produce output like this.

```
$ packer build template.json 
ticker output will be in this color.
other output will be in this color.

==> ticker: Running(  2,   5)...
==> other: Running(  1,  10)...
==> other: Building... 1.00007767s
==> ticker: Building... 2.000076307s
==> other: Building... 2.000087931s
==> other: Building... 3.00008243s
==> ticker: Building... 4.000094309s
==> other: Building... 4.00009628s
==> ticker: Done! Stopping...
Build 'ticker' finished.
==> other: Building... 5.000100202s
==> other: Building... 6.000102877s
==> other: Building... 7.000111495s
==> other: Building... 8.000111263s
==> other: Building... 9.00008563s
==> other: Building... 10.000107023s
==> other: Done! Stopping...
Build 'other' finished.

==> Builds finished. The artifacts of successful builds are:
```
