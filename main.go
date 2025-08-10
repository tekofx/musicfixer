package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bogem/id3v2"
)

func main() {

	tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
	if err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}
	defer tag.Close()

	// // Read tags
	// fmt.Println(tag.Artist())
	// fmt.Println(tag.Title())

	// // Set tags
	// tag.SetArtist("Aphex Twin")
	// tag.SetTitle("Xtal")
	//
	//
	frames := tag.AllFrames()

	for k, v := range frames {

		fmt.Println(k, v)
	}

	pictures := tag.GetFrames("APIC")
	for _, f := range pictures {
		pic, ok := f.(id3v2.PictureFrame)
		if !ok {
			log.Fatal("Couldn't assert picture frame")
		}

		// Use a fallback filename if description is empty
		filename := pic.Description
		if filename == "" {
			filename = "cover"
		}

		// Ensure the file extension matches the actual MIME type
		ext := ".jpg"
		if pic.MimeType == "image/png" {
			pictures := tag.GetFrames(tag.CommonID("Attached picture"))
			for _, f := range pictures {
				pic, ok := f.(id3v2.PictureFrame)
				if !ok {
					log.Fatal("Couldn't assert picture frame")
				}

				// Do something with picture frame
				fmt.Println(pic.Description)
				fmt.Println(pic.MimeType)
				outputFile, err := os.Create(pic.Description + ".jpg")
				if err != nil {
					fmt.Println(err)
				}
				defer outputFile.Close()
				_, err = pic.WriteTo(outputFile)
				if err != nil {
					fmt.Println(err)
				}
			}
		} else if pic.MimeType == "image/gif" {
			ext = ".gif"
		}

		outputFile, err := os.Create(filename + ext)
		if err != nil {
			log.Println("Error creating file:", err)
			continue
		}

		// Write raw picture data (not the frame structure)
		_, err = outputFile.Write(pic.Picture)
		if err != nil {
			log.Println("Error writing image data:", err)
		}

		err = outputFile.Close()
		if err != nil {
			log.Println("Error closing file:", err)
		}

		fmt.Printf("Saved image: %s%s (MIME: %s)\n", filename, ext, pic.MimeType)
	}

	// Write tag to file.mp3
	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}

	// outputDir, dry, removeOriginalFolder := flags.SetupFlags()
	// dir := flags.GetDir()

	// albumSongs, errors := model.ReadAlbums(dir)
	// if errors != nil {
	// 	fmt.Println("Error reading songs metadata:")
	// 	for _, error := range errors {
	// 		fmt.Printf("%v\n", error)
	// 	}
	// 	os.Exit(0)
	// }

	// if albumSongs == nil {
	// 	fmt.Printf("No songs found in %s", dir)
	// 	os.Exit(0)
	// }

	// model.SetNewFilePaths(albumSongs)

	// if dry {
	// 	flags.DryRun(albumSongs, outputDir)
	// }

	// err := model.RenameSongs(*albumSongs, outputDir)
	// if err != nil {
	// 	fmt.Printf("Error renaming songs: %v\n", err)
	// 	os.Exit(0)
	// }

	// fmt.Printf("Done! All song renamed in %s", dir)

	// if removeOriginalFolder {
	// 	err := os.RemoveAll(dir)
	// 	if err != nil {
	// 		fmt.Printf("Error removing original directories: %v", err)
	// 	}
	// }

}
