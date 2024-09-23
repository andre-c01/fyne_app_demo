///https://stackoverflow.com/a/73703527

package main

import (
	"time"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/data/binding"
	"database/sql"
    _ "fmt"
    _ "github.com/go-sql-driver/mysql"
	"strconv"
)

func Delete(student Student) {
	db, err := sql.Open("mysql", "teacher:System32@tcp(10.0.3.89:3306)/school")
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    stmtOut, err := db.Prepare("delete from students where id=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtOut.Close()

	_, err = stmtOut.Exec(student.Id)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

}

func Add(name string) {
	db, err := sql.Open("mysql", "teacher:System32@tcp(10.0.3.89:3306)/school")
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    stmtOut, err := db.Prepare("INSERT INTO students VALUES(default,?)")
    if err != nil {
        panic(err.Error())
    }
    defer stmtOut.Close()

    _, err = stmtOut.Exec(name)
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

}

func Edit(id int, name string) {
	
	db, err := sql.Open("mysql", "teacher:System32@tcp(10.0.3.89:3306)/school")
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    stmtOut, err := db.Prepare("update students set name = ? where id = ?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtOut.Close()

    _, err = stmtOut.Exec(name,id)
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

}

func EditWindow(student Student) {
	
	//fmt.Println(strconv.Itoa(student.Id))
	
	edit_input_field := widget.NewEntry()
	edit_input_field.Text = student.Name

	window := this_app.NewWindow("Edit")
    window.Resize(fyne.NewSize(300,30))

	container := container.NewBorder(
		nil,nil,
		widget.NewLabel("ID: " + strconv.Itoa(student.Id)),
		widget.NewButton("Confirm", func() {
			Edit(student.Id, edit_input_field.Text)
			window.Close()
		}),
		edit_input_field, 
	)

	window.SetContent(container)

    window.Show()
}

func MainContainer() *fyne.Container {

	user_add_field := widget.NewEntry()
	user_add_field.PlaceHolder = "Aluno"
	user_add_btn := widget.NewButton("Add", func() {
		Add(user_add_field.Text)
		user_add_field.Text = ""
	})
	user_add_btn.Disable()	
	user_add_field.OnChanged = func(s string) {
		user_add_btn.Disable()

		if len(s) >= 3 {
			user_add_btn.Enable()
		}
	}

	user_input_container := container.New(
		layout.NewVBoxLayout(),
		//widget.NewLabel("hello"),
		container.NewBorder(
			nil, nil, nil,
			user_add_btn,
			user_add_field,
		),
		widget.NewLabel(""),
	)

	students := Fetch_all()

	students_bind_list := binding.NewUntypedList()
    for _,name := range students {
      //fmt.Println(name)
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
				c2 := ctr.Objects[1].(*fyne.Container).Objects[1].(*widget.Button)
				/*
					diu, _ := di.(binding.Untyped).Get()
					todo := diu.(models.Todo)
				*/
				v, _ := di.(binding.Untyped).Get()
				student := v.(Student)
				//fmt.Sprintf("%v",student.Name)
				l.SetText(student.Name)
				c1.OnTapped= func () {EditWindow(student)}
				c2.OnTapped = func () {Delete(student)}
			},
		),
	)

	

	container := container.NewBorder(
		user_input_container,
		nil, nil, nil,
		table_list_container,
	)

	return container
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

var this_app fyne.App = app.New()

func main() { 

	window := this_app.NewWindow("Alunos")
    //window.SetFullScreen(true)
    window.SetMaster()
    window.Resize(fyne.NewSize(500,500))

    window.SetContent(MainContainer())

	go func() {
		for range time.Tick(time.Second * 3){
			window.SetContent(MainContainer())
		}
	}()
	
	window.Show()
	this_app.Run()
}

