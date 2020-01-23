package main

import (
	"log"
	"net/http"
	"os/exec"
	"time"
	"fmt"
	"regexp"
)

func main() {
/*
    This go server listens at  http://localhost:8888/apk.
    You will need to add a rule to your .htaccess file like this:
    RewriteRule mediaplayer.apk http://localhost:8888/apk [P]
    This allows you to visit http://yourdomain.com/mediaplayer.apk
    and get back the apk file as "mediaplayer.apk", rather then
    whatever it is called on your server.
    
    You can probably use *.apk in your .htaccess so that any name works.
    
    You need to install go, then you can run:
    go install $GOPATH/apkserver.go (Make the binary file. Make sure apkserver.go is in your $GOPATH folder.)
    $GOBIN/apkserver (Launch the go server.)
    
    Of course, you don't need to use go. If you have issues with it
    you can use php, nodejs or any other server side language.
*/

	http.HandleFunc("/apk", func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		//log.Println("in the matrix")
        //Get clickid from the URL. eg: http://domain.com/mediaplayer.apk?cid=clickid
		keys, ok := r.URL.Query()["cid"]

		if !ok || len(keys[0]) < 1 {
			//log.Println("No cid dude...")
			return
		}

		//cid found
		log.Println(keys[0])

		//Make sure cid does not contain any special chars. Don't want to get hacked.
		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			return
		}
		cid := reg.ReplaceAllString(keys[0], "")

        //Shell script location + clickid
		cmdStr:="/root/Desktop/apk2.sh " + cid

		//Execute shell script to generate apk
		cmd := exec.Command("/bin/sh", "-c", cmdStr)
		//if err != nil {
		//	log.Println(err)
		//}

		cmd.Run()

        //Serve the apk. Replace this with wherever your apk is.
        http.ServeFile(w, r, "/root/Android/src/droidware/app/target/outputs/apk/release/signed_"+cid+".apk")

        //Checking how long it takes to serve the file. Not sure if this is accurate or not.
		t2 := t1.Add(time.Second)
		diff := t2.Sub(t1)
		fmt.Println(diff)

		//Cleanup. Delete clickid folder and new apk file to prevent filling up our disk.
        //Replace this with the proper paths for your setup.
		cmdStr = "rm -r /root/Android/src/droidware/app/target/outputs/apk/release/"+cid
                cmd = exec.Command("sh", "-c",  cmdStr)
                cmd.Run()

		cmdStr = "rm -r /root/Android/src/droidware/app/target/outputs/apk/release/signed_"+cid+".apk"
                cmd = exec.Command("sh", "-c",  cmdStr)
                cmd.Run()
	})
	http.ListenAndServe(":8888", nil)
}