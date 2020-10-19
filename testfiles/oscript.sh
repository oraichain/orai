#!/bin/bash

# add #!bin/sh to allow this file to be called using shell script 

aggregation () {
  if [ $# != 2 ]; then
    echo "2 arguments are required "
    exit
  else
    local func_result=$(echo "$1 + $2" |bc) # |bc allows adding two floating point numbers
    echo "$func_result"
  fi
}

func_result="$(aggregation $1 $2)"
echo $func_result

# aggregation() {
#     local first_source=$1
#     local sec_source=$2
#     echo $(( $1 + $2 )) >&2
# } 

# result=$(aggregation)