Table library converts []string to a nice looking text-table.

Default Table:
```golang
import "github.com/thibran/table"

t, _ := table.New(
    []string{"Fruits:", "Count:"},
    [][]string{
        {"Apple", "4"},
        {"Banana", "25"}}...)
fmt.Println(t)
```

Default result:
```
Fruits: Count:
==============
Apple   4     
Banana  25    
```

The table object contains two Style objects for the head and the body.  
To e.g. add vertical lines to the body set `t.BodyStyle.vLines = true`.

Default with vertical body lines:
```
 Fruits:  Count:
=================
|Apple   |4     |
|Banana  |25    |
```

You can customize or pick a pre-defined style.

```golang
import "github.com/thibran/table"

func TestTable_foo(t *testing.T) {
	t, _ := table.New(
		[]string{"Fruits:", "Count:"},
		[][]string{
			{"Apple", "4"},
			{"Banana", "25"}}...)
	t.HeadStyle = table.NewStyle('o', '=', '|', true)
	t.BodyStyle = table.NewStyle('•', '–', '|', true)
	t.HeadOnlyBottomLine = false
	t.BottomLine = true
	fmt.Println(t.String())
}
```

Custom style result:
```
o========o======o
|Fruits: |Count:|
o========o======o
|Apple   |4     |
•––––––––•––––––•
|Banana  |25    |
•––––––––•––––––•
```

To get even more control over the table-drawing, use the `Draw` object from the package `github.com/thibran/table/draw`.
