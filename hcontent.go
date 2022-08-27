package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"
)

//Handling Response body
func content_type_format(url string, resp *http.Response) (string, int, error) {
	//Get content-type header
	var contentType string
	//Check if the header has value or None
	if len(resp.Header["Content-Type"]) != 0 {
		contentType = strings.Split(resp.Header["Content-Type"][0], ";")[0]
	}else {
		contentType = ""
	}

	//get extension from url
	extion := strings.ReplaceAll(regexp.MustCompile(`\.([a-zA-Z0-9]+)$`).FindString(url), ".", "")

	rand := random_string(10)

	filename := "/tmp/" + rand

	Body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", 0, err
	}

	err = ioutil.WriteFile(filename, []byte(Body), 0644)

	if err != nil {
		return "", 0, err
	}

	//Check the content-type and extension
	//And return beautify content
	//And overwrite in the same file
	if contentType == "text/javascript" || extion == ".js" {
		return rand,
			len(
				Execute_command(
					"cat " + filename + " | js-beautify | sponge " + filename + ";cat " + filename,
				),
			), nil

	} else {
		return rand,
			len(
				Execute_command(
					"cat " + filename + " | js-beautify --type html | sponge " + filename + ";cat " + filename,
				),
			), nil
	}
}

func check_dynamic(url string, resps []*http.Response, headers map[string]string) (map[string]int, error) {
	var (
		file_name     = make(map[string]int)
		result        []string
		filename_save string
		save          string
		err           error
	)

	//for loop response
	for _, resp := range resps {
		rand, length, err := content_type_format(url, resp)

		file_name[rand] = length

		if err != nil {
			return file_name, err
		}
	}

	content_length_first_file := file_name[reflect.ValueOf(file_name).MapKeys()[0].Interface().(string)]
	content_length_last_file := file_name[reflect.ValueOf(file_name).MapKeys()[len(file_name)-1].Interface().(string)]

	//check if the first file is equal the last file
	//If true then it is dynamic handling it
	//If false handling it as normal
	if content_length_first_file != content_length_last_file {
		//check each file with next file
		for filename, _ := range file_name {

			//Get the bigest file size and save it in result directory
			if save == "" {

				number := 0
				filenamed := ""

				//get the big response length
				for filenamec, lengthc := range file_name {
					if lengthc > number {
						number = lengthc
						filenamed = filenamec
					}
				}

				save = "done"
				filename_save, err = save_file("/tmp/"+filenamed, url)

				if err != nil {
					return file_name, err
				}

			}

			//append all length in one array and add in variable result

			for filename2, _ := range file_name {
				//Do not append the same file
				if filename != filename2 {
					result = append(result, get_diff_line_len(filename, filename2)...)
				}

			}
		}

		//sort the result array
		result = unique(result)

		remove_line_from_file(filename_save, result, url)
		append_to_json(url, result, filename_save, headers, true)
		return file_name, nil
		//Here if the responses have no dynamic lines
	} else {

		//get the first file name
		filename_save = reflect.ValueOf(file_name).MapKeys()[0].Interface().(string)
		//save in result folder
		filename_save, err = save_file("/tmp/"+filename_save, url)
		if err != nil {
			return file_name, err
		}

		append_to_json(url, result, filename_save, headers, false)
	}
	return file_name, nil
}

//Get the dynamic lines from response
func get_diff_line_len(file_old, file_new string) []string {

	var argument string
	var out_command string

	argument = "diff --unchanged-line-format='' --old-line-format='%dn,' --new-line-format='%dn,'"

	//Get the diffrent line between two files
	out_command = Execute_command(
		argument + fmt.Sprintf(
			" /tmp/%s /tmp/%s |  sed 's/,$//g'",
			file_old, file_new,
		),
	)

	//split diff lines and remove dublicate
	split_regex := unique(regexp.MustCompile(",").Split(string(out_command), -1))

	return split_regex
}

//Remove dublicate diff value
func unique(s []string) []string {

	inResult := make(map[string]bool)

	var result []string

	for _, str := range s {

		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)

		}

	}

	return result
}

//Remove lines from content
func remove_line_from_file(filename string, lines_number []string, url string) {
	split_arr := regexp.MustCompile(`^(\|)|(\|)$`).ReplaceAllString(strings.Join(lines_number, "|"), "")
	//Remove the dynamic lines
	Execute_command(
		fmt.Sprintf(
			"cat -n %s | sed -E 's/^( )+(%s)(\t)(.*)/---removed---/g;s/^( )+([0-9]+)(\t)//g' | sponge %s",
			filename, split_arr, filename,
		),
	)
}

func getUrlByFileName(filenames []string) []string {
	fileUrls := make(map[string]string)

	for _, target := range data.Target {
		filename := strings.Split(target.Filename, "/")
		fileUrls[filename[len(filename)-1]] = target.Url

	}

	for fileElemeNum, fileName := range filenames {
		filenames[fileElemeNum] = fileUrls[fileName]
	}

	return filenames
}

func Diff() {
	var files_has_change []string

	cmd := Execute_command("cd " + data.Config.Directory_work + ";git diff --name-only")
	split_regex := regexp.MustCompile(`\r?\n`).Split(string(cmd), -1)

	//check if there line null
	for _, change_file := range split_regex {
		if string(change_file) != "" {
			//append filename has change to file_has_change array
			files_has_change = append(files_has_change, strings.ReplaceAll(change_file, " ", ""))
		}

	}

	time_file := fmt.Sprintf("%s%v.html", data.Config.Directory_result, time.Now().Unix())
	if len(files_has_change) < 1 {
		fmt.Println("There is no change")

	} else {
		Execute_command(
			fmt.Sprintf("cd %s ;diff2html -d char --summary open -o stdout . > %v",
				data.Config.Directory_work,
				time_file,
			),
		)
		if data.Config.Channel_use != "None" {
			push_notifcation(time_file, files_has_change)
			fmt.Println(fmt.Sprintf("There is change in this file [%s]\n%s", time_file, strings.Join(getUrlByFileName(files_has_change), "\n")))
		} else {
			fmt.Println(fmt.Sprintf("There is change in this file [%s]\n%s", time_file, strings.Join(getUrlByFileName(files_has_change), "\n")))
		}
	}
}
