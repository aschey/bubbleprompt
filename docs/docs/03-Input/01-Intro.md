## Tokenization

Bubbleprompt's input handling uses the concept of **tokenization** to parse
user input. This makes it easier to perform operations on specific parts of
the user input. A common tokenization method is to split input on
whitespace. For example, if the input is `this is some text`, the
tokenized input would be `["this", " ", "is", " ", "some", " ", "text"]`.
Each word would be considered a word token and each space would be a delimiter token.
Manipulating inputs in terms of tokens makes it easier to do common operations such as
getting the current word under the cursor.

The above example could be handled by a simple `strings.Split(str, " ")`, but what if
we wanted to parse `this is "some text"` as `["this", " ", "is", " ", "some text"]`?
Something like this would be complicated to do using simple string manipulation.
Instead we can use regex to define our tokens. A regex to define a token which either
contains no whitespace or does contain whitespace only if it's enclosed by quotes could look
like this: `("[^"]*"?)|[^\s]+`. And a regex for whitespace delimiters would be `\s+`.

Under the hood, Bubbleprompt uses [Participle](https://github.com/alecthomas/participle)
to handle input tokenization. Defining these rules yourself can be complicated, which is why
Bubbleprompt offers the [Simple](./Simple/Usage) input for common uses like handling whitespace delimited
tokens as described above.

If you need more control over your tokenization, you can use the [Lexer](./Lexer/Usage) input to
define custom tokenization strategies.

## Grammars and Parsing

Sometimes, operating on a simple structure like a list of tokens isn't enough if your app requires
a more structured input. Take a simple CLI command for example. It contains a command, zero or more subcommands,
zero or more arguments, and any number of short flags and long flags. Each of these components has unique meaning
which would be difficult to parse with a list of tokens alone. For these use cases, you can define a custom grammar
that represents your input structure.

For handling POSIX-style CLI input as described above, we provide the [Command](./Command/Usage) input.
To implement custom grammars, you can use the [Parser](./Parser/Usage) input.
