package main

import (
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

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
		cmd = "open -a Safari"
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
	//url := "https://forecast.weather.gov/MapClick.php?lat=41.25&lon=-81.44&unit=0&lg=english&FcstType=dwml"

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
	t := time.Now()
	formattedTime := t.Format(time.Kitchen)
	xdata = xdata + "Current Time : " + formattedTime + "<BR>"
	chr := ""
	ton := false
	word := ""
	loc := ""
	//------------------------------------------------------------------------ Location
	for x := 1; x < len(body); x++ {
		chr = string(body[x : x+1])
		if chr == "<" {
			ton = true
		}
		if chr == ">" {
			ton = false
			word = ""
		}
		if ton {
			word = word + chr
		}
		if word == "<location" {
			tmp := ""
			tdata := string(body[x+20 : x+170])
			for xx := 1; xx < len(tdata)-18; xx++ {
				if tdata[xx:xx+18] == "<area-description>" {
					xx = xx + 18
					for xx := xx; xx < len(tdata)-18; xx++ {
						chr = string(tdata[xx : xx+1])
						if chr == "<" {
							break
						}
						tmp = tmp + chr
					}

				}

			}
			loc = tmp
		}

	}
	xdata = xdata + "<BR>Location : " + loc + "<BR>"
	//------------------------------------------------------------------------ Hazzard Warning
	hazctl := false
	haz := ""
	for x := 1; x < len(body); x++ {
		chr = string(body[x : x+1])
		if chr == "<" {
			ton = true
		}
		if chr == ">" {
			ton = false
			word = ""
		}
		if ton {
			word = word + chr
		}
		if word == "<hazard headline" {
			haz = ""
			hazctl = true
			tdata := string(body[x : x+30])
			tt := false
			for xx := 1; xx < len(tdata); xx++ {
				chr = string(tdata[xx : xx+1])
				if tt {
					if asciistring.StringToASCII(chr) != 34 {
						haz = haz + chr
					}
				}
				switch {
				case asciistring.StringToASCII(chr) == 34 && tt == false:
					tt = true
				case asciistring.StringToASCII(chr) == 34 && tt == true:
					tt = false
				}
			}
		}

	}
	if hazctl {
		xdata = xdata + "<BR>Hazard Warning : " + haz + "<BR>"
	}

	//------------------------------------------------------------------------ Temperature
	for x := 1; x < len(body); x++ {
		chr = string(body[x : x+1])
		if chr == "<" {
			ton = true
		}
		if chr == ">" {
			ton = false
			word = ""
		}
		if ton {
			word = word + chr
		}
		if word == "<temperature" {
			if string(body[x+8:x+16]) == "apparent" {
				temp := ""
				tdata := string(body[x+20 : x+100])
				for xx := 1; xx < len(tdata)-7; xx++ {
					if tdata[xx:xx+7] == "<value>" {
						xx = xx + 7
						for xx := xx; xx < len(tdata)-7; xx++ {
							chr = string(tdata[xx : xx+1])
							if chr == "<" {
								break
							}
							temp = temp + chr
						}
					}
				}
				xdata = xdata + "<BR>Current Temperature : " + temp + "<BR>"
			}
		}

	}
	//------------------------------------------------------------------------ Current Conditions
	cond := ""
	for x := 1; x < len(body); x++ {
		chr = string(body[x : x+1])
		if chr == "<" {
			ton = true
		}
		if chr == ">" {
			ton = false
			word = ""
		}
		if ton {

			word = word + chr
		}

		if word == "<weather-conditions w" {
			cond = ""
			tdata := string(body[x+10 : x+50])
			tt := false
			for xx := 1; xx < len(tdata); xx++ {
				chr = string(tdata[xx : xx+1])
				if tt {
					if asciistring.StringToASCII(chr) != 34 {
						cond = cond + chr
					}
				}
				switch {
				case asciistring.StringToASCII(chr) == 34 && tt == false:
					tt = true
				case asciistring.StringToASCII(chr) == 34 && tt == true:
					tt = false
				}
			}
		}

	}
	xdata = xdata + "<BR>Current Conditions : " + cond + "<BR>"

	//------------------------------------------------------------------------ Wind
	gust := ""
	sust := ""
	for x := 1; x < len(body); x++ {
		chr = string(body[x : x+1])
		if chr == "<" {
			ton = true
		}
		if chr == ">" {
			ton = false
			word = ""
		}
		if ton {
			word = word + chr
		}
		if word == "<wind-speed" {
			if string(body[x+8:x+12]) == "gust" {
				tdata := string(body[x+20 : x+100])
				for xx := 1; xx < len(tdata)-7; xx++ {
					if tdata[xx:xx+7] == "<value>" {
						xx = xx + 7
						for xx := xx; xx < len(tdata)-7; xx++ {
							chr = string(tdata[xx : xx+1])
							if chr == "<" {
								break
							}
							gust = gust + chr
						}
					}
				}
			}
			if gust == "NA" {
				gust = ""
			}
			if string(body[x+8:x+17]) == "sustained" {
				tdata := string(body[x+20 : x+100])
				for xx := 1; xx < len(tdata)-7; xx++ {
					if tdata[xx:xx+7] == "<value>" {
						xx = xx + 7
						for xx := xx; xx < len(tdata)-7; xx++ {
							chr = string(tdata[xx : xx+1])
							if chr == "<" {
								break
							}
							sust = sust + chr
						}
					}
				}

			}

		}
	}
	xdata = xdata + "<BR>Sustained Wind at " + sust + " MPH"

	if len(gust) > 0 {
		xdata = xdata + " Gusting to " + gust + " MPH"
	}
	xdata = xdata + "<BR>"
	//------------------------------------------------------------------------ Humidity
	//------------------------------------------------------------------------

	xdata = xdata + "<BR><BR>RSS Feed Reader"

	xdata = xdata + "</center>"
	xdata = xdata + " </body>"
	xdata = xdata + " </html>"

	return xdata

}
