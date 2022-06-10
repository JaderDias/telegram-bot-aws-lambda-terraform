#!/bin/bash

language="$1"
temp_memory_file_name="/dev/shm/$language.csv"

i=0
# read 100 lines at a time from the csv file
while mapfile -t -n 100 ary && ((${#ary[@]})); do
    printf '%s\n' "${ary[@]}" > $temp_memory_file_name
    # and upload those chunks to AWS S3
    aws s3api put-object --bucket my-bucket-legible-quetzal --key language/$language/$i --body $temp_memory_file_name
    ((i=i+1))
done <$language.csv
