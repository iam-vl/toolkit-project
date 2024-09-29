

## Setup 

```
/toolkit-project$ ls
app  README.md  toolkit  toolkit.code-workspace
/toolkit-project$ go work init toolkit app
```

## Plan S01 

* Creating a random string 
* Writing a test 
* Trying things app w/ a simple app
* Pushing to Github 

Let's create a func: 
```go
const randStringSource = "abcdefghijklmnopqrstuvwxyzACDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"
type Tools struct{}
func (t *Tools) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}
```
Let's create a test (tools_test.go):
```go
func TestTools_RandomString(t *testing.T) {
	var testTools Tools
	s := testTools.RandomString(10)
	if len(s) != 10 {
		t.Error("wrong length random string returned")
	}
}
```
Let's run it: 
```sh
/toolkit-project$ cd toolkit
/toolkit-project/toolkit$ go test .
ok      github.com/iam-vl/toolkit       0.007s
```

GIt:
```
touch .gitignore
git init 
git add .
git config ...
git commit -am ...
```

## Plan S02

* Upload 1 or+ files from browser to server 
* Limit uploads by file size
* Limit uploads my mime types 
* Write a test 
* Write a simple app cd 

### Add new app

```
go work use app-upload
```