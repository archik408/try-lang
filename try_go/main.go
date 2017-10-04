package main


func main() {
	a := App{}
	a.Initialize("user", "pswd", "db")
	a.Run(":8080")
}