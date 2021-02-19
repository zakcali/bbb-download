package main

import ( 
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"math"
	"log"
	"io"
	"bufio"
	"net/http"
	"io/ioutil"
	"strings"
)

var webcamsFile = "webcams.webm"
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
	metaUrl := baseUrl + "/metadata.xml"
//	fmt.Println ("baseUrl= ", baseUrl)

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

  	shapesBody, err := ioutil.ReadAll (responseShapes.Body) 
    if err != nil {
        log.Fatal(err)
    }
    shapes:=string (shapesBody)

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

     // Find and print slide timings, image Urls
	durations := make(map[int]float64)
	vidnames := make(map[int]string)
	imgnames := make(map[int]string)
	inValue,  outValue, videoLength, truncated:= 0.0, 0.0, 0.0, 0.0
	inSrc, outSrc, pngSrc := "0.0", "10.5", "presentation/"  
	i:=1 // number of png pictures for slide

	//parse for in= out= href= from /shapes.svg
	ins := strings.Split(shapes, "in=\"")

	fmt.Println ("Downloading Slides")
// download all off the slides loop
	for k :=1; k < len (ins); k++ {
	intext := strings.SplitAfter (ins[k], "\"")
	realin:=strings.Split(intext[0],"\"")
	inSrc = (realin[0])

	outtext := strings.SplitAfter (ins[k], "out=\"")
	realout:=strings.Split(outtext[1],"\"")
	outSrc = (realout[0])

	imgtext := strings.SplitAfter (ins[k], "xlink:href=\"")
	realpng:=strings.Split(imgtext[1],"\"")
	pngSrc = (realpng[0])

	inValue, _ = strconv.ParseFloat (inSrc,64)
	outValue, _ = strconv.ParseFloat (outSrc,64)
	truncated = ( outValue*10-inValue*10)/10
	durations[i] = truncated
	imgnames [i] = "s" + strconv.Itoa(i)+".png"
	vidnames[i] = "v" + strconv.Itoa(i) + ".mp4"

	imgUrl := baseUrl + "/" + pngSrc
//	fmt.Println (inSrc, " ", outSrc, " ", durations[i], " ", pngSrc, imgnames [i], " ", vidnames[i])
	fmt.Print ("Downloading: ", imgnames [i], "\r") // print to same line just like a counter
		if err := DownloadFile(presentationId+"/"+imgnames[i], imgUrl); err != nil {
				panic(err)
			}	
	i++
	}
nSlides=i-1  // if only one slide is converted, than nSlides = 1
fmt.Println("Number of slides =",nSlides)
// end of download all off the slides loop
		//correct duration of last slide
outValue, _ = strconv.ParseFloat(duration[0],64)
outValue= outValue/1000
videoLength=math.Round ( outValue*100)/100
fmt.Println("Length of presentation =",  videoLength) 
truncated = ( videoLength*10-inValue*10)/10
durations[i-1] = math.Round (truncated*100)/100
fmt.Println("Duration of last slide according to meta.xml =", durations[i-1]) 
	
		// create mp4 files from png files
fmt.Println ("Creating videos from slide pictures, duration is given as seconds")
	for j:=1; j<=nSlides; j++ 	{
	fmt.Print (imgnames[j], " ", vidnames[j], " ", durations [j], "\r") // print to same line just like a counter
	cmd := exec.Command("ffmpeg","-loop", "1", "-r", "5", "-f", "image2", 
						"-i", presentationId +"/"+imgnames[j],
						"-c:v", "libx264", "-r", "24", "-t", fmt.Sprint(durations[j]), "-pix_fmt", "yuv420p", 
						"-vf", "scale='if(gt(a,1024/768),1024,-2)':'if(gt(a,1024/768),-2,768)',pad=1024:768:(ow-iw)/2:(oh-ih)/2:color=white",        // as close as 800x600
						presentationId+"/"+vidnames[j]  )
	cmd.Run()

	}

if nSlides == 1 { 
			slidesFile="v1.mp4" // none of the videos are merged
} else { // there are more than one video file, so merge them	
	//create video_list.txt file to cancat with ffmpeg
	    f, err := os.Create("video_list.txt")
    if err != nil {
        fmt.Println(err)
        return
    }
	for j:=1; j<=nSlides; j++ 	{
	    _, err := f.WriteString("file " + presentationId+"/"+vidnames[j] + "\n")
    if err != nil {
        fmt.Println(err)
        f.Close()
        return }
	}
	err = f.Close()
    if err != nil {
        fmt.Println(err)
        return }
	//concat slide videos to create one piece of video file: slides.mp4

fmt.Println ("merging slide videos to create: slides.mp4")
	cmd := exec.Command("ffmpeg","-f", "concat", "-safe", "0", "-i", "video_list.txt", 
						"-c", "copy", presentationId+"/"+slidesFile)
	cmd.Run()
fmt.Println ("slide videos merged")

} 

//convert webcams video file to  webcamsRight.mp4
fmt.Println ("converting ",webcamsFile, " to  webcamsRight.mp4")
	cmd := exec.Command("ffmpeg","-i", presentationId+"/"+webcamsFile,
						 "-q:a", "0", "-q:v", "0", 
						 "-vf", "scale=512:-2,pad=height=768:color=white", 
						presentationId+"/"+"webcamsRight.mp4")
	cmd.Run()
fmt.Println ("merging slides and webcams side by side")
	cmd = exec.Command("ffmpeg", "-i", presentationId+"/"+slidesFile,
				"-i", presentationId+"/"+"webcamsRight.mp4",
				"-filter_complex", "[0:v][1:v]hstack=inputs=2[v]", 
				"-t", fmt.Sprint(videoLength),
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
