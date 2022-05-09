package main

type Json struct {
	Config Config
	Target []Target
}

type Config struct {
	Is_first         bool
	Directory_work   string
	Directory_result string
	Channel_use      string
	Channel          Channel
}

type Target struct {
	Url      string
	Filename string
	Dynamic  Dynmic
	Headers  map[string]string
}

type Channel struct {
	Slack string
	Mail  Email
}

type Email struct {
	Tls      bool
	From     string
	To       string
	Host     string
	Port     int
	Email    string
	Password string
}

type Dynmic struct {
	Status      bool
	DynamicLine []int
}

type arrayFlags []string

var (
	configFileFlag string
	configFile     string
	flagRepeat     string
	Run            string
	add            string
	flagNoColor    bool
	flagDynmaic    bool
	flagTimeout    int
	Thread         int
	data           Json
	Headers        arrayFlags
	Header         map[string]string
)
