#!/bin/bash

# the number of parameters are fixed depending on the total data sources run
route() {
  if [[ $1 = "aiDataSource" ]]
  then 
    echo "image_ocr image_ocr image_ocr image_ocr" # return names of the data sources
  elif [[ $1 = "testcase" ]] # return names of the test cases
  then
    echo "testcase_ocr"
  elif [[ $1 = "aggregation" ]] # $2 is true output, $3 is expected output
  then
    echo "collected the following result from" $2 "data sources that passed the test case": $3
    #echo $2
  else
    echo 0
  fi
}

route_name="$(route $1 $2 $3 $4)"
echo $route_name