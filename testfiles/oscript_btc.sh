#!/bin/bash

# the number of parameters are fixed depending on the total data sources run
route() {
  if [[ $1 = "aiDataSource" ]]
  then 
    echo "crypto_compare_btc coindesk_btc coingecko_btc" # return names of the data sources
  elif [[ $1 = "testcase" ]] # return names of the test cases
  then
    echo "testcase_price"
  elif [[ $1 = "aggregation" ]] # $2 is true output, $3 is expected output
  then
    # collect input string with a delimiter of each data source value separated with a '-' delimiter
    IFS='-' read -ra array <<< "$2"
    aggregation_result=0
    size=0
    # here is the algorithm for each oracle script. This should be different based on the oscript
    for i in "${array[@]}"; do
      let "size+=1"
      # process "$i"
      aggregation_result=`echo $aggregation_result + $i | bc`
      #$aggregation_result = $aggregation_result + $i
    done
    temp=$aggregation_result
    # scale=2 allows division with floating points (two decimals)
    aggregation_result=$(echo "scale=2; ($temp) / $size" |bc) # |bc allows adding two floating point numbers
    echo "$aggregation_result"
  else
    echo 0
  fi
}

route_name="$(route $1 $2 $3 $4)"
echo $route_name