///https://stackoverflow.com/a/73703527

package main

import (
	//"time"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/data/binding"
	"database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
	"strconv"
)

func Delete(app fyne.App, student Student) {

}

func Edit(app fyne.App, student Student) {
	
	fmt.Println(strconv.Itoa(student.Id))
	
	add_input_field := widget.NewEntry()
	add_input_field.PlaceHolder = student.Name

	container := container.NewBorder(
		nil,nil,
		widget.NewLabel("ID: " + strconv.Itoa(student.Id)),
		widget.NewButton("Confirm", func() {}),
		add_input_field, 
	)

	window := app.NewWindow("Edit")
    window.Resize(fyne.NewSize(300,30))

    window.SetContent(container)

    window.Show()
}

func MainWindow(app fyne.App) fyne.Window {

	user_input_container := container.New(
		layout.NewVBoxLayout(),
		//widget.NewLabel("hello"),
		container.NewBorder(
			nil, nil, nil,
			widget.NewButton("Add", func() {}),
			widget.NewEntry(),
		),
		widget.NewLabel(""),
	)

	students := Fetch_all()

	students_bind_list := binding.NewUntypedList()
    for _,name := range students {
      fmt.Println(name)
      students_bind_list.Append(name)
    }

	table_list_container := container.NewBorder(
		layout.NewSpacer(),
		nil, nil, nil,
		widget.NewListWithData(
			students_bind_list,
			func() fyne.CanvasObject {
				return container.NewBorder(
					nil, nil, nil,
					// left of the border
					container.NewHBox(
        				widget.NewButton("Edit", func() {}),
        				widget.NewButton("Delete", func() {}),
    				),
					// takes the rest of the space
					widget.NewLabel(""),
				)
			},

			func(di binding.DataItem, o fyne.CanvasObject) {
				ctr, _ := o.(*fyne.Container)
				// ideally we should check `ok` for each one of those casting
				// but we know that they are those types for sure
				l := ctr.Objects[0].(*widget.Label)
				c1 := ctr.Objects[1].(*fyne.Container).Objects[0].(*widget.Button)
				//c2 := ctr.Objects[1].(*fyne.Container).Objects[1].(*widget.Button)
				/*
					diu, _ := di.(binding.Untyped).Get()
					todo := diu.(models.Todo)
				*/
				v, _ := di.(binding.Untyped).Get()
				student := v.(Student)
				fmt.Sprintf("%v",student.Name)
				l.SetText(student.Name)
				c1.OnTapped= func () {Edit(app,student)}
			},
		),
	)

	container := container.NewBorder(
		user_input_container,
		nil, nil, nil,
		table_list_container,
	)
	
    window := app.NewWindow("Alunos")
    //window.SetFullScreen(true)
    window.SetMaster()
    window.Resize(fyne.NewSize(500,500))

    window.SetContent(container)

	return window
}


////#################

type Student struct {
    Id int
    Name string
}



func Fetch_all() []Student {
    db, err := sql.Open("mysql", "teacher:System32@tcp(10.0.3.89:3306)/school")
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    stmtOut, err := db.Prepare("SELECT id,name FROM students")
    if err != nil {
        panic(err.Error())
    }
    defer stmtOut.Close()

    var students []Student
    rows, err := stmtOut.Query()
    if err != nil {
        panic(err.Error())
    }

    for rows.Next() {
        var std Student
        rows.Scan(&std.Id,&std.Name)
        students = append(students, std)
    }

    //students_bind_list := binding.NewUntypedList()
    //for _, t := range students {
    //  fmt.Println(t)
    //  students_bind_list.Append(t)
    //}    
    //return students_bind_list

    return students
}



////////////////////////////

func main() {
    app := app.New()
	MainWindow(app).Show()
    app.Run()
}

