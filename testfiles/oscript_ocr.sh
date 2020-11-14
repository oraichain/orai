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
    # collect input string with a delimiter of each data source value separated with a '-' delimiter
    IFS='-' read -ra array <<< "$2"
    echo $IFS
    aggregation_result=0
    size=0
    result=""
    # here is the algorithm for each oracle script. This should be different based on the oscript
    for i in "${array[@]}"; do
      let "size+=1"
    done
    echo "collected the following result from" $size "data sources that passed the test case": ${array[0]}
    #echo $2
  else
    echo 0
  fi
}

route_name="$(route $1 $2 $3 $4)"
echo $route_name