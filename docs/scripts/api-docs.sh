(cd "$(dirname $0)/../.." &&
    gomarkdoc --template-file package=./docs/scripts/package.gotxt --template-file index=./docs/scripts/index.gotxt --output ./docs/docs/03-Input/02-Simple/API.mdx ./input/simpleinput  &&
    gomarkdoc --template-file package=./docs/scripts/package.gotxt --template-file index=./docs/scripts/index.gotxt --output ./docs/docs/03-Input/03-Command/API.mdx ./input/commandinput  &&
    gomarkdoc --template-file package=./docs/scripts/package.gotxt --template-file index=./docs/scripts/index.gotxt --output ./docs/docs/03-Input/04-Lexer/API.mdx ./input/lexerinput &&
gomarkdoc --template-file package=./docs/scripts/package.gotxt --template-file index=./docs/scripts/index.gotxt --output ./docs/docs/03-Input/05-Parser/API.mdx ./input/parserinput)