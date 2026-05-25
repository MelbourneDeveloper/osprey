#!/bin/bash
while read p || [ -n "$p" ]
do
    # Escape forward slashes for sed
    escaped_p=$(echo "$p" | sed 's/\//\\\//g')
    sed -i "/${escaped_p}/d" ./coverage.out
done < ./exclude-from-code-coverage.txt