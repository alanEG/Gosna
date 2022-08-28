package main

import (
	"bufio"
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Running() {
	_, errDiff := exec.LookPath("diff2html")
	_, errJsB  := exec.LookPath("js-beautify")
	if errDiff != nil && errJsB != nil {
		logger("error", "\n[command404]", "diff2html doesn't exist:\n\tnpm install diff2html\nor\nnpm i js-beautify", "", 1)
	}

	Banner()
	options_parse()
	configFile = func() string {
		if os.Getenv("gosna_config") != "" {
			return os.Getenv("gosna_config")
		} else {
			return configFileFlag
		}
	}()

	config_handling("load") //load config file

	setupFirst() //setup for run in first

	wg := &sync.WaitGroup{}

	if Run == "add" {
		Dynamic_pre(wg)

	} else if Run == "check" {
		if flagRepeat != "" {
			for {
				Check_pre(wg)
				config_handling("save")
				Handling_time(flagRepeat)
			}
		} else {
			Check_pre(wg)

		}

	} else { //exit if the run flag is null or not [check,add]
		fmt.Println("Enter Valid run type [check,add]")
		flag.Usage()
		os.Exit(0)

	}

	config_handling("save")
}

func Handling_time(Time string) {
	Type := string(Time[len(Time)-1])
	TimeSleep, _ := strconv.Atoi(Time[0 : len(Time)-1])
	switch Type {
	case "m":
		logger("success", "\n[Sleep]", "Sleep "+strconv.Itoa(TimeSleep)+" Minute\n", "", 0)
		time.Sleep(time.Duration(TimeSleep) * time.Minute)

	case "h":
		logger("success", "\n[Sleep]", "Sleep "+strconv.Itoa(TimeSleep)+" Hour\n", "", 0)
		time.Sleep(time.Duration(TimeSleep) * time.Hour)

	case "d":
		logger("success", "\n[Sleep]", "Sleep "+strconv.Itoa(TimeSleep)+" day\n", "", 0)
		time.Sleep(time.Duration(TimeSleep*24) * time.Hour)
	default:
		logger("error", "\n[Error]", "Repeat Option not valid", "Enter valid repeat option [m,h,d]", 1)
		os.Exit(0)
	}
}

func Dynamic_pre(wg *sync.WaitGroup) {
	w := make(chan string)
	go func() { //add for add new targets for check they later
		q := bufio.NewScanner(os.Stdin) //read from stdin
		for q.Scan() {
			w <- q.Text()
		}
		close(w)
	}()

	for i := 0; i < Thread; i++ {
		wg.Add(1)
		go Dynamic(w, wg, Header)

	}
	wg.Wait()
}

func Check_pre(wg *sync.WaitGroup) {
	c := make(chan Target)
	go func() { //check the target
		a := data.Target //Json target Object
		for _, b := range a {
			c <- b
		}
		close(c)

	}()

	for i := 0; i < Thread; i++ {
		wg.Add(1)
		go Check(c, wg)

	}

	wg.Wait()
	Diff()
}

//Check request if there is change in urls
func Check(targets chan Target, mainWg *sync.WaitGroup) {
	defer mainWg.Done()
	for target := range targets {

		response, err := send_request(
			target.Url, "GET", target.Headers, []byte{},
		)

		if err != nil {
			logger("error", "[Error]", target.Url, err, 0)
			break
		}

		rand, _, err := content_type_format(target.Url, response)

		if err != nil {
			logger("error", "[Error]", target.Url, err, 0)
			break
		}

		//if the dynamic status in the target is true
		//Remove the dynamic lines from new response
		//And save it in the target file

		if target.Dynamic.Status == true {
			DynmaicLine := strings.Split(
				strings.Trim(
					strings.Replace(
						fmt.Sprint(
							target.Dynamic.DynamicLine,
						), " ", ",", -1,
					), "[]",
				), ",",
			)
			remove_line_from_file("/tmp/"+rand, DynmaicLine, target.Url)
		}
		_, err = save_file("/tmp/"+rand, target.Url)

		if err != nil {
			logger("error", "[Error]", target.Url, err, 0)
			break
		}

		logger("success", "[Check]", target.Url, "", 0)
	}
}

func Dynamic(urls chan string, wa *sync.WaitGroup, headers map[string]string) {
	var wg sync.WaitGroup
	loop := 1
	errorDublicate := ""

	defer wa.Done()
	for url := range urls {
		for targerUrlIndex := range data.Target {
			if url == data.Target[targerUrlIndex].Url {
				errorDublicate = "err"
				break
			}
		}

		if errorDublicate != "err" {
			if flagDynmaic == true {
				loop = 5
			}

			var resp []*http.Response

			//Forloop for get 3 request for check the lines between them later
			for i := 0; i < loop; i++ {
				wg.Add(1)
				go func() {
					response, err := send_request(url, "GET", headers, []byte{})

					if err != nil {
						logger("error", "[Error]", url, err, 0)
						wg.Done()
						return
					}

					resp = append(resp, response)
					wg.Done()
				}()
			}

			wg.Wait()

			file_name, err := check_dynamic(url, resp, headers)

			if err != nil {
				for filen, _ := range file_name {
					os.Remove("/tmp/" + filen)
				}
				break
			}
		
			for filen, _ := range file_name {
				os.Remove("/tmp/" + filen)
			}

			if err != nil {
				logger("error", "[Error]", url, err, 0)
				break
			}
			logger("success", "[Add]", url, "", 0)
		
		} else {
			logger("error", "[Error]", url, "Dublicate Url", 0)
		}
	}
}

//Generator random string value for name tmep files
func random_string(length int) string {

	b := make([]byte, length)

	rand.Seed(time.Now().UnixNano())
	rand.Read(b)

	return fmt.Sprintf("%x", b)[:length]
}

//remove temp file and save the content to directory result
func save_file(file_old, url string) (string, error) {
	content, err := ioutil.ReadFile(file_old)

	file_new := fmt.Sprintf("%s%x", data.Config.Directory_work, md5.Sum([]byte(url)))

	err = ioutil.WriteFile(file_new, content, 0644)

	if err != nil {
		logger("error", "[Error]", url, err, 0)
		return file_new, err
	}

	err = os.Remove(file_old)

	if err != nil {
		logger("error", "[Error]", url, err, 0)
		return file_new, err
	}

	return file_new, nil
}

func append_to_json(url string, result []string, filename_save string, headers map[string]string, dynamic_status bool) {

	//convert string to int
	var target_data Target
	var result_int []int

	for _, item := range result {
		new_item, _ := strconv.Atoi(item)

		result_int = append(result_int, new_item)
	}

	target_data.Url = url
	target_data.Filename = filename_save
	target_data.Dynamic.Status = dynamic_status
	target_data.Dynamic.DynamicLine = result_int
	target_data.Headers = headers
	data.Target = append(data.Target, target_data)
}

func Execute_command(command string) string {
	cmd, err := exec.Command("/bin/bash", "-c", command).Output()
	if err != nil {
		fmt.Println(err)
	}

	return string(cmd)
}
