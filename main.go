package main

import (
    "fmt"
    "os"
    "sort"
    "time"

    "github.com/PuloV/ics-golang"
)

var zone time.Location

func main() {
    args := os.Args[1:]
    if len(args) != 1 {
        fmt.Println("ics 파일 경로를 입력해주세요.")
        fmt.Println("example: ics-viewer ./basic.ics")
        os.Exit(1)
    }

    parser := ics.New()
    input := parser.GetInputChan()
    input <- args[0]

    parser.Wait()

    cal, err := parser.GetCalendars()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if cal[0].String() != "" {
        fmt.Printf("%s \n\n", cal[0].String())
    }
    zone = cal[0].GetTimezone()

    eventsMap := cal[0].GetEventsByDates()
    keys := make([]string, len(eventsMap))
    for k := range eventsMap {
        keys = append(keys, k)
    }

    for i, k := range keys {
        fmt.Printf("%d: %s \n", i, k)
    }

    var idx int
    fmt.Print("어느 날짜의 이벤트를 보시겠습니까?")
    fmt.Scanln(&idx)

    eventViewer(eventsMap[keys[idx]])
}

func eventViewer(events []*ics.Event) {
    timeEvents := make(map[string]string)
    var eventsKeys []string
    for _, event := range events {
        start := event.GetStart().In(&zone)
        end := event.GetEnd().In(&zone)
        key := fmt.Sprintf("%02d/%02d, %02d:%02d~%02d:%02d", start.Month(),
                                                             start.Day(),
                                                             start.Hour(),
                                                             start.Minute(),
                                                             end.Hour(),
                                                             end.Minute())
        timeEvents[key] = event.GetSummary()
        eventsKeys = append(eventsKeys, key)
    }
    sort.Strings(eventsKeys)

    for _, k := range eventsKeys {
        fmt.Printf("[%s] %s \n", k, timeEvents[k])
    }
}
