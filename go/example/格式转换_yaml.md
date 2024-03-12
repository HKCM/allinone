```go
package main
import (
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

type MarksStruct struct {
    Sceince     int8 `yaml:"science"`
    Mathematics int8 `yaml:"mathematics"`
    English     int8 `yaml:"english"`
}

type Student struct {
    Name   string      `yaml:"student-name"`
    Age    int8        `yaml:"student-age"`
    Marks  MarksStruct `yaml:"subject-marks"`
    Sports []string
}

func main() {
    s1 := Student{
        Name: "Sagar",
        Age:  23,
        Marks: MarksStruct{
            Sceince:     95,
            Mathematics: 90,
            English:     90,
        },
        Sports: []string{"Cricket", "Football"},
    }

    yamlData, err := yaml.Marshal(&s1)

    if err != nil {
        fmt.Printf("Error while Marshaling. %v", err)
    }

    fileName := "test.yaml"
    err = os.WriteFile(fileName, yamlData, 0644)
    if err != nil {
        panic("Unable to write data into the file")
    }
}
```