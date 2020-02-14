package version

import (
	"log"
	"os"
	"time"
)

var (
	Commit    string    = ""
	Branch    string    = ""
	Tag       string    = ""
	Timestamp time.Time = time.Now()
)

func Log(log *log.Logger) {
	log.Printf("  ** %s\n", os.Args[0])
	log.Printf("  * Build Tiem:  %s\n", Timestamp.Format("2006-01-02 15:04:05"))
	log.Printf("  * Commit:      %s\n", Commit)
	log.Printf("  * Tag:         %s\n", Tag)
	log.Printf("  * Branch:      %s\n", Branch)
	log.Println("  **")
}
