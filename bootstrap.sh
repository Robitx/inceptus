#!/bin/bash

orig_repository="github.com/robitx/inceptus/_server_template"

delimeter="========================================================================"

usage()
{
echo """Hello there,
this script allows you to quickly bootstrap new golang server based on inceptus/_server_template.

Example:
./bootstrap.sh -n shiny_new_project -r github.com/XXXX/shiny_new_project -d /tmp/shiny_new_project

cd /tmp/shiny_new_project; docker-compose -f docker-compose-dev.yaml up;

Usage:
    bootstrap.sh -h Display this help message.

    bootstrap.sh -n PROJECT_NAME -r REPOSITORY -d DIRECTORY
    where:
      - PROJECT_NAME will be used instead of \"server_template\"
      - REPOSITORY (like github.com/XXX/YYY) will be used instead of:
        \"$orig_repository\"
      - DIRECTORY is a folder where you want to boostrap the server
"""
}
 while getopts "n:r:d:h" opt; do
  case "${opt}" in
    h)
      usage;
      exit 0;
      ;;
    \? )
      echo "Invalid Option: -$OPTARG" 1>&2
      exit 1
      ;;
    n)
      name=$OPTARG
      ;;
    d)
      directory=$OPTARG
      ;;
    r)
      repository=$OPTARG
      ;;
  esac
done 

if [ -z "$name" ] || [ -z "$directory" ] || [ -z "$repository" ];
then
  usage;
  exit 1;
fi

echo -e "Trying to make $directory in case it doesn't exists:"
mkdir -p "$directory";

echo -e "Copying files from _server_template/* to $directory:"
cp -v -r ./_server_template/* "$directory";
echo -e "DONE\n";

echo "Renaming some files:"
while read file; do
  mv -v $file $(echo $file | sed -e "s/server_template/$name/");
done < <(find "$directory" | grep server_template);
echo -e "DONE\n";

echo "Editing files with imports:"
while read file; do
  echo $file;
  sed -i.REMOVE_BACKUP -e "s|$orig_repository|$repository|g" $file;
done < <(grep -irl $orig_repository $directory);
echo -e "DONE\n";


echo "Editing files containing \"server_template\":"
name_upper=$(echo $name | tr [:lower:] [:upper:]);
echo $name_upper;
while read file; do
  echo $file;
  sed -i.REMOVE_BACKUP -e "s|SERVER_TEMPLATE|$name_upper|g" $file;
  sed -i.REMOVE_BACKUP -e "s|server_template|$name|g" $file;
done < <(grep -irl "SERVER_TEMPLATE" $directory);
echo -e "DONE\n";


echo "Removing sed backup files:"
while read file; do
  rm -v $file;
done < <(find "$directory" | grep ".REMOVE_BACKUP");
echo -e "DONE\n";

echo "Remainging instances of server_template:"
find "$directory" | grep server_template;
grep -irl "server_template" $directory;
echo -e "DONE\n";

echo "Prepping git repository:"
cd $directory;
git init -b main;

echo -e "Making first commit:"
git add .;
git commit -m "bootstraping $name project from inceptus server_template";

echo -e "Setting remote git origin (ssh type):"
git remote add origin "git@$repository.git"
git remote -v;
echo -e "DONE\n";

echo -e "\n\n\n\n$delimeter\nREADME:\n"
cat README.md

echo -e "\n\n\n\nYour $name project is ready!\n"
echo "Please read README.md printed above ^"
echo "(especially the section about HTTPS for local development)"

