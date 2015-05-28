# Contributing Guide Lines

## Code Submission

* All code must be run through `go fmt` before commiting
* There isn't a strict line length limit but try to not to have crazy long lines (120ish chars)
* A single commit is prefered to multiple unless they are clearly seperate 
  (touching two different packages for example)
* Commit messages should be in the form `packagename: commit message` where the commit message
  should be in lower case (exluding acronyms like HTTP and JSON)
* For the commit messages the root package `github.com/thinkofdeath/steven` is referred to as `steven`
* In the case where multiple packages are touched and it doesn't make sense to split the commits 
  the package names can be seperated in a list via a comma e.g. `render,steven: commit message`
* The code must work on Linux, Windows and Mac unless its in a platform specific area file 
  (with a build tag or file extension)

## Issue Submission

* steven-log.txt should be submitted with all bug reports.
* Logs/Crash logs should be wrapped in \`\`\` to keep it readable. 
* In the case the log is too large use https://gist.github.com
* The title of the issue should clearly state the issue and the package (if known) that the issue occurs in
  e.g. `steven: missing model for red flowers`
* Please include details about your operating system and graphics card in the issue to help with
  tracking down issues