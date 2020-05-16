# bbb-download
Downloader for bigbluebutton presentations

I need a simple executable file wroks on Windows, MAc-Os or Linux systems. Program only needs ffmepg executable on your path, capable of encoding with parameter libx264

Program asks for url of presentation, then downloads slides, creates slides.mp4, downloads webcams.webm, creates webcams.mp4, then merges slides.mp4 and webcams.mp4 side by side
Unfortunately cursor movements, desktop share, chat windows is not included. Best way to include them is, clone and build the code from https://github.com/jibon57/bbb-recorder on a Linux system

This is a proof of concept software, programmed in Oldskool fashion, straight 200 lines of code, by an c64/Amiga/pocketpc (6502, 68x00, arm) assembler coder.
This is also my first program coded with go language
If you want to see the downloaded pictures, and created videos before deleted, comment two lines near the end of the code

    os.RemoveAll(presentationId+"/")  // delete temporary dir
    err = os.Remove("video_list.txt") // delete video-list file
  
 I didn't tried the code on Linux or Mac-OS systems, but I expect no problem
