#!/bin/bash

orig_repository="github.com/robitx/inceptus/server_template"

usage()
{
echo """Hey,
this script expects you already have existing git repository, into
which you want to bootstrap server based on inceptus server_template.

Usage:
    bootstrap.sh -h Display this help message.

    bootstrap.sh -n PROJECT_NAME -r REPOSITORY -d DIRECTORY
    where:
    - PROJECT_NAME will be used instead of \"server_template\"
    - REPOSITORY (like github.com/XXX/YYY) will be used instead of:
      \"$orig_repository\"
    - DIRECTORY is a folder where you want to boostrap the server:
      ├── cmd
      │   └── server
      │       └── main.go
      ├── conf
      │   ├── PROJECT_NAME.env
      │   └── PROJECT_NAME.yaml
      └── internal
          ├── do
          │   └── do.go
          └── env
              ├── config.go
              └── environment.go
      ...
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

echo -e "Copying files from server_template/* to $directory:"
cp -v -r ./server_template/* "$directory";
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

echo "Server is ready:"
tree $directory;