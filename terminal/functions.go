package terminal


import(
	"fmt"
	"github.com/chzyer/readline"
	"strings"
	"io"
	"go-deliver/model"
	"go-deliver/database"
)

var context string = "main"
var prompt string = "go-deliver (\033[0;32m%s\033[0;0m)\033[31m >> \033[0;0m"



var MainCompleter = readline.NewPrefixCompleter(
	readline.PcItem("payload",
		readline.PcItem("add",
			// Vlerat ketu do merren nga databaza ne te ardhmen.
			readline.PcItem("mshta"),
			readline.PcItem("regsrv32"),
			readline.PcItem("powershell"),
			readline.PcItem("javascript"),
		),
		readline.PcItem("delete"),
		readline.PcItem("list"),

	),
	readline.PcItem("host"),
	//	readline.PcItem("listeners") To be implemented later .
)

var PayloadCompleter = readline.NewPrefixCompleter(
	readline.PcItem("set",
		readline.PcItem("name"),
		readline.PcItem("content_type"),
		readline.PcItem("host_blacklist"),
		readline.PcItem("host_whitelist"),
		readline.PcItem("data_file"),
		readline.PcItem("data_b64"),
		readline.PcItem("ptype"),
		//readline.PcItem("listener"), // This is will be implemented later.
		),
	readline.PcItem("unset",
		readline.PcItem("content_type"),
		readline.PcItem("host_blacklist"),
		readline.PcItem("host_whitelist"),
		readline.PcItem("data_file"),
		readline.PcItem("data_b64"),
		readline.PcItem("type"),
		//readline.PcItem("listener"), // This is will be implemented later.
	),
	readline.PcItem("options"),
	readline.PcItem("create"),
	readline.PcItem("back"),

)


var HostCompleter = readline.NewPrefixCompleter(
	readline.PcItem("set",
		readline.PcItem("name"),
		readline.PcItem("type",
			readline.PcItem("ip"),
			readline.PcItem("subnet"),
		),
		readline.PcItem("data"),
		),
)



func handlePayloadCreation(ptype string, l *readline.Instance)  {
	payload := model.Payload{}
	payload.Ptype = ptype

	//fmt.Println(fmt.Sprintf("Will create a payload named with the ptype %s",payload.ptype))
	l.Config.AutoComplete = PayloadCompleter
	l.SetPrompt(fmt.Sprintf(prompt,"payload-options"))

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)

		temp := strings.Split(line," ")
		command := temp[0]
		switch command{
		case "back":
			backMain(l)
			return
		case "options":
			// To be fixed
			fmt.Println(payload)
		case "set":
			if len(temp) == 3{
				key := temp[1]
				value := temp[2]
				switch key{
				case "name":
					payload.Name = value
				case "content_type":
					payload.Content_type = value
				case "host_whitelist":
					payload.Host_whitelist = value
				case "host_blacklist":
					payload.Host_blacklist = value
				case "data_file":
					payload.Data_file = value
				case "data_b64":
					payload.Data_b64 = value
				case "ptype":
					payload.Ptype = value
				case "type_id":
					payload.Type_id, _ = fmt.Sscanf("%d",value)
				}

			}
		case "unset":
			fmt.Println("Unset the payload.")
		case "create":

			database.InsertPayload(payload)
		default:


		}
	}
}



func backMain(l *readline.Instance){
	context = "main"
	l.SetPrompt(fmt.Sprintf(prompt,"main"))
	l.Config.AutoComplete = MainCompleter
}

func handleInput(line string ,l *readline.Instance)  {


	line = strings.TrimSpace(line)
	temp := strings.Split(line," ")

	if len(temp) > 2 {

		command := temp[1]
		switch {

		// Handle the payload functions
		case strings.HasPrefix(line, "payload "):

			var ptype string = temp[2]
			switch  command{
			case "add":
				handlePayloadCreation(ptype,l)
			case "delete":
				fmt.Println("Remove a payload")
			default:
				fmt.Println("Invalid command")
			}

		// Handle the Hosts functions
		case strings.HasPrefix(line, "host "):
			switch command {
			case "add":
				fmt.Println("Add a host")
			case "delete":
				fmt.Println("Remove a host")
			default:
				fmt.Println("Invalid command")
			}

		}
	}

	
}

func StartTerminal()  {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          fmt.Sprintf(prompt,"main"),
		HistoryFile:     "history.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		AutoComplete:	 MainCompleter,

	})


	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {

		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		handleInput(line,l)


	}
}



