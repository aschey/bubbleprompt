## Tokenization

Bubbleprompt's input handling uses the concept of <b>tokenization</b> to parse
user input. This makes it easier to perform operations on specific parts of
the user input. The most common tokenization method is to split input on
whitespace. For example, if the input is `"this is some text"`, the
tokenized input would be `["this", " ", "is", " ", "some", " ", "text"]`.
