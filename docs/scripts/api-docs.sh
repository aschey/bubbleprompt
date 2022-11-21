(cd "$(dirname $0)/../.." &&
    godocdown -template=./docs/scripts/md.tmpl .\\editor\\simpleinput > ./docs/docs/03-Editor/01-Simple/API.mdx &&
    godocdown -template=./docs/scripts/md.tmpl .\\editor\\commandinput > ./docs/docs/03-Editor/02-Command/API.mdx &&
godocdown -template=./docs/scripts/md.tmpl .\\editor\\parserinput > ./docs/docs/03-Editor/03-Parser/API.mdx)