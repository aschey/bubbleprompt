(cd "$(dirname $0)/../.." &&
    godocdown -template=./docs/scripts/md.tmpl .\\editor\\simpleinput > ./docs/docs/02-Editor/01-Simple/API.md &&
godocdown -template=./docs/scripts/md.tmpl .\\editor\\commandinput > ./docs/docs/02-Editor/02-Command/API.md)