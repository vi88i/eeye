# Code Review Guidelines for eeye

This document outlines the code review standards and best practices for the eeye stock screener project.

## General Principles

1. **Readability First**: Code should be self-documenting and easy to understand
2. **Consistency**: Follow existing patterns and conventions in the codebase
3. **Simplicity**: Prefer simple, straightforward solutions over clever ones
4. **Performance**: Consider performance implications, especially for data-intensive operations
5. **Error Handling**: Always handle errors explicitly and provide meaningful error messages

- Don't complain about closure issue because in Go 1.25 and above loop variables have per-iteration scope.

- Don't complain about test cases.

- Don't complain about SQL injection because no inputs are received from user, unless you feel something is coming from external source let me know.

- Consistent logging messages

**Improvements:**
- ✅ Proper error handling and propagation
- ✅ Parameterized query (prevents SQL injection)
- ✅ Resource cleanup with defer
- ✅ Error checks for scan operations
- ✅ Returns error type
- ✅ Clear documentation
- ✅ Checks rows.Err() after iteration

## References

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
