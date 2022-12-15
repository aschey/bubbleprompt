(cd "$(dirname $0)/../.." &&
    gomarkdoc --template-file package=./docs/scripts/package.gotxt --template-file text=./docs/scripts/text.gotxt --output ./docs/docs/03-Input/07-API.mdx ./input  &&
    gomarkdoc --template-file package=./docs/scripts/package.gotxt --template-file text=./docs/scripts/text.gotxt --output ./docs/docs/03-Input/02-Simple/02-API.mdx ./input/simpleinput  &&
    gomarkdoc --template-file package=./docs/scripts/package.gotxt --template-file text=./docs/scripts/text.gotxt --output ./docs/docs/03-Input/03-Command/02-API.mdx ./input/commandinput  &&
    gomarkdoc --template-file package=./docs/scripts/package.gotxt --template-file text=./docs/scripts/text.gotxt --output ./docs/docs/03-Input/04-Lexer/02-API.mdx ./input/lexerinput &&
    gomarkdoc --template-file package=./docs/scripts/package.gotxt --template-file text=./docs/scripts/text.gotxt --output ./docs/docs/03-Input/05-Parser/02-API.mdx ./input/parserinput &&
    gomarkdoc --template-file package=./docs/scripts/package.gotxt --template-file text=./docs/scripts/text.gotxt --output ./docs/docs/05-Suggestion/02-API.mdx ./suggestion
)