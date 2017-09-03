# rndr

rndr serves templates that are rendered using
https://github.com/unrolled/render. I mostly use it when I'm designing
templates for my Go webapps.

Installation: 

```
go install github.com/ejamesc/rndr
```

Usage: 

```
rndr -t <templates directory> -s <static directory>
```

For more instructions just run `rndr`.

Typical workflow: point this to a directory of templates, and point the static
flag to the static directory of my actual webapp. rndr is set to reload
templates for each request.

### misc
- The base template is base.html, and should be located in
  `<templatesdir>/base.html`
- All static assets are served under `/static`
