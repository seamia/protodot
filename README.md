# protodot
transforming your `.proto` files into `.dot` files (and `.svg`, `.png` if you happen to have `graphviz` installed)

## online demo
you can find an online demo of this tool [here](https://protodot.seamia.net)

## data pipeline
<p align="center">
  <img src="https://protodot.seamia.net/pipeline.svg">
</p>


## installation
you can download the sources (from this page) and build `protodot` yourself, or
you can download pre-built binaries for a few selected environments (or from [here](https://github.com/seamia/protodot/tree/master/binaries)):

   * [Darwin](https://protodot.seamia.net/binaries/darwin)
   * [Linux](https://protodot.seamia.net/binaries/linux)
   * [Windows](https://protodot.seamia.net/binaries/windows)

`protodot` output is highly customizable through its configuration file and a set of templates files.
if you installed `protodot` from the binary distribution - you may want to extract aforementioned files by running:

```
   ./protodot -install
```

## command line arguments

   * `-src what.proto` - location and name of the source file, required
   * `-config config.json` - location and name of the configuration file, optional
   * `-select .one.two;three.four` - name(s) of the selected elements to show, optional, explained later in this document
   * `-output save-it-here` - name of the output file, optional
   * `-inc /abc/def;/xyz` - (semicolon separated) list of the include directories, optional


## configuration file
tbd

## selected output
sometimes the resulting diagram can be overwhelming.
you have an option to limit the output to the elements that interest you the most, hence `-select args` command line option.
so far, `args` in `-select args` can take one of the three available forms:
   * list of the elements (and their dependencies) that you want to see included (separated by `;`). the elements can be `enums`, `messages`, `rpc` methods and `services`.
   * if you specify `*` as an argument - this will result in the inclusion of the elements declared in the **main** `.proto` file (specified in `-src` argument) and their dependencies. in other words: all the **unused** elements declared in all the **included** `.proto` files will not be shown.
   * if you specify `imports` as an argument - `protodot` will generate import dependency graph (see an example below)


## an example of output
<p align="center">
  <img src="https://protodot.seamia.net/pipeline/svg">
</p>




## an illustration of effects of different `-select` options
using `https://github.com/googleapis/googleapis/blob/master/google/privacy/dlp/v2/dlp.proto` as the source

### not using `-select`
everything declared in the root `.proto` file and all the `imports` will be produced:
<p align="center">
  <img src="https://protodot.seamia.net/demo/dlp_full.svg">
</p>

### using `-select *`
only elements declared in the root `.proto` file and their dependencies will be produced:
<p align="center">
  <img src="https://protodot.seamia.net/demo/dlp_star.svg">
</p>

### using `-select imports`
only `imports` dependency graph will be produced:
<p align="center">
  <img src="https://protodot.seamia.net/demo/dlp_imports.svg">
</p>

### using `-select .ListDlpJobs`
in this particular case, the name of the `rpc` method was specified: this will result in production of the requested `rpc` method, it's encompasing `service` element and all the dependencied of the method:
<p align="center">
  <img src="https://protodot.seamia.net/demo/dlp_rpc.svg">
</p>

### using `-select .AnalyzeDataSourceRiskDetails`
same as above, but instead of `rpc` method, name of the `message` was specified:
<p align="center">
  <img src="https://protodot.seamia.net/demo/dlp_message.svg">
</p>


## how to (automatically) generate `.svg` and/or `.png` images from produced `.dot` file
1. install `graphviz` (see https://graphviz.gitlab.io/download/ for the instructions)
2. specify the location of the `dot` utility (which is a part of `graphviz)` in your version of configuration file, e.g.
```
{
	"locations": {
		"graphviz":	"/path/to/dot",
```
3. set approptiate `options` in your version of configuration file, e.g.
```
{
	"options" : {
		"generate .png file":		false,
		"generate .svg file":		true,

```

