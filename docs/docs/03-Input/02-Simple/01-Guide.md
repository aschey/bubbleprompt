# Guide

## When to Use

This package can be used to tokenize simple delimiter-based inputs escaped by quotes.
For example, it can parse `hello "this is" 'some text'` as `["hello", "this is", "some text"]`.
It can also handle parsing with custom delimiters, such as `hello."this is".'some text'`
which can be parsed the same as the above using the [WithTokenRegex](API#func-withtokenregex) and [WithDelimiterRegex](API#func-withdelimiterregex) options.

## Example

See the package-level example in the [API Docs](./API)

## Options

See available options with examples in the [API Docs](./API#type-option)
