---
name: golang-refactor-tools
description: >
  Go code refactoring with automated tools: gofmt -r, gopatch, and rsc.io/rf.
  Use this skill when performing mechanical rewrites across a codebase — renaming
  identifiers, migrating API calls, moving types between packages, rewriting
  expressions, or replacing deprecated patterns. Triggers when the task involves
  bulk code transformation, API migration, structural rewrite, or the user mentions
  gofmt -r, rf, gopatch, coccinelle-style patching, or large-scale rename/move in Go.
---

# Go Refactoring Tools

## Decision table

| Situation | Tool |
|---|---|
| Swap one expression for another | `gofmt -r` |
| Structural multi-line or multi-statement rewrite | `gopatch` |
| Type-sensitive expression rewrite | `rf ex` |
| Rename / move identifier, type, function, package | `rf mv` |
| Delete declarations | `rf rm` |

> **Workflow**: dry-run first (`gopatch -d` / `rf -diff`), review the diff, then apply.
> After any tool run: `goimports -w . && go build ./... && go test ./...`. Commit before bulk rewrites.

---

## `gofmt -r` — quickest, ships with Go

Syntax-aware but **not type-aware**. Pattern variables must be **single lowercase letters**; everything else is matched literally. `-w .` recurses naturally — `./...` is a Go package pattern and will error.

```bash
gofmt -r 'a.OldMethod(b) -> a.NewMethod(b)' -w .
gofmt -r 'ioutil.ReadAll(a) -> io.ReadAll(a)' -w .   # ioutil/io are literal; a is wildcard
```

After running: fix imports with `goimports -w .` (gofmt -r never touches imports).

---

## `gopatch` — structural patch files

Operates on syntax, understands statements and imports. Can work on partially-invalid code.

```bash
go install github.com/uber-go/gopatch@latest
gopatch -p my.patch ./...      # rewrite in place
gopatch -d -p my.patch ./...   # dry-run: print diff
gopatch ./... <<'EOF'          # read patch from stdin (omit -p)
...
EOF
```

### Patch anatomy

```patch
# Description (shown in -d output)
@@
var x expression    # matches any Go expression
var n identifier    # matches any single identifier
@@
-old code
+new code
```

**Metavariable types**: `expression` (calls, field accesses, literals…), `identifier` (single name).
**Undeclared identifiers** are matched literally — `a`, `b`, `c` in the body match only args named exactly `a`, `b`, `c`. Use `...` (elision) or declare `var a expression` explicitly.
**`...`**: matches zero or more statements/arguments. Multiple `@@…@@` blocks = multiple patches run in order.

### Examples

**Replace call + fix import in one patch:**
```patch
@@
@@
-import "io/ioutil"
+import "io"

-ioutil.ReadAll(...)
+io.ReadAll(...)
```

**Remove boilerplate with elision:**
```patch
# Delete redundant gomock.Controller.Finish()
@@
var ctrl, gomock identifier
var t expression
@@
 import gomock "github.com/golang/mock/gomock"

 ctrl := gomock.NewController(t)
 ...
-defer ctrl.Finish()
```

**Replace expression with sub-expression wildcard:**
```patch
@@
var x expression
@@
-time.Now().Sub(x)
+time.Since(x)
```

---

## `rsc.io/rf` — rename, move, expression rewrite

> ⚠ **Experimental**: rf's README says *"incredibly rough and likely to be buggy and change incompatibly."* Prefer `gopatch` for structural work; use `rf` when type-awareness is required (renames, moves, typed expression rewrites).

Takes a **single script argument** and operates on the module in the current directory — no `./...`.

```bash
go install rsc.io/rf@latest
rf -diff 'mv Foo Bar'    # show diff without writing
```

### `mv` — rename or move

```bash
rf 'mv OldType NewType'
rf 'mv pkg.OldFunc pkg.NewFunc'
rf 'mv mypkg.Foo otherpkg.Foo'
rf 'mv MyStruct.OldField MyStruct.NewField'
```

### `ex` — type-aware expression rewrite

Variables inside `{ }` are **typed** metavariables — only values of the declared type match, so there are no false positives on same-named identifiers from other packages. If a type doesn't precisely match the AST node, the rule is silently skipped. Multiple rules go in one block.

```bash
rf 'ex {
  var s string
  var r io.Reader
  fmt.Sprintf("%s", s) -> s
  ioutil.ReadAll(r)    -> io.ReadAll(r)
}'
```

`avoid` prevents rewriting inside named method bodies:
```bash
rf 'ex {
  var m myMap
  var k, v int
  avoid myMap.Get
  m[k] -> m.Get(k)
}'
```

### `rm` / `add`

```bash
rf 'rm OldFunc'
rf 'rm MyType.DeprecatedField'
rf 'add MyStruct.NewField int `json:"new_field"`'
```

---

## Common recipes

### Rename a public type

```bash
rf 'mv pkg.OldName pkg.NewName'
```
