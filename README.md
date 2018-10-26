# Expect

expect is a Go package to auto input on terminal.


# Installation

```
go get github.com/standsun/expect
```


# Usage

```
// auto login 127.0.0.1
e, err := expect.New("ssh", "standsun@127.0.0.1")
if err != nil {
    panic(err)
}
defer e.Close()

// add output expect string and then input string

e.AddGroup(expect.NewGroup("Password:", "pass").Show(false))     // .Show(false) set password not show
e.AddGroup(expect.NewGroup("$", "cd /etc"))
e.AddGroup(expect.NewGroup("$", "ls"))

e.Run()
```

the execute result on mac

```
Password:•••••••
standsun@work:~ $ cd /etc
standsun@work:/etc $ ls

standsun@work:/etc $ ls
afpovertcp.cfg                        emond.d                               mail.rc~orig                          paths~orig                            rpc~previous
afpovertcp.cfg~orig                   find.codes                            man.conf                              periodic                              rtadvd.conf
...
standsun@work:/etc $
standsun@work:/etc $
```
