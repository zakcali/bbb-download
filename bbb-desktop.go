package main

import ( 
	"fmt"
	"os"
	"os/exec"
	"log"
	"io"
	"bufio"
	"net/http"
	"io/ioutil"
	"strings"
)

var webcamsFile = "webcams.webm"
var deskshareFile = "deskshare.webm"
var slidesFile = "slides.mp4"
var nSlides int =1

func main () {

	fmt.Println ("Bigbluebutton video creator/downloader")
	fmt.Print ("Enter url of conference/lecture: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
    	presentationUrl := scanner.Text()
	result := strings.SplitAfter(presentationUrl,"?meetingId=")
	presentationId := result [1]

	result2 := strings.Split(result [0], "/playback/" )
	baseUrl := result2 [0] + "/presentation/" + presentationId
	shapesUrl := baseUrl + "/shapes.svg"
	webcamsUrl := baseUrl + "/video/"
	deskshareUrl := baseUrl + "/deskshare/"
	metaUrl := baseUrl + "/metadata.xml"
	fmt.Println ("baseUrl= ", baseUrl)

	//read duration of recording and meeting name from meta.xml
    responseMeta, err := http.Get(metaUrl)
    if err != nil {
        log.Fatal(err)
    }
    defer responseMeta.Body.Close()
	metaBody, err := ioutil.ReadAll (responseMeta.Body) 
    if err != nil {
        log.Fatal(err)
    }
    //finding correct duration-ending of last slide
    timeString := strings.SplitAfter(string(metaBody),"<duration>")
    duration := strings.Split(timeString[1],"</duration>")
    fmt.Println ("duration of recording=", duration[0], "ms")

    meetingString := strings.SplitAfter(string(metaBody),"<meetingName>")
    meetingName := strings.Split(meetingString[1],"</meetingName>")
    fmt.Println ("name of the meeting=", meetingName[0])
	//read content of the shapes.svg file, and assign shapes to it
    responseShapes, err := http.Get(shapesUrl)
    if err != nil {
        log.Fatal(err)
    }
    defer responseShapes.Body.Close()
    if err != nil {
        log.Fatal(err)
    }
   	fmt.Println ("creating directory: ", presentationId)
	if _, err := os.Stat(presentationId); os.IsNotExist(err) {
    os.Mkdir(presentationId, 0700) // create temporary dir
}
// download webcams
	fmt.Print ("downloading webcams",  "\r")
	if err := DownloadFile(presentationId+"/"+webcamsFile, webcamsUrl+"/"+webcamsFile); err != nil {
		panic(err) 	}
	fmt.Println (webcamsFile, " file is downloaded",  "\r")
	fi, err := os.Stat(presentationId+"/"+webcamsFile)
		if err != nil {
    		panic(err)
			}
	fileSize := fi.Size()
//	fmt.Print ("file size =", fileSize,  "\r")
	if fileSize<1000 { // returned 404 error text
// webcams.webm file is so small that real webcams file must be in mp4 format
		webcamsFile = "webcams.mp4"
		fmt.Print ("downloading webcams", "\r")
		if err := DownloadFile(presentationId+"/"+webcamsFile, webcamsUrl+"/"+webcamsFile); err != nil {
			panic(err) 	}
		fmt.Println (webcamsFile, "file is downloaded",  "\r")
	}

// download deskshare
	fmt.Print ("downloading deskshare",  "\r")
	if err := DownloadFile(presentationId+"/"+deskshareFile, deskshareUrl+"/"+deskshareFile); err != nil {
		panic(err) 	}
	fmt.Println (deskshareFile, " file is downloaded",  "\r")
	fi, err = os.Stat(presentationId+"/"+deskshareFile)
		if err != nil {
    		panic(err)
			}
	fileSize = fi.Size()
//	fmt.Print ("file size =", fileSize,  "\r")
	if fileSize<1000 { // returned 404 error text
// webcams.webm file is so small that real webcams file must be in mp4 format
		deskshareFile = "deskshare.mp4"
		fmt.Print ("downloading deskshare", "\r")
		if err := DownloadFile(presentationId+"/"+deskshareFile, webcamsUrl+"/"+deskshareFile); err != nil {
			panic(err) 	}
		fmt.Println (deskshareFile, "file is downloaded",  "\r")
	}

//convert webcams video file to  webcamsRight.mp4
fmt.Println ("converting ",webcamsFile, " to  webcamsRight.mp4")
	cmd := exec.Command("ffmpeg","-i", presentationId+"/"+webcamsFile,
			 "-q:a", "0", "-q:v", "0", 
			 "-vf", "scale=512:-2,pad=height=768:color=white", 
			presentationId+"/"+"webcamsRight.mp4")
	cmd.Run()
//convert deskshare video file to  deskshare.mp4
fmt.Println ("converting ",deskshareFile, " to  deskshare.mp4")
	cmd = exec.Command("ffmpeg","-i", presentationId+"/"+deskshareFile,
			 "-q:a", "0", "-q:v", "0", 
			 "-vf", "scale=1024:-2,pad=height=768:color=white", 
			presentationId+"/"+"deskshare.mp4")
	cmd.Run()

fmt.Println ("merging slides and webcams side by side")
	cmd = exec.Command("ffmpeg", "-i", presentationId+"/"+"deskshare.mp4",
			"-i", presentationId+"/"+"webcamsRight.mp4",
			"-filter_complex", "[0:v][1:v]hstack=inputs=2[v]", 
			"-map", "[v]", "-map", "1:a", meetingName[0]+".mp4")
	cmd.Run()
fmt.Println ("Name of the final video is: ", meetingName[0])
os.RemoveAll(presentationId+"/")  // delete temporary dir
err = os.Remove("video_list.txt") // delete video-list file
}

// DownloadFile will download a url to a local file. 
func DownloadFile(filepath string, url string) error {

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()
    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    return err
}