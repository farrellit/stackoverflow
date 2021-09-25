#!/bin/bash -e
sed 's/ .*//' < input.txt | sort | uniq | xargs -n 1 -I % \
  sh -c "grep ^% < input.txt | awk '{print \$2}' > %.txt"
