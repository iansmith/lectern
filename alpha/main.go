package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/coocood/qbs"
	"github.com/coreos/go-etcd/etcd"
	_ "github.com/lib/pq"
)

const (
	//unclear if we should really link against beta
	USERPROP = "postgres/host_count/username"
	PWDPROP  = "postgres/host_count/password"
)

func etcdConfig() []string {
	return []string{"http://" + os.Getenv("ETCD_HOST") + ":" + os.Getenv("ETCD_PORT")}
}

type HostCount struct {
	Hostname string `qbs:"pk"`
	Count    int
}

func ReadKV(name string) (string, error) {
	client := etcd.NewClient(etcdConfig())
	resp, err := client.Get(name, false, false)
	if err != nil {
		//special case not found
		if err.(*etcd.EtcdError).ErrorCode == 100 {
			return "", nil
		}
		return "", err
	}
	return resp.Node.Value, nil
}

func tryPostgres(user string, pwd string) (*sql.DB, error) {
	app := "alpha"
	host := os.Getenv("POSTGRES_HOST")
	url := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pwd, host, app)
	log.Printf("trying postgres url: %s", url)
	return sql.Open("postgres", url)
}

func readConfig() error {
	//try to contact the DB
	resp, err := ReadKV(USERPROP)
	if err != nil {
		return err
	}
	if resp == "" {
		return fmt.Errorf("no configuration found for database!")
	}
	user := resp
	resp, err = ReadKV(PWDPROP)
	if err != nil {
		return err
	}
	if resp == "" {
		return fmt.Errorf("can't find a password, but got a username")
	}
	pwd := resp

	log.Printf("read the username and password for db from consul: %s, %s", user, pwd)

	db, err := tryPostgres(user, pwd)
	if err != nil {
		return fmt.Errorf("failed to connect to the database (networking failed): %v", err)
	}
	qbs.RegisterWithDb("postgres", db, qbs.NewPostgres())

	//at this point we will only have failed if the connectivity is bad
	//not anything with auth because postgres doesn't try that until
	//you do sql
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	if err := readConfig(); err != nil {
		fmt.Fprintf(w, "read config: %v", err)
		return
	}
	q, err := qbs.GetQbs()
	if err != nil {
		fmt.Fprintf(w, "GetQbs: %v", err)
	}
	defer q.Close()

	h, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "hostname() failed: %v", err)

	}

	var count HostCount
	count.Hostname = h
	if err := q.Find(&count); err != nil {
		if err != sql.ErrNoRows {
			fmt.Fprintf(w, "find failed (probably your crendentials are bad): %v", err)
			return
		}
	}
	//why doesn't this happen?
	if err == sql.ErrNoRows {
		fmt.Fprintf(w, "no count found for %s", h)
	}
	count.Count++
	fmt.Fprintf(w, "new count for %s is %d", count.Hostname, count.Count)

	if _, err := q.Save(&count); err != nil {
		fmt.Fprintf(w, "save failed: %v", err)
		return
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":80", nil)
}
