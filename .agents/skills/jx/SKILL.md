---
name: jx
description: >
  Correct, idiomatic, high-performance JSON encoding and decoding with
  github.com/go-faster/jx. Use this skill whenever you write or review
  jx-based Encode/Decode methods, implement JSON marshaling over jx,
  use jx.Decoder / jx.Encoder / jx.Writer, or encounter questions about
  jx buffer safety, Capture, pooling, or Writer vs Encoder trade-offs.
  Also trigger when the file imports "github.com/go-faster/jx" and you
  are adding or changing any encoding / decoding logic.
---

# jx — JSON encoding and decoding

`github.com/go-faster/jx` is a zero-allocation, RFC 7159 JSON library
used as the foundation of [ogen](https://github.com/ogen-go/ogen). It
trades the convenience of `encoding/json` for direct control over every
byte — which means you get to make choices that matter for correctness
and performance.

---

## Three rules that matter most

### 1. `StrBytes`, `Num`, and `Raw` reference the internal buffer

These methods return slices that alias the decoder's read buffer.
They are valid only until the **next decoder call**.

```go
// BUG: key is overwritten before use
keys = append(keys, key)

// OK: copy immediately
keys = append(keys, string(key))  // or append([]byte{}, key...)
```

The same applies to `jx.Num` returned by `d.Num()` and `jx.Raw`
returned by `d.Raw()`. Use the `*Append` variants to copy into your
own buffer:

```go
n, err = d.NumAppend(n[:0])   // safe: appends into your slice
raw, err = d.RawAppend(raw[:0])
```

### 2. `Capture` enables multi-pass decoding — byte buffers only

`d.Capture(f)` saves decoder state, runs `f`, then rolls back.
Use it to peek at a discriminator field before deciding how to decode:

```go
var kind string
if err := d.Capture(func(d *jx.Decoder) error {
    return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
        if string(key) == "type" {
            v, err := d.StrBytes()
            if err != nil {
                return err
            }
            kind = string(v) // copy — buffer is temporary
            return err
        }
        return d.Skip()
    })
}); err != nil {
    return err
}

// Decoder is reset to before the Capture call.
switch kind {
case "foo":
    return s.Foo.Decode(d)
case "bar":
    return s.Bar.Decode(d)
}
```

**`Capture` does not work with `io.Reader` decoders.** It only works
when the decoder was created with `DecodeBytes` / `DecodeStr` /
`ResetBytes` (i.e., the full input is already in memory).

### 3. `jx.Writer` is faster but requires manual commas

`jx.Encoder` tracks a comma-state stack and inserts commas
automatically. `jx.Writer` omits that bookkeeping — ~1.7× faster in
benchmarks, but every non-first element must be preceded by an explicit
`w.Comma()` call. Use `Writer` in generated or hot-path code where the
structure is statically known; use `Encoder` for hand-written code.

---

## Decoder

### Creating a decoder

```go
d := jx.DecodeBytes(data)        // byte slice
d := jx.DecodeStr(`{"k":"v"}`)   // string literal
d := jx.Decode(r, 512)           // io.Reader, buffer size 512

// Pool reuse — decoder is reset on PutDecoder
d := jx.GetDecoder()
defer jx.PutDecoder(d)
d.ResetBytes(data)
```

### Decoding objects

Prefer `ObjBytes` over `Obj` — it avoids allocating a `string` for
every key. The key slice is only valid inside the callback; using
`string(key)` in the switch statement is safe and cheap.

```go
func (s *MyStruct) Decode(d *jx.Decoder) error {
    return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
        switch string(key) {
        case "name":
            v, err := d.Str()
            if err != nil {
                return err
            }
            s.Name = v
        case "count":
            v, err := d.Int()
            if err != nil {
                return err
            }
            s.Count = v
        default:
            return d.Skip() // must skip unknown fields or the decoder stalls
        }
        return nil
    })
}
```

### Decoding arrays

```go
return d.Arr(func(d *jx.Decoder) error {
    v, err := d.Str()
    if err != nil {
        return err
    }
    s.Items = append(s.Items, v)
    return nil
})
```

### Peeking at the next type

`d.Next()` returns the type of the next value without consuming it.

```go
switch d.Next() {
case jx.String:
    v, err := d.Str()
    ...
case jx.Null:
    if err := d.Null(); err != nil { return err }
    // value is null
case jx.Number:
    v, err := d.Int()
    ...
}
```

### Zero-copy string

`StrBytes` avoids a heap allocation by returning a slice into the
decoder buffer. Use it immediately — it is only valid until the next
decoder call. Do not store the result.

```go
raw, err := d.StrBytes()
if err != nil {
    return err
}
// Use raw here — e.g. switch, hex.Decode, or hand to a library that
// accepts []byte (like uuid.ParseBytes). Do NOT store raw in a field.
```

When you need to retain the string but want allocation-free reuse
across calls, use `StrAppend` with your own buffer:

```go
var buf []byte // reuse across iterations

buf, err = d.StrAppend(buf[:0]) // appends decoded string into buf
name = string(buf)              // copy once, into the final destination
```

### Numbers

`jx.Num` is a `[]byte` type. Like `StrBytes`, the value returned by
`d.Num()` aliases the decoder buffer. Use `NumAppend` to own the bytes.

```go
n, err := d.NumAppend(n[:0]) // append into your own slice
i64, err := n.Int64()
f64, err := n.Float64()
```

### String-encoded numbers

Some APIs encode numbers inside JSON strings (`"count": "42"`).
Decode them by extracting the string bytes and feeding them to a
nested decoder — no allocation because `StrBytes` is zero-copy:

```go
raw, err := d.StrBytes() // e.g. []byte("42")
if err != nil {
    return err
}
n, err := jx.DecodeBytes(raw).Int()
```

---

## Encoder

Commas are inserted automatically. The zero value is valid.

### Encoding a struct

```go
func (s *MyStruct) Encode(e *jx.Encoder) {
    e.ObjStart()
    defer e.ObjEnd()

    e.FieldStart("name")
    e.Str(s.Name)

    e.FieldStart("count")
    e.Int(s.Count)

    if s.Optional != "" {   // omit zero-value optional fields
        e.FieldStart("optional")
        e.Str(s.Optional)
    }
}
```

### Encoding arrays

```go
e.ArrStart()
for _, v := range items {
    e.Str(v)   // comma inserted automatically before each non-first element
}
e.ArrEnd()
```

### Nullable values

```go
e.FieldStart("value")
if ptr == nil {
    e.Null()
} else {
    e.Str(*ptr)
}
```

### Pool reuse

```go
e := jx.GetEncoder()
defer jx.PutEncoder(e)   // resets on return

s.Encode(e)
result := append([]byte{}, e.Bytes()...)  // copy before defer runs
```

### Primitive methods

```go
e.Str(v string)
e.ByteStr(v []byte)   // encodes bytes as a JSON string, no allocation
e.Int(v int) / e.Int64 / e.UInt64 / e.Int32 ...
e.Float64(v float64)
e.Bool(v bool)
e.Null()
e.Raw(v []byte)       // embed pre-encoded JSON verbatim
e.RawStr(v string)    // same, from string
e.Base64(v []byte)
e.Num(v jx.Num)
```

### Encoding pre-formatted values efficiently

When a value has a fixed or bounded byte representation, you can encode
it with zero allocations by formatting into a stack-allocated array and
writing the result as raw JSON.

**Fixed-size value (e.g. UUID — 38 bytes including quotes):**
```go
// Pre-encode with quotes included; write as raw bytes.
const quoted = 38 // 36 chars + 2 quotes
var dst [quoted]byte
dst[0] = '"'
dst[quoted-1] = '"'
hexEncode((*[36]byte)(dst[1:37]), id) // your hex-encode function
e.Raw(dst[:])
```

**Variable but bounded value (e.g. timestamp, duration):**
```go
// AppendFormat into a stack buffer; ByteStr encodes it as a JSON string.
var buf [64]byte
b := v.AppendFormat(buf[:0], time.RFC3339)
e.ByteStr(b)   // writes "...", no allocation
```

The key rule: `e.Raw(b)` embeds `b` verbatim (caller provides quotes if
needed); `e.ByteStr(b)` wraps `b` in JSON string quotes and escaping.

---

## Writer (faster, manual commas)

`jx.Writer` exposes `Buf []byte` directly and skips the comma-state
stack, making it the fastest way to produce JSON. You are responsible
for calling `w.Comma()` between every pair of adjacent values.

### Encoding a struct with Writer

For statically known fields, embed the comma and colon directly into
`RawStr` literals — this is what code generators do:

```go
func (s *MyStruct) Write(w *jx.Writer) {
    w.ObjStart()
    w.RawStr(`"name":`)       // first field — no leading comma
    w.Str(s.Name)
    w.RawStr(`,"count":`)     // subsequent fields — leading comma in literal
    w.Int(s.Count)
    w.ObjEnd()
}
```

For dynamic fields (map keys, optional fields), track first-element
state explicitly:

```go
w.ObjStart()
first := true
for k, v := range m {
    if !first {
        w.Comma()
    }
    first = false
    w.FieldStart(k)
    w.Str(v)
}
w.ObjEnd()
```

### Encoding arrays with Writer

```go
w.ArrStart()
for i, v := range items {
    if i != 0 {
        w.Comma()
    }
    w.Str(v)
}
w.ArrEnd()
```

### Pool reuse

```go
w := jx.GetWriter()
defer jx.PutWriter(w)

s.Write(w)
result := append([]byte{}, w.Buf...)  // copy before defer runs
```

---

## Encoder vs Writer

| | `jx.Encoder` | `jx.Writer` |
|---|---|---|
| Commas | Automatic | Manual (`w.Comma()`) |
| Indentation | `e.SetIdent(n)` | Not supported |
| Benchmark speed | Fast | ~1.7× faster |
| When to use | Hand-written code | Generated / hot-path code |

---

## Conventional interfaces

ogen-generated code uses these signatures — implement them on your
types to stay interoperable:

```go
func (s *MyStruct) Encode(e *jx.Encoder) { ... }  // encoding
func (s *MyStruct) Decode(d *jx.Decoder) error { ... }  // decoding
```

---

## Common mistakes

| Mistake | Fix |
|---|---|
| Storing `StrBytes` / `Num` / `Raw` result beyond the callback | Copy: `string(b)` or `append([]byte{}, b...)` |
| No `d.Skip()` in `ObjBytes` default branch | Decoder stalls; always skip |
| `Capture` on an `io.Reader`-backed decoder | Only works with byte-backed decoders |
| Adding your own commas with `Encoder` | `Encoder` is automatic; extra commas corrupt output |
| Forgetting commas with `Writer` | `Writer` never adds commas; every non-first value needs `w.Comma()` |
| Forgetting `e.Bytes()` before `PutEncoder` | Buffer is reset on pool return; copy first |
