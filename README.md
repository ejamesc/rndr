# rndr

rndr serves templates that are rendered using
https://github.com/unrolled/render. I mostly use it when I need to quickly
prototype something, to be used in a Go webapp.

Installation: 

```
go install github.com/ejamesc/rndr
```

Usage: 

```
rndr -t <templates directory> -s <static directory>
```

For more instructions just run `rndr`.

Typical workflow: point this to a repo of templates, and point the static
folder to the static folder of my actual webapp. rndr is set to reload
templates for each request.
