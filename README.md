# CTOP ("Top for Cassandra")

CTOP is a tool which allows you to quickly find out what's happening on a machine running Cassandra. It is particularly useful on a cluster with multiple-tenants, multiple-applications, and large numbers of tables. If you suspect that the performance is not good, then you can use this to figure out which table is giving you trouble.

## Screenshots:
![Ordered by read-latency](ordered-by-read-latency.png "Ordered by read-latency")
![Ordered by disk-space](ordered-by-diskspace.png "Ordered by disk-space")

## How to use it:
* Run the binary
* Allow CTOP some time to collect metrics, then press SPACE to refresh the display. You will see your tables listed by Reads/s (in descending order).
* You can press SPACE at any time to refresh the display
* The numbers 1 through 5 change the sorting order:
  1: Order by Reads/s
  2: Order by Writes/s
  3: Order by the amount of disk-space used (in Bytes)
  4: Order by read-latency (in miliseconds)
  5: Order by write-latency (in miliseconds)
* Pressing "Q" will quit
* Metrics are for one node only (not cluster-wide)

## Notes
* CTOP should run on anything, but I've only tested it on Ubuntu.

## Pre-requisites:
* CTOP uses "libmx4j-java" (an HTTP -> JMX) to retrieve JMX metrics from Cassandra (there was no easy way to query JMX directly), so this needs to be in the class-path where Cassandra can find it upon startup (with Ubuntu it is enough to do "apt-get install libmx4j-java", then to re-start Cassandra).
* Cassandra needs to be told how to run MX4J (the default is to listen to the interface that Cassandra listens on):  
```
  MX4J_ADDRESS="-Dmx4jaddress=127.0.0.1"  
  MX4J_PORT="-Dmx4jport=8081"  
  JVM_OPTS="$JVM_OPTS $MX4J_ADDRESS"  
  JVM_OPTS="$JVM_OPTS $MX4J_PORT"  
```
