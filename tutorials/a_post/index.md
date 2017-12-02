
# A post

This is an example of using _shorthand_ to build a blog post.

## Parts 

1. [build.shorthand](build.shorthand) - this is the _shorthand_ file orchestrating everything
2. [post.md](post.md) - the Post we want to make in Markdown
3. [template.md](template.md) - is a markdown file with _shorthand_ embedded, it is our template because we embedded _shorthand_
    + [template.html](template.html) - an HTML rendering of the template to show HTML but un-resolved _shorthand_ still in it.
4. [post.html](post.html) - the final post

## Process

Run the build process

```shell
    shorthand build.shorthand
```

This reads in build.shorthand which intern pulls in template.md to use
as a template with content taken from post.md rendering the final post.html.




