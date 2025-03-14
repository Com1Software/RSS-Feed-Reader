package main

import (
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"fmt"

	asciistring "github.com/Com1Software/Go-ASCII-String-Package"
)

// ----------------------------------------------------------------
func main() {
	fmt.Println("Video Web Server")
	fmt.Printf("Operating System : %s\n", runtime.GOOS)
	xip := fmt.Sprintf("%s", GetOutboundIP())
	port := "8080"
	fmt.Println("Server running....")
	fmt.Println("Listening on " + xip + ":" + port)

	fmt.Println("")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		xdata := InitPage(xip)
		fmt.Fprint(w, xdata)
	})
	Openbrowser(xip + ":" + port)
	if err := http.ListenAndServe(xip+":"+port, nil); err != nil {
		panic(err)
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// Openbrowser : Opens default web browser to specified url
func Openbrowser(url string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start msedge"}

	case "linux":
		cmd = "chromium-browser"
		args = []string{""}

	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func InitPage(xip string) string {
	//---------------------------------------------------------------------------
	//----------------------------------------------------------------------------
	xdata := "<!DOCTYPE html>"
	xdata = xdata + "<html>"
	xdata = xdata + "<head>"
	//------------------------------------------------------------------------
	xdata = xdata + "<title>RSS Feed Reader</title>"
	//------------------------------------------------------------------------
	xdata = xdata + "</head>"
	//------------------------------------------------------------------------

	xdata = xdata + "<body>"
	xdata = xdata + "<center>"
	xdata = xdata + "<H1>RSS Feed Reader</H1>"
	//---------
	xdata = xdata + "<body>"

	xdata = xdata + "Host Port IP : " + xip
	xdata = xdata + "<BR><BR>"

	url := "https://forecast.weather.gov/MapClick.php?lat=41.5&lon=-81.7&unit=0&lg=english&FcstType=dwml"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching URL: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "HTTP error: %v\n", resp.Status)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
		os.Exit(1)
	}

	//	fmt.Println(string(body))
	fmt.Printf("\n\n len %d\n", len(body))
	xdata = xdata + string(body)
	chr := ""
	line := ""
	linecnt := 0
	for x := 1; x < len(body); x++ {
		chr = string(body[x : x+1])
		if asciistring.StringToASCII(chr) == 10 {
			fmt.Println(line)
			line = ""
			linecnt++
		} else {
			line = line + chr
		}

		//		fmt.Println(asciistring.StringToASCII(chr))
		//		fmt.Println(chr)

	}
	fmt.Println(linecnt)
	xdata = xdata + "<BR><BR>RSS Feed Reader"

	//------------------------------------------------------------------------

	xdata = xdata + "</center>"
	xdata = xdata + " </body>"
	xdata = xdata + " </html>"
	return xdata
}
