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

- [Download](https://github.com/alaneg/gosna/releases/latest) a prebuilt binary from [releases page](https://github.com/alaneg/gosna/releases/latest), unpack and run!
  
  _or_
- If you have recent go compiler installed: `go install github.com/alanEG/Gosna` (the same command works for updating)
  
  _or_
- `git clone https://github.com/alaneg/gosna ; cd gosna ; go build`

Gosna depends on Go go1.18.1 or greater.


## Configuration file

Gosna doesn't work if you don't setup config file for save the urls and file names and other config   

There are three ways the Gosna tool take config file from them 

1 - `gosna_config` env

2 - `~/.gosna_config.json`

3 - `-config` option

The defualt config file will be 

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

`Is_first` Don't edit this will edit by the tool 

`Directory_work` Edit this for the directory the tool will add the content of urls in it 

`Directory_result` Edit this for the directory the tool will save html result in it 

`Channel_use` Edit this for the channel the tool will use if it found change or use `None` if you need print the result in stdout channel available `[slack,mail]`

- `channel`

  - `Slack` slack webhook
  
  - `Mail` Mail config

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
Add for add the url to check it later 

`cat url.txt | gosna -run add`

Gosna will gets the url content then format it 

Then save the content in directory work file 

And add it to config file in Target json object  


##### Example: `check`
Chcek for check the url if there change 

`gosna -run check`

Gosna will gets the urls from config file then get the new content 
And save it in the work directory then diff the new content with new content 

Reapet check for diffrent time by flag `-repeat` or `-r`
In this options you have 3 subs options [m,h,d] minute,hour,day
You can use one of this options in flag `-repeat` or `-r` by 

`gosna -run check -repeat 1m`
### Dynamic Usage

Here simply how check dynamic works when add

1 - Send 5 request and save them responses in `tmp` directory

2 - Get the the dynamic lines from the responses in the `tmp` directory   

3 - Save the bigest response length 

4 - Remove the dynamic lines from the file saved 

5 - Add `true` in object is `Target[n].Status` for knowing the tool it is dynamic url 

Here simply how check dynamic works when Check

1 - Check the url is dynamic or not by `Target[n].Status`

2 - Send request and remove the dynamic lines then save it in the file in work directory 

3 - Diff the old content with new contnet 


##### Example
`cat urls.txt | gosna -run add -dynamic`

##### Note
Some time Dynamic option not working well 

Let's imagen we need to monitor the `https://www.youtube.com/`

Because the it printing diffrent line number in each request 

The tool will printed this url in each check  

## License

Gosna is released under MIT license. See [LICENSE](https://github.com/alanEG/Gosna/blob/main/LICENSE).
