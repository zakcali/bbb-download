# bbb-download
Downloader for bigbluebutton presentations

I need a single executable file works on Windows, Mac-Os or Linux systems. This software also needs ffmepg executable on your path, capable of encoding with parameter libx264

Program asks for url of presentation, then downloads slides, creates slides.mp4, downloads webcams.webm, creates webcams.mp4, then merges slides.mp4 and webcams.mp4 side by side.
Unfortunately cursor movements, desktop share, chat windows is not included. Best way to include them is, clone and build the code from https://github.com/jibon57/bbb-recorder on a Linux system.

What I share is a proof of concept software, programmed in Oldskool fashion, straight 200 lines of code. This is also my first program coded with go language. If you want to see the downloaded pictures, and created videos before deleted, comment two lines near the end of the code

    os.RemoveAll(presentationId+"/")  // delete temporary dir
    err = os.Remove("video_list.txt") // delete video-list file
  
 I didn't try the code on Linux or Mac-OS systems, but I expect no problem
 
Do-not use it. It consumes so much processor sources, and can't show cursor movements, desktop share, chat windows. Use https://github.com/jibon57/bbb-recorder on Linux instead. bbb-recorder uses codes from puppetcam depends on xvfb which doesn't exist for Windows operating system. For windows 10 compatible puppetcam look at: https://github.com/Ventricule/html2screen 

skips svg files created during a poll, because ffmpeg cannot convert svg files to mp4 format.

there are two builds, one for webcams.webm file and other is for webcams.mp4 file resides in the server. I have no time to check if there is really a file resides in the server. golang Download function retrieves a html file named webcams.webm containing 404 code even webcams.webm file is not on the server.

# important
webroot secure anywhere antivirus cause crashing executable files build with go language. I build executables with an older version of go, v 1.13.15 for intel386. I use webroot at home, and I don't want to uninstall it. for more information look at: https://github.com/golang/go/issues/40878 and https://community.webroot.com/webroot-secureanywhere-antivirus-12/apps-built-in-go-language-are-crashing-on-windows-345321
