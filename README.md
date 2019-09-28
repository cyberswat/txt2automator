# txt2autmoator

Simple command line utility that accepts a file as input and converts it into an applescript file. You can then use the applescript file in Apple's Automator program to simulate typing the original files contents as a user in iTerm.

# Example Usage
```
txt2automator convert ~/my-first-script.sh ~/my-second-script.sh

~/my-first-script.sh successfully converted to my-first-script-sh.scpt
~/my-second-script.sh successfully converted to my-second-script-sh.scpt
```
# Using Automator
Once you have the file converted you'll want to copy the converted file's contents in Apple's Automator program.
* Open Automator
* Select New Document
* Select the Workflow icon and press the Choose button
* Add **Launch Application** as the first step
* Select iTerm as the application to launch
* Add **Run AppleScript** as the second step
* Replace the auto generated Applescript with the contents of the converted file
* Press the run button to see your original file magically typed into the iTerm window

# Potential Uses
I use this by writing a script file that contents what I would normally type while giving a presentation. Converting the script into an automator workflow allows me to use Quicktime to record a video of what I would normally type. Embedding the video into a program like Keynote allows me to capture the output into a presentation slide for consistent playback without human error.