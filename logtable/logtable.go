package main

import "bufio"
import "fmt"
import "log"
import "os"
import "regexp"
import "sort"
import "strconv"
import "strings"

// N.B. Our log level is surrounded by \033[ for highlighting.
var logLineRegex = regexp.MustCompile(`^\s*(\S*) \[.+\].+\033\[.+\033\[[^\[]+\[([^\]]+)\] (.*)$`)
var commaRegex = regexp.MustCompile(`,`)
var threadRegex = regexp.MustCompile(`[^-]+-(\d+)`)
var threads = make(map[string]string)
var threadToCol = make(map[string]int)
var colToThread = make(map[int]string)

func main() {
//    testRegex()
//    t := mungeThread("ServiceBrokerAsyncThread-6")
//    fmt.Println(t)
//    return
    if (len(os.Args) != 2) {
        fmt.Printf("Usage: logtable filename\n")
        os.Exit(1)
    }
	filename := os.Args[1]
	//fmt.Printf("filename: %v\n", filename)
    scanFile(filename)
}

func scanFile(filename string)  {
	//fmt.Printf("scanFile: %v\n", filename)
    f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
    if err != nil {
        log.Fatalf("open file error: %v", err)
    }
    defer f.Close()

    sc := bufio.NewScanner(f)
    for sc.Scan() {
        line := sc.Text()  // GET the line string
        findThreads(line)
    }

    var keys []string
    for k := range threads {
        keys = append(keys, k)
    }

    sort.Strings(keys)

    for i := range keys {
       s := keys[i]
       t := threads[s]
       threadToCol[t] = i
       colToThread[i] = t
       i++
    }

    //fmt.Println(threadToCol)
    //fmt.Println(colToThread)

    printFirstLine()
    f.Seek(0,0)

    sc = bufio.NewScanner(f)
    for sc.Scan() {
        line := sc.Text()  // GET the line string
        scanLine(line)
    }
    if err := sc.Err(); err != nil {
        log.Fatalf("scan file error: %v", err)
    }

    //fmt.Println(threads)
}

func testRegex() {
	//str := "  2019-03-26T17:38:46.10-0400 [APP/PROC/WEB/0] OUT 2019-03-26 21:38:46,099 INFO  [ServiceBrokerAsyncThread-19] com."
    str := "2019-03-26T17:38:46.18-0400 [APP/PROC/WEB/0] OUT 2019-03-26 21:38:46,185 INFO  [ServiceBrokerAsyncThread-20] com.sola.clou.serv.serv.HAServicePlanConfigurationHandler - isConfigured requestedMessageRouter: MessageRouter semp: 192.168.16.118 existingMessageRouter: MessageRouter semp: 192.168.16.118"

	re := regexp.MustCompile(`^\s*(\S*) \[.+\][^\[]+\[(.+)\] (.*)$`)
	matches := re.FindStringSubmatch(str)
	if (len(matches) > 0) {
        fmt.Printf("Matches: %v|%v|%v\n", matches[1], matches[2], matches[3])
    } else {
	    fmt.Printf("Doesn't match: %v\n", str)
    }

    str2 := "This contains, a comma."
    str3 := commaRegex.ReplaceAllString(str2, " ")
    fmt.Println(str3)

    os.Exit(0)
}

func findThreads(line string) {
    matches := logLineRegex.FindStringSubmatch(line)
	if (len(matches) > 2) {
        thread := matches[2]
        thread = mungeThread(thread)
        if (thread != "") {
            threads[thread] = thread
        }
        //fmt.Println("Thread", thread)
   } 
}

func printFirstLine() {
    str := "Datetime,"
    ln := len(threads)
    for i := 0; i < ln; i++ {
       thread :=  colToThread[i]
       str += thread

       if (i < ln - 1) {
           str += ","
       }
    }

    fmt.Println(str)
}

func scanLine(line string) {
    matches := logLineRegex.FindStringSubmatch(line)

	if (len(matches) > 2) {
        date := matches[1]
        thread := matches[2]
        message := matches[3]
        thread = mungeThread(thread)

        if (thread != "") {
            col := threadToCol[thread] + 1
            colsBefore := col
            colsAfter := len(threads) - col
            commasBefore := strings.Repeat(",", colsBefore)
            commasAfter := strings.Repeat(",", colsAfter)
            message = commaRegex.ReplaceAllString(message, " ")
            fmt.Printf("%s%s%s%s\n", date, commasBefore, message, commasAfter)
        }
        //fmt.Printf("Matches: %v|%v|%v\n", matches[1], matches[2], matches[3])
    } else {
	    //fmt.Printf("Doesn't match:  %d %v\n", len(matches), line)
    }
}


func convertTo2Digits(s string) (string) {
    num, err := strconv.Atoi(s)
    if (err != nil) {
        log.Fatalf("Can't convert string to int: %v %v", num, err)
    }
    str := fmt.Sprintf ("%02d", num)
    return str
}

func mungeThread(thread string) (string) {
    if (strings.HasPrefix(thread, "scheduling")) {
        thread = compressThread("s", thread)
    } else if strings.HasPrefix(thread, "ServiceBrokerAsyncThread") {
        thread = compressThread("t", thread)
    } else {
        //fmt.Fprintf(os.Stderr, "Rejecting thread %s\n", thread)
        return ""
    }
    return thread
}

func compressThread(prefix, thread string) (string) {
    
	matches := threadRegex.FindStringSubmatch(thread)

    if (len(matches) < 2) {
        log.Fatalf("Can't munge thread %s", thread)
    }

    num := matches[1]
    num =  convertTo2Digits(num)
    ret := fmt.Sprintf("%s-%s", prefix, num)
    return ret
}

