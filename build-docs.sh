#!/bin/bash

set -e

if [ ! -d ./docs-bosh ]; then
  echo 'Missing docs-bosh repository'
  exit 1
fi

echo "Regenerate doc assets"
( cd ./docs-bosh-io && bundle exec bookbinder publish local )

echo "Remove old copy of docs"
rm -rf ./templates/docs

echo "Copy out assets generated by the bookbinder"
cp -R ./docs-bosh-io/final_app/public/docs ./templates/

echo "Done building docs"
