set -ex

cd app1
make run &
P1=$!

cd ../app2
make run &
P2=$!

cd ../app3
make run &
P3=$!

wait $P1 $P2 $P3
