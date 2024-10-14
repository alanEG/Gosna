# Gosna

A Mointor url for change 

- [Requirements](https://github.com/alanEG/Gosna#Requirements)
- [Installation](https://github.com/alanEG/Gosna#installation)
- [Configuration file](https://github.com/alanEG/Gosna#configuration-file)
- [Example usage](https://github.com/alanEG/Gosna#example-usage)
    - [Normal_Usage](https://github.com/alanEG/Gosna#normal-Usage)
    - [Dynamic_Usage](https://github.com/alanEG/Gosna#dynamic-Usage)

## Requirements 
   
   - [diff2html-cli](https://github.com/rtfpessoa/diff2html-cli)
   - [js-beautify](https://www.npmjs.com/package/js-beautify)
    
## Installation

- Run the following command to install the tool `go install github.com/alanEG/Gosna` (the same command works for updating)
  _or_
- `git clone https://github.com/alaneg/gosna ; cd gosna ; go build`

Gosna depends on go1.18 or greater.

## Configuration file

Gosna won't work if you don't initialize the config file.  

Gosna will search for the config file in three different ways.

1 - `gosna_config` env

2 - `~/.gosna_config.json`

3 - `-config` option

Here is an example of the configurations 

```json
{
  "Config": {
    "Is_first": false,
    "Directory_work": "/home/user/work/",
    "Directory_result": "/home/user/result/",
    "Channel_use": "None",
    "Channel": {
      "Slack": null,
      "Mail": {
        "Tls": true,
        "From": null,
        "To": null,
        "Host": null,
        "Port": null,
        "Email": null,
        "Password": null
      }
    }
  },
  "Target": []
}
```

`Is_first`: It helps Gosna with important setup tasks during its initial setup. If not, set it to true.

`Directory_work`: Specify the folder where the tool will save the content from the URLs.

`Directory_result`: Specify the directory where the tool will save the results.

`Channel_use`: Specify the channel the tool will use if it detects a change, or set it to None if you want to print the result to the stdout. Available options are `[slack, mail]`.

- `channel`

  - `Slack`: Slack webhook
  
  - `Mail`: Mail config

## Example usage

The usage examples below show just the simplest tasks you can accomplish using `gosna -h`. 

```
    ______      _____   ______  __   __    _____    
   /_/\___\    ) ___ ( / ____/\/_/\ /\_\  /\___/\   
   ) ) ___/   / /\_/\ \) ) __\/) ) \ ( ( / / _ \ \  
  /_/ /  ___ / /_/ (_\ \\ \ \ /_/   \ \_\\ \(_)/ /  
  \ \ \_/\__\\ \ )_/ / /_\ \ \\ \ \   / // / _ \ \  
   )_)  \/ _/ \ \/_\/ /)____) ))_) \ (_(( (_( )_) ) 
   \_\____/    )_____( \____\/ \_\/ \/_/ \/_/ \_\/  
                  v1.0 @alanEG

Usage: ./gosna [options]
   -run          Run type [add,check]
   -timeout      Requests timeout (default 5)
   -thread       Requests thread
   -header  -H   Requests header
   -repeat  -r   Repeat Check  [m,h,d]
   -dynamic      Check dynamic (default false)
   -config       Config file   (default ~/.gosna_config.json)
   -no-color     Disable color
```

### Normal Usage

There are two run type 
- add
- check 

##### Example: `add` 

Add the url to check it later 

`cat url.txt | gosna -run add`

Gosnal will retrieve the HTML content, save it in the work directory, and then store the URL along with its hash and options in the config file under the target JSON object.

##### Example: `check`

Check the url if there change 

`gosna -run check`

Gosna will retrieve the URLs from the config file, then fetch the new content, save it in the work directory, and compare the new content with the existing content.

Reapet check for diffrent time by flag `-repeat` or `-r`
In this options you have 3 subs options [m,h,d] minute,hour,day
You can use one of this options in flag `-repeat` or `-r` by 

`gosna -run check -repeat 1m`
### Dynamic Usage

Here's how the dynamic check operates when adding a new URL:

1 - Send five requests and save the responses in the tmp directory.

2 - Extract the dynamic lines from the response, then remove the dynamic lines from a selected file.

3 - Save response of the selected one to the work directory  

Here's how the dynamic check functions when the tool checks for any changes:

1 - Fetch the latest response.

2 - Remove the dynamic lines from the response.

3 - Compare the new response with the old response to check for changes.


##### Example
`cat urls.txt | gosna -run add -dynamic`

##### Note
Sometimes, the dynamic option may not work effectively.

For instance, consider that we need to monitor `https://www.youtube.com/`. Since it prints different line numbers with each request, the tool will display this URL in each check.

## License

Gosna is released under MIT license. See [LICENSE](https://github.com/alanEG/Gosna/blob/main/LICENSE).
