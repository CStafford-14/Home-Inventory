package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/aquilax/truncate"
)

type Location struct {
	Name       string
	Color      string
	Identifier string
	Containers []string
}

type Container struct {
	Name     string
	Color    string
	Type     string
	Location string
	ID       string
	Contents []string
}

var devMode bool
var library = "/var/lib/HomeInventory/"

func main() {
	// Capture the desired port. Defaults to 8090.
	port := flag.String("port", "8090", "Port to be used for the server")
	devModeRaw := flag.Bool("dev", false, "Development mode (Routes events to working directory)")
	flag.Parse()

	if *devModeRaw {
		devMode = true
		library = ""
	}

	// Start the Webserver
	staticFileServer := http.FileServer(http.Dir("./websource/static/"))

	http.HandleFunc("/locations/", locationHandler)
	http.HandleFunc("/containers/", containerHandler)
	http.Handle("/common/", staticFileServer)
	http.HandleFunc("/", homeHandler)

	if devMode {
		fmt.Print(strings.Join([]string{"Starting server at port ", *port, " in development mode\n"}, ""))
	} else {
		fmt.Print(strings.Join([]string{"Starting server at port ", *port, "\n"}, ""))
	}
	if err := http.ListenAndServe(strings.Join([]string{":", *port}, ""), nil); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" && r.Method != "GET" {
		http.Error(w, "Method is not supported.\n\nInstead, please enjoy this ASCII Rick Astley.\n\ntttttfttttttttffffftfffftfftttttttttt111tttt1111111ttttttt111tttt111111tttttttttttttttttttt1111111tt\nttfttttttttttttttttffLLftfffffffftttt1t1111111t11tfffffftttt1111111111111111ttftttftttttttt111111ttt\nttttttttttttttttffffLftttfffffffLLfftt1tttttt111tfffffffftttttttt11111tttttffffftffLfttttttt11111ttt\nttttttttttttttttfffffttffffffLLLfffttt1ttttttt1tffffffffffttttttt11111tfffffffLffttfffffttt1111ttttt\nttttttttttttffffttfffttffffffffftttftttttt11ttt1ttfffffftttttt11tt111111ttffffLLLfttfLLft111111ttttt\nttttttttttffLLLLfttttffLLLfftttttfffttttt1tfffft1ttfffffttfffft11t111tt111ttfffLLftttftt1ttt11tttttt\nttttttttfffLLLLLLffttfLLLfttfffftfLftttttttfffffftttffttffffffftt1111tftt111ttffLLfttt11tffftttttttt\ntttttttfffLLLLLLLLLfttfttttfLLffftffttttt111i;iitfftttfffffffffftt111tft1tttttttfLLfttttfffffftttttt\ntttttttffLLLLLLLLLfttttffffffffftfftttt1;:,::,,,:;ittttffffffffftt1111t11ttttttttttttttffffffffffttt\ntttffftfLLLLLLLLLLfttfftfffttfffLLLftt1i:,,,,,,,,,,:1tttttfffffftt1111tft1ttt1tffttt11tffffLLLLLfttt\nttffffffffLLLLLLLfftfLLftttfLLffLLLftt1i;;;;;;;i;:,:tftt11tfffft11111tffft1tt11ttttt1ttfffLLLLLLLfff\ntfffffttttfLLLffttttfffffttLLLLfLLLft1;iiiii11111i::tffft11tttt11t111tfffttfftt1ttfftttffLLLLLLLLffL\ntttfftttttffLftttttttttttttfLLLffLLftt11;;;iiiiiii:;tffffft11tt1111111tfftfffftttttttttttfLLLLffffff\nttfLffffffftttfffffffffttffttLLfffLft11i;;;;i;;iii;tffffttt1tffft11111tfttfffttt1tttttttttfLLftffttt\nttffffffffftttffffffffftffffttfftfttt11i;;;i1iii1iitffft11tttttfft1111tt1tffttffttfffffLLfffftfLffLL\nttffffffffftttffffffffttttttttttffftt11i;;;iii111ii1tt111tfft1ttft111111t1ttttfftttffffLLLfttfLLLLLL\ntttffffffffttfffffffftt1ttt111ttffftttti;;;;iiiiii11111ttfffft11t11111tfftttttttt1tfffffLLftttfLLLLf\nttttttttfffttfffttttt1ttfft1111ttffttt1i;;;iiii;;itt11tffffffft1111111ttttt1ttfft11ttffffLfftfLLLLLf\ntftttttttttttt1tttt11tfffft1tt11tfft111;;;;;iii;;it111ttfffffft11t1111tttt1tt1tffftt1111ttfftfLffftt\ntfttffffftt11ttftfftttffftttfft11tti;;1;;;:;iii::i111111ttfft111tt1111ttt11tf1tffffftttttttttttttfft\n111tfffffftttffftttt11tfftttft1i;:,..,1i;;;;iii;i;:;;1t111t11111tt1111tt111tft1fffffttfffftt11tffLLf\n1tttffffft11ttffftt1111tt1ii:,,......,i1i;;;ii;11;...,:;i111tt11111111111t11tf1tfffttfffffffttfLffLf\ntfttttftt11t11ttt11ttt11;,,..........,11t1iiii1t1:,,.....,:i111111t1111111111t1tffttttffffttt1tfffLf\ntffttttttttfttt111ttttti,.............:;;iii11t1i,,.........,ittt11111111tttt11tttfLffttffttttttffft\nttttttttttttttt111tttt1:................:;;;;;i;:,...........:ttt1111111tttffttttffLLffttttffftttttf\ntfffttttftttttt11ttttti.................,;;;ii;;:............,1ttt11111ttttffftttffffffttttttttffttf\ntffft1ttfffftf111ttttti..................:;;;i;;:,...........,1tt11111111ttffttttffLLfLfttfffffffttf\ntffft1ttfttttt111ttttt;..................,;;;;;::,...........,1t1111111111tttttttffLLfLfttfffffffttf\ntffft1tttttttt111ttttt;...........,,.....,;;:;:::............,111111111tt111tttttfffLLLfttfffftffttf\nttttt1tttttttt111tttti,........,;ii:......:;:::::............,ittttt111ttt11tft1tfffffLftttfftfffttt\nttttt1tttt111t111tttt;........:;;;;:,.....,;:::::.............ittttt1111t1111tt1ttfttttttttttttffttt\nttttt1tttttttt11tttt1,.......,:::;;;,......::::::.............;tt1t1111111ttt1tttfttttttt1tttttffttt\n1tttt1tttttttt11tttti.........,:::;:.......,:::::..........,,:;1t11111111ttff1111ttttttt111ttttffttt\n111111tttttt111111t1,..........,:::........,:::::..........,:;;i1111t111ttttt11111ttttt1111ttttttt1t\n1111111111111111111;............,,.........,::::,.........,:;;;i1111111111111111111111111111111tt111\n11111111111111111111:.......................,:::,.. ......,:;;ii111111111111111111111111111111111111\n111111111111111111111:.......,..............,:::,..........,::i11t1111111111111111111111111111111111\n1111111111111111111111:,,..,, ..............,,,::,............,1111111111111111111111111111111111111\n1111111111111111111111111i11:...............,..,:,,.....    ..:1t11111111111111111111111111111111111\n111111111111111111111111111;................,::::,,.....::::;i11tttt111111t1111111111111111111111111\n111111111111111111111111111,.................ii;;::,....:1tttt1ttttttttttttttttt11tt111111ttttt11111\n1111111111111111111111111ti..................;i;;;:,.....;t1ttttttttttttttttttttt1ttttttttttttt11111", http.StatusNotFound)
		return
	}

	rawIndex, err := os.ReadFile(library + "websource/static/home.html")

	if err != nil {
		log.Fatal(err)
	}

	var locfiles []Location
	var contfiles []Container
	var locs string
	var boxs string

	loclib, err := os.ReadDir(library + "locations/")
	for _, f := range loclib {
		var add *Location
		work, err := os.ReadFile(library + "locations/" + f.Name())
		if err != nil {
			fmt.Println("Panic caused by read JSON file. \n\n\n")
			log.Fatal(err)
		}
		err = json.Unmarshal(work, &add)
		if err != nil {
			fmt.Println("Panic caused by parse JSON file. \n\n\n")
			log.Fatal(err)
		}
		locfiles = append(locfiles, *add)
	}

	contlib, err := os.ReadDir(library + "containers/")
	for _, f := range contlib {
		var add *Container
		work, err := os.ReadFile(library + "containers/" + f.Name())
		if err != nil {
			fmt.Println("Panic caused by read JSON file. \n\n\n")
			log.Fatal(err)
		}
		err = json.Unmarshal(work, &add)
		if err != nil {
			fmt.Println("Panic caused by parse JSON file. \n\n\n")
			log.Fatal(err)
		}
		contfiles = append(contfiles, *add)
	}

	var buff1, buff2 bytes.Buffer
	body := string(rawIndex)
	temp, err := template.New("tmp").Parse(body)
	loc, err := template.New("loc").Parse("<section><h2 style=\"font-size: 40px; margin-bottom: 20px;\">Locations</h2><div style=\"display: flex;justify-content: start;\">{{range .}}<a href=\"/locations/{{.Identifier}}\"><div class=\"locationBox\" style=\"background-color:#{{.Color}};\"><span class=\"label\">{{.Name}}</span></div></a>{{end}}</div></section>")
	box, err := template.New("box").Parse("<section><h2 style=\"font-size: 40px; margin-bottom: 20px;\">Containers</h2><div style=\"display: flex;justify-content: start;\">{{range .}}<a href=\"/containers/{{.ID}}\"><div class=\"containerBox\" style=\"background-color:#{{.Color}};\"><span class=\"label\">{{.Name}}</span><span class=\"subLabel\">{{.Type}} #{{.ID}}</span><span class=\"subLabel\">{{.Location}}</span></div></a>{{end}}</div></section>")

	loc.Execute(io.Writer(&buff1), locfiles)
	box.Execute(io.Writer(&buff2), contfiles)
	if len(locfiles) > 0 {
		locs = buff1.String()
	} else {
		locs = "<section><h2 style=\"font-size: 40px; margin-bottom: 20px;\">Locations</h2><div style=\"display: flex;justify-content: start;\"><p>No locations.</p></div></section>"
	}
	if len(contfiles) > 0 {
		boxs = buff2.String()
	} else {
		boxs = "<section><h2 style=\"font-size: 40px; margin-bottom: 20px;\">Containers</h2><div style=\"display: flex;justify-content: start;\"><p>No containers.</p></div></section>"
	}

	if err != nil {
		fmt.Println("Panic caused by template parse.\n\n\n")
		log.Fatal(err)
	}
	err = temp.ExecuteTemplate(w, "tmp", locs+"\n"+boxs)
	if err != nil {
		fmt.Println("Panic caused by template fill.\n\n\n")
		log.Fatal(err)
	}
}

func locationHandler(w http.ResponseWriter, r *http.Request) {
	var contData []Container

	if r.Method != "POST" && r.Method != "GET" {
		http.Error(w, "Method is not supported.\n\nInstead, please enjoy this ASCII Rick Astley.\n\ntttttfttttttttffffftfffftfftttttttttt111tttt1111111ttttttt111tttt111111tttttttttttttttttttt1111111tt\nttfttttttttttttttttffLLftfffffffftttt1t1111111t11tfffffftttt1111111111111111ttftttftttttttt111111ttt\nttttttttttttttttffffLftttfffffffLLfftt1tttttt111tfffffffftttttttt11111tttttffffftffLfttttttt11111ttt\nttttttttttttttttfffffttffffffLLLfffttt1ttttttt1tffffffffffttttttt11111tfffffffLffttfffffttt1111ttttt\nttttttttttttffffttfffttffffffffftttftttttt11ttt1ttfffffftttttt11tt111111ttffffLLLfttfLLft111111ttttt\nttttttttttffLLLLfttttffLLLfftttttfffttttt1tfffft1ttfffffttfffft11t111tt111ttfffLLftttftt1ttt11tttttt\nttttttttfffLLLLLLffttfLLLfttfffftfLftttttttfffffftttffttffffffftt1111tftt111ttffLLfttt11tffftttttttt\ntttttttfffLLLLLLLLLfttfttttfLLffftffttttt111i;iitfftttfffffffffftt111tft1tttttttfLLfttttfffffftttttt\ntttttttffLLLLLLLLLfttttffffffffftfftttt1;:,::,,,:;ittttffffffffftt1111t11ttttttttttttttffffffffffttt\ntttffftfLLLLLLLLLLfttfftfffttfffLLLftt1i:,,,,,,,,,,:1tttttfffffftt1111tft1ttt1tffttt11tffffLLLLLfttt\nttffffffffLLLLLLLfftfLLftttfLLffLLLftt1i;;;;;;;i;:,:tftt11tfffft11111tffft1tt11ttttt1ttfffLLLLLLLfff\ntfffffttttfLLLffttttfffffttLLLLfLLLft1;iiiii11111i::tffft11tttt11t111tfffttfftt1ttfftttffLLLLLLLLffL\ntttfftttttffLftttttttttttttfLLLffLLftt11;;;iiiiiii:;tffffft11tt1111111tfftfffftttttttttttfLLLLffffff\nttfLffffffftttfffffffffttffttLLfffLft11i;;;;i;;iii;tffffttt1tffft11111tfttfffttt1tttttttttfLLftffttt\nttffffffffftttffffffffftffffttfftfttt11i;;;i1iii1iitffft11tttttfft1111tt1tffttffttfffffLLfffftfLffLL\nttffffffffftttffffffffttttttttttffftt11i;;;iii111ii1tt111tfft1ttft111111t1ttttfftttffffLLLfttfLLLLLL\ntttffffffffttfffffffftt1ttt111ttffftttti;;;;iiiiii11111ttfffft11t11111tfftttttttt1tfffffLLftttfLLLLf\nttttttttfffttfffttttt1ttfft1111ttffttt1i;;;iiii;;itt11tffffffft1111111ttttt1ttfft11ttffffLfftfLLLLLf\ntftttttttttttt1tttt11tfffft1tt11tfft111;;;;;iii;;it111ttfffffft11t1111tttt1tt1tffftt1111ttfftfLffftt\ntfttffffftt11ttftfftttffftttfft11tti;;1;;;:;iii::i111111ttfft111tt1111ttt11tf1tffffftttttttttttttfft\n111tfffffftttffftttt11tfftttft1i;:,..,1i;;;;iii;i;:;;1t111t11111tt1111tt111tft1fffffttfffftt11tffLLf\n1tttffffft11ttffftt1111tt1ii:,,......,i1i;;;ii;11;...,:;i111tt11111111111t11tf1tfffttfffffffttfLffLf\ntfttttftt11t11ttt11ttt11;,,..........,11t1iiii1t1:,,.....,:i111111t1111111111t1tffttttffffttt1tfffLf\ntffttttttttfttt111ttttti,.............:;;iii11t1i,,.........,ittt11111111tttt11tttfLffttffttttttffft\nttttttttttttttt111tttt1:................:;;;;;i;:,...........:ttt1111111tttffttttffLLffttttffftttttf\ntfffttttftttttt11ttttti.................,;;;ii;;:............,1ttt11111ttttffftttffffffttttttttffttf\ntffft1ttfffftf111ttttti..................:;;;i;;:,...........,1tt11111111ttffttttffLLfLfttfffffffttf\ntffft1ttfttttt111ttttt;..................,;;;;;::,...........,1t1111111111tttttttffLLfLfttfffffffttf\ntffft1tttttttt111ttttt;...........,,.....,;;:;:::............,111111111tt111tttttfffLLLfttfffftffttf\nttttt1tttttttt111tttti,........,;ii:......:;:::::............,ittttt111ttt11tft1tfffffLftttfftfffttt\nttttt1tttt111t111tttt;........:;;;;:,.....,;:::::.............ittttt1111t1111tt1ttfttttttttttttffttt\nttttt1tttttttt11tttt1,.......,:::;;;,......::::::.............;tt1t1111111ttt1tttfttttttt1tttttffttt\n1tttt1tttttttt11tttti.........,:::;:.......,:::::..........,,:;1t11111111ttff1111ttttttt111ttttffttt\n111111tttttt111111t1,..........,:::........,:::::..........,:;;i1111t111ttttt11111ttttt1111ttttttt1t\n1111111111111111111;............,,.........,::::,.........,:;;;i1111111111111111111111111111111tt111\n11111111111111111111:.......................,:::,.. ......,:;;ii111111111111111111111111111111111111\n111111111111111111111:.......,..............,:::,..........,::i11t1111111111111111111111111111111111\n1111111111111111111111:,,..,, ..............,,,::,............,1111111111111111111111111111111111111\n1111111111111111111111111i11:...............,..,:,,.....    ..:1t11111111111111111111111111111111111\n111111111111111111111111111;................,::::,,.....::::;i11tttt111111t1111111111111111111111111\n111111111111111111111111111,.................ii;;::,....:1tttt1ttttttttttttttttt11tt111111ttttt11111\n1111111111111111111111111ti..................;i;;;:,.....;t1ttttttttttttttttttttt1ttttttttttttt11111", http.StatusNotFound)
		return
	}

	var arc *Location

	id := truncate.Truncate(r.URL.String(), len(r.URL.String())-11, "", truncate.PositionStart)
	work, err := os.ReadFile(library + "locations/" + id + ".json")
	if err != nil {
		fmt.Println("Panic caused by read JSON file. \n\n\n")
		log.Fatal(err)
	}
	err = json.Unmarshal(work, &arc)
	if err != nil {
		fmt.Println("Panic caused by parse JSON file. \n\n\n")
		log.Fatal(err)
	}

	if arc != nil {
		rawIndex, err := os.ReadFile(library + "websource/static/locations.html")
		body := string(rawIndex)
		var buff1, buff2 bytes.Buffer

		for _, r := range arc.Containers {
			var add *Container
			work, err := os.ReadFile(library + "containers/" + r + ".json")
			if err != nil {
				fmt.Println("Panic caused by read JSON file. \n\n\n")
				log.Fatal(err)
			}
			err = json.Unmarshal(work, &add)
			if err != nil {
				fmt.Println("Panic caused by parse JSON file. \n\n\n")
				log.Fatal(err)
			}
			contData = append(contData, *add)
		}

		temp, err := template.New("tmp2").Parse(body)
		head, err := template.New("head").Parse("<div id=\"banner\" style=\"background-color:#{{.Color}}\"><h1 style=\"font-size: 20pt;\">{{.Name}}</h1></div>")
		info, err := template.New("info").Parse("<div id=\"wrapper\"><ul>{{range .}}<a href=\"/containers/{{.ID}}\"><li>{{.Name}} - <em>{{.ID}}</em></li></a>{{end}}</ul></div>")
		if err != nil {
			fmt.Println("Panic caused by template parse.\n\n\n")
			log.Fatal(err)
		}

		head.Execute(io.Writer(&buff1), arc)
		info.Execute(io.Writer(&buff2), contData) // Make this an array of a struct that has the ID and name of the identified container.

		head2 := buff1.String()
		info2 := buff2.String()

		err = temp.ExecuteTemplate(w, "tmp2", head2+"\n"+info2)
		if err != nil {
			fmt.Println("Panic caused by template fill.\n\n\n")
			log.Fatal(err)
		}

		if err != nil {
			log.Fatal(err)
		}
	} else {
		rawIndex, err := os.ReadFile("websource/static/404.html")
		body := string(rawIndex)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, body)
	}
}

func containerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "GET" {
		http.Error(w, "Method is not supported.\n\nInstead, please enjoy this ASCII Rick Astley.\n\ntttttfttttttttffffftfffftfftttttttttt111tttt1111111ttttttt111tttt111111tttttttttttttttttttt1111111tt\nttfttttttttttttttttffLLftfffffffftttt1t1111111t11tfffffftttt1111111111111111ttftttftttttttt111111ttt\nttttttttttttttttffffLftttfffffffLLfftt1tttttt111tfffffffftttttttt11111tttttffffftffLfttttttt11111ttt\nttttttttttttttttfffffttffffffLLLfffttt1ttttttt1tffffffffffttttttt11111tfffffffLffttfffffttt1111ttttt\nttttttttttttffffttfffttffffffffftttftttttt11ttt1ttfffffftttttt11tt111111ttffffLLLfttfLLft111111ttttt\nttttttttttffLLLLfttttffLLLfftttttfffttttt1tfffft1ttfffffttfffft11t111tt111ttfffLLftttftt1ttt11tttttt\nttttttttfffLLLLLLffttfLLLfttfffftfLftttttttfffffftttffttffffffftt1111tftt111ttffLLfttt11tffftttttttt\ntttttttfffLLLLLLLLLfttfttttfLLffftffttttt111i;iitfftttfffffffffftt111tft1tttttttfLLfttttfffffftttttt\ntttttttffLLLLLLLLLfttttffffffffftfftttt1;:,::,,,:;ittttffffffffftt1111t11ttttttttttttttffffffffffttt\ntttffftfLLLLLLLLLLfttfftfffttfffLLLftt1i:,,,,,,,,,,:1tttttfffffftt1111tft1ttt1tffttt11tffffLLLLLfttt\nttffffffffLLLLLLLfftfLLftttfLLffLLLftt1i;;;;;;;i;:,:tftt11tfffft11111tffft1tt11ttttt1ttfffLLLLLLLfff\ntfffffttttfLLLffttttfffffttLLLLfLLLft1;iiiii11111i::tffft11tttt11t111tfffttfftt1ttfftttffLLLLLLLLffL\ntttfftttttffLftttttttttttttfLLLffLLftt11;;;iiiiiii:;tffffft11tt1111111tfftfffftttttttttttfLLLLffffff\nttfLffffffftttfffffffffttffttLLfffLft11i;;;;i;;iii;tffffttt1tffft11111tfttfffttt1tttttttttfLLftffttt\nttffffffffftttffffffffftffffttfftfttt11i;;;i1iii1iitffft11tttttfft1111tt1tffttffttfffffLLfffftfLffLL\nttffffffffftttffffffffttttttttttffftt11i;;;iii111ii1tt111tfft1ttft111111t1ttttfftttffffLLLfttfLLLLLL\ntttffffffffttfffffffftt1ttt111ttffftttti;;;;iiiiii11111ttfffft11t11111tfftttttttt1tfffffLLftttfLLLLf\nttttttttfffttfffttttt1ttfft1111ttffttt1i;;;iiii;;itt11tffffffft1111111ttttt1ttfft11ttffffLfftfLLLLLf\ntftttttttttttt1tttt11tfffft1tt11tfft111;;;;;iii;;it111ttfffffft11t1111tttt1tt1tffftt1111ttfftfLffftt\ntfttffffftt11ttftfftttffftttfft11tti;;1;;;:;iii::i111111ttfft111tt1111ttt11tf1tffffftttttttttttttfft\n111tfffffftttffftttt11tfftttft1i;:,..,1i;;;;iii;i;:;;1t111t11111tt1111tt111tft1fffffttfffftt11tffLLf\n1tttffffft11ttffftt1111tt1ii:,,......,i1i;;;ii;11;...,:;i111tt11111111111t11tf1tfffttfffffffttfLffLf\ntfttttftt11t11ttt11ttt11;,,..........,11t1iiii1t1:,,.....,:i111111t1111111111t1tffttttffffttt1tfffLf\ntffttttttttfttt111ttttti,.............:;;iii11t1i,,.........,ittt11111111tttt11tttfLffttffttttttffft\nttttttttttttttt111tttt1:................:;;;;;i;:,...........:ttt1111111tttffttttffLLffttttffftttttf\ntfffttttftttttt11ttttti.................,;;;ii;;:............,1ttt11111ttttffftttffffffttttttttffttf\ntffft1ttfffftf111ttttti..................:;;;i;;:,...........,1tt11111111ttffttttffLLfLfttfffffffttf\ntffft1ttfttttt111ttttt;..................,;;;;;::,...........,1t1111111111tttttttffLLfLfttfffffffttf\ntffft1tttttttt111ttttt;...........,,.....,;;:;:::............,111111111tt111tttttfffLLLfttfffftffttf\nttttt1tttttttt111tttti,........,;ii:......:;:::::............,ittttt111ttt11tft1tfffffLftttfftfffttt\nttttt1tttt111t111tttt;........:;;;;:,.....,;:::::.............ittttt1111t1111tt1ttfttttttttttttffttt\nttttt1tttttttt11tttt1,.......,:::;;;,......::::::.............;tt1t1111111ttt1tttfttttttt1tttttffttt\n1tttt1tttttttt11tttti.........,:::;:.......,:::::..........,,:;1t11111111ttff1111ttttttt111ttttffttt\n111111tttttt111111t1,..........,:::........,:::::..........,:;;i1111t111ttttt11111ttttt1111ttttttt1t\n1111111111111111111;............,,.........,::::,.........,:;;;i1111111111111111111111111111111tt111\n11111111111111111111:.......................,:::,.. ......,:;;ii111111111111111111111111111111111111\n111111111111111111111:.......,..............,:::,..........,::i11t1111111111111111111111111111111111\n1111111111111111111111:,,..,, ..............,,,::,............,1111111111111111111111111111111111111\n1111111111111111111111111i11:...............,..,:,,.....    ..:1t11111111111111111111111111111111111\n111111111111111111111111111;................,::::,,.....::::;i11tttt111111t1111111111111111111111111\n111111111111111111111111111,.................ii;;::,....:1tttt1ttttttttttttttttt11tt111111ttttt11111\n1111111111111111111111111ti..................;i;;;:,.....;t1ttttttttttttttttttttt1ttttttttttttt11111", http.StatusNotFound)
		return
	}

	var arc *Container

	id := truncate.Truncate(r.URL.String(), len(r.URL.String())-11, "", truncate.PositionStart)
	work, err := os.ReadFile(library + "containers/" + id + ".json")
	if err != nil {
		fmt.Println("Panic caused by read JSON file. \n\n\n")
		log.Fatal(err)
	}
	err = json.Unmarshal(work, &arc)
	if err != nil {
		fmt.Println("Panic caused by parse JSON file. \n\n\n")
		log.Fatal(err)
	}

	if arc != nil {
		rawIndex, err := os.ReadFile(library + "websource/static/locations.html")
		body := string(rawIndex)
		var buff1, buff2 bytes.Buffer

		temp, err := template.New("tmp2").Parse(body)
		head, err := template.New("head").Parse("<div id=\"banner\" style=\"background-color:#{{.Color}}\"><h1 style=\"font-size: 20pt;\">{{.Name}}</h1><span>{{.Type}} #{{.ID}} | {{.Location}}</span></div>")
		info, err := template.New("info").Parse("<div id=\"wrapper\"><ul>{{range .}}<li>{{.}}</li>{{end}}</ul></div>")
		if err != nil {
			fmt.Println("Panic caused by template parse.\n\n\n")
			log.Fatal(err)
		}

		head.Execute(io.Writer(&buff1), arc)
		info.Execute(io.Writer(&buff2), arc.Contents) // Make this an array of a struct that has the ID and name of the identified container.

		head2 := buff1.String()
		info2 := buff2.String()

		err = temp.ExecuteTemplate(w, "tmp2", head2+"\n"+info2)
		if err != nil {
			fmt.Println("Panic caused by template fill.\n\n\n")
			log.Fatal(err)
		}

		if err != nil {
			log.Fatal(err)
		}
	} else {
		rawIndex, err := os.ReadFile("websource/static/404.html")
		body := string(rawIndex)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, body)
	}
}
