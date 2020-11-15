#!/bin/bash

func_result="`/bin/bash .oraifiles/"$1" "$2"`"
# $3 is the expected output
expected_output=$(echo "$3" |base64 --decode)
if [[ $func_result = $expected_output ]]
then 
  echo $func_result
else
  echo null
fi