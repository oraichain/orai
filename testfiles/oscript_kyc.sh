#!/bin/bash

route() {
  if [[ $1 = "aiDataSource" ]]
  then 
    echo "test test test" # return names of the data sources
  elif [[ $1 = "testcase" ]] # return names of the test cases
  then
    echo "test"
  elif [[ $1 = "aggregation" && $2 > $3 ]] # $2 is true output, $3 is expected output
  then 
    echo 1
  else
    echo 0
  fi
}

route_name="$(route $1 $2 $3)"
echo $route_name