Multidoc yamldiff
=================

A CLI tool to diff two YAML/JSON files.

Inspired by https://github.com/sahilm/yamldiff but tries to work with multi-document YAML files as well.

It’s meant to be useful while diffing kubernetes yaml configs.

Installation
------------

```
go get github.com/mmikulicic/multidoc_yamldiff
```

Usage
-----

```
multidoc_yamldiff old.yaml new.yaml
```

Caveats
-------

This tool tries to deal with unordered documents but it’s not very smart: documents are simply sorted
by their (canonicalized) content. This means that it might work when there are small diffs at the tail of the yaml objects.
This seems to hold true for many changes to kubernetes object for which this tool has been designed.

In particular it can tell you for sure when two document sets are not different.

Also, the output of the diff is not yaml.

NOTE: This software has been developed in a rush while working on something that required a diffing tool.
