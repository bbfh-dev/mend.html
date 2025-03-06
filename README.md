# 🔩 mend.html

Mend is a simple HTML template processor designed to, but not limited to be used to generate static websites.

> [!CAUTION]
> This project is currently in **Beta**, meaning that it is <u>NOT</u> production ready.

<!-- vim-markdown-toc GFM -->

* [Installation](#installation)
* [Usage](#usage)
    * [Example Usage](#example-usage)
* [Language Specification](#language-specification)
    * [1. Comment Statements](#1-comment-statements)
        * [Extend Statement](#extend-statement)
        * [If Statement](#if-statement)
        * [Range Statement](#range-statement)
        * [Include Statement](#include-statement)
        * [Slot Statement](#slot-statement)
    * [2. Expression Statements](#2-expression-statements)

<!-- vim-markdown-toc -->

# Installation

Download the [latest release](https://github.com/bbfh-dev/mend.html/releases/latest) or install via the command line with:

```bash
go install github.com/bbfh-dev/mend.html
```

# Usage

Run `mend --help` to display usage information.

## Example Usage

```bash
mend build -s '{"title":"Hello World!","filename":"index.html","items":[]}' example/index.html
```

This command builds the `example/index.html` file along with all its dependencies using the provided JSON parameters.

# Language Specification

> [!TIP]
> In the documentation, `[argument]` denotes required arguments and `(argument)` denotes optional ones.

Mend processes a file's content unchanged until it encounters one of two types of mend statements:

## 1. Comment Statements

> [!IMPORTANT]
> Mend comments must be on separate, single lines to be recognized and processed.

There are two types of comment statements:

- `<!-- @... -->` — A simple inline mend statement.
- `<!-- #... -->` paired with `<!-- /... -->` — A mend block used to wrap a section of content.

### Extend Statement

**Syntax:**\
`<!-- #extend [filename] (parameters) -->` ... `<!-- /extend -->`

This statement extends a referenced file. Use the [Slot statement](#slot-statement) within the parent file to define where the child content should be inserted.

### If Statement

**Syntax:**\
`<!-- #if [name] [operator] [value] -->` ... `<!-- /if -->`

This block conditionally removes its enclosed content if the specified condition evaluates to false.

**Supported Operators:**

- `==` (equals)
- `!=` (does not equal)
- `has` (checks if an array contains a specified element)
- `lacks` (checks if an array does not contain a specified element)

### Range Statement

**Syntax:**\
`<!-- #range [parameter] -->` ... `<!-- /range -->`

This block iterates over an array. To access properties of the current item, prefix the expression with `#` (for example, `{{ #.child_property }}`).

### Include Statement

**Syntax:**\
`<!-- @include [filename] (parameters) -->`

This inline statement inserts the contents of the referenced file directly into the current document.

### Slot Statement

**Syntax:**\
`<!-- @slot -->`

This statement marks the insertion point for content when a file is extended. **Note:** Each file can declare only one slot.

## 2. Expression Statements

Expression statements insert values directly into the output.

- An expression beginning with `.` accesses a property from the input JSON (e.g. `{{ .path.to.name }}` retrieves the value at `path.to.name`).
- Using just `.` outputs the entire JSON object (the root).

Expressions can also include modifiers, using the format: `{{ modifier_name .path.to.property }}`.

**Supported Modifiers:**

- `length`, `len`, or `size` — Returns the length of an array.
- `quote` — Wraps the output value in double quotes.
