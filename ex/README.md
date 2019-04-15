exception
=========

This is a simple library for wrapping an `error` with a stack trace.

## Key Concepts

An exception is an error with additional context; message and most importantly, the stack trace at creation.

Concepts:
- `Class`: A grouping error; you should be able to test if exceptions are similar by comparing the exception classes. When an exception is created from an error, the class is set to the original error.
- `Message`: Additional context that is variable, that would otherwise break equatibility with exceptions. You should put extra descriptive information in the message.
- `Inner`: A causing exception or error; if you have to chan multiple errors together as a larger grouped exception chain, use `WithInner(...)`.
- `StackTrace`: A stack of function pointer / frames giving important context to where an exception was created.

## Usage

If we want to create a new exception we can use `New`

```go
return ex.New("this is a test exception")
```

`New` will create a stack trace at the given line. It ignores stack frames within the `exception` package itself. If you'd like to add variable context to an exception, you can use `WithMessagef(...)`:

If we want to wrap an existing golang `error` all we have to do is call `New` on that error.

```go
file, err := os.ReadFile("my_file.txt")
if err != nil {
    return ex.New(err)
}
```

If we want to add an inner exception, i.e. a causing exception, we can just add it with `.WithInner(...)`

```go
file, err := os.ReadFile("my_file.txt")
if err != nil {
    return ex.New("problems reading the config").WithInner(err)
}
```

A couple properties of `New`:
* It will return nil if the input `class` is nil.
* It will not modify an error that is actually an exception, it will simply return it untouched.
* It will create a stack trace for the class if it is not nil, and assign the class from the existing error.

## Formatted Output

If we run `fmt.Printf("%+v", ex.New("this is a sample error"))` we will get the following output (assuming we're running the statement in an http server somewhere):

```text
Exception: this is a sample error
       At: foo_controller.go:20 testExceptions()
           http.go:198 func1()
           http.go:213 func1()
           http.go:117 func1()
           router.go:299 ServeHTTP()
           server.go:1862 ServeHTTP()
           server.go:1361 serve()
           asm_amd64.s:1696 goexit()
```
