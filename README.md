# Billie Muk Personal Site

## Commands

- Build site: `go run . build`
- Dev server (live reload): `go run . serve`
- New post scaffold: `go run . new "Post Title"`

## Write a new post

1) Run the scaffold command:

```
go run . new "My Great Post"
```

2) Edit the generated file in `content/posts/`.

The file includes frontmatter:

```
---
title: "My Great Post"
date: YYYY-MM-DD
summary: "..."
draft: true
---
```

3) Set `draft: false` when ready to publish, then run `go run . build`.
