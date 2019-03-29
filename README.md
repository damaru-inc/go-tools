# go-tools
A random collection of golang projects. So far they are pretty much only useful to the Cloud Integration team at Solace Inc.

## logtable
This takes a log file from the Solace PCF Service Broker app and creates a csv file that can be loaded into a spreadsheet, displaying the output of each thread in its own column. The timestamp is in the first column, the second column contains the 'scheduled-1' thread labelled 's-1', and each subsequent column contains one of the ServiceBrokerAsyncThread-xx threads, labelled t-xx.

To use it run

```go run logtable.go service-broker.log > out.csv```

and then load out.csv into a spreadsheet program.

## server
This is a mini http server. Currently it just does a redirect - it was meant to test that the http clients in the Solace PCF Service Sroker can handle redirects.


