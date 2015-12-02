Version: 0.1

Table library converts `[][]string` or `io.Reader` to a nice text-table.

Default example:
```golang
import "github.com/thibran/table"

ta, _ := table.New(true, [][]string{
        {"Fruits:", "Count:"},
        {"Apple", "4"},
        {"Banana", "25"}}...)
fmt.Println(ta)
```
or
```golang
import "github.com/thibran/table"

text := "Fruits:,Count:\nApple,4\nBanana,25\n"
r := strings.NewReader(text)
ta, _ := table.ReadFrom(r, true, []int{7, 6})
fmt.Println(ta)
```

Result:
```
Fruits: Count:
===============
Apple   4      
Banana  25     
```

The table contains two Style objects for the head and the body.  
To e.g. add vertical lines to the body set `t.BodyStyle.vLines = true`.

Default with vertical body lines:
```
 Fruits:  Count:  
==================
|Apple   |4      |
|Banana  |25     |
```

You can customize or pick a pre-defined style.

```golang
import "github.com/thibran/table"

ta, _ := table.New(true, [][]string{
    {"Fruits:", "Count:"},
    {"Apple", "4"},
    {"Banana", "25"}}...)
ta.HeadStyle = table.NewStyle('o', '=', '|', true)
ta.BodyStyle = table.NewStyle('•', '–', '|', true)
ta.HeadOnlyBottomLine = false
ta.BottomLine = true
fmt.Println(ta)
```

Custom style result:
```
o========o=======o
|Fruits: |Count: |
o========o=======o
|Apple   |4      |
•––––––––•–––––––•
|Banana  |25     |
•––––––––•–––––––•
```

To get even more control over the table-drawing, use the `Draw` object from the package `github.com/thibran/table/draw`.
