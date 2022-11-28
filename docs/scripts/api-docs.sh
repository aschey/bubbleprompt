(cd "$(dirname $0)/../.." &&
    godocdown -template=./docs/scripts/md.tmpl ./input/simpleinput > ./docs/docs/03-Input/02-Simple/API.mdx &&
    godocdown -template=./docs/scripts/md.tmpl ./input/commandinput > ./docs/docs/03-Input/03-Command/API.mdx &&
    godocdown -template=./docs/scripts/md.tmpl ./input/lexerinput > ./docs/docs/03-Input/04-Lexer/API.mdx &&
godocdown -template=./docs/scripts/md.tmpl ./input/parserinput > ./docs/docs/03-Input/05-Parser/API.mdx)