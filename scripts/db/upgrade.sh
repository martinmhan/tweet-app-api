# Use this script to run all the MongoDB scripts in order
# Structural changes to the database should only be done via this script
# Scripts should be labeled and numbered in order with the appropriate leading 0s
# Scripts should be idempotent, i.e., running all of them multiple times in order should not cause any problems.

set -e

source ../../.env;

for FILE in `ls *.js | sort -V`; do
  mongo $FILE --eval "const dbName='"$DB_NAME"'";
  echo "Ran $FILE";
done;

echo "Database updated successfully";

exit 0;