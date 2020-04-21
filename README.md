# doulivery-go

Installation
------------

Use go get.

	go get github.com/cloudoti/doulivery-go

Then import the validator package into your own code.

	import "github.com/cloudoti/doulivery-go"
	
Usage and documentation
------

Send email without files

```go
client := CreateClient("app_id", "key", "secret")

//Here, you can use GMAIL or MAILGUN
mailer := CreateMailer(client, "GMAIL").
    AddFrom("daniel.roncal.87@gmail.com").
    AddTo("daniel.roncal@doous.com").
    AddSubject("Test Doulivery").
    AddBody("Hi Test")

err := mailer.SendEmail()

if err != nil {
    t.Errorf("error test: " + err.Error())
}
 ```

Send email with files

```go
client := CreateClient("app_id", "key", "secret")

b, _ := ioutil.ReadFile("file_path")

//Here, you can use GMAIL or MAILGUN
mailer := CreateMailer(client, "GMAIL").
    AddFrom("daniel.roncal.87@gmail.com").
    AddTo("daniel.roncal@doous.com").
    AddSubject("Test Doulivery file").
    AddBody("Hola Prueba file").
    AddFile(File{Name: "README.md", Content: b})

err := mailer.SendEmail()

if err != nil {
    t.Errorf("error test: " + err.Error())
}
 ```