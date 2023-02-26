package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Equal(a, b []string) bool {
	for _, ea := range a {
		found := false
		for _, eb := range b {
			if ea == eb {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func InstallTW() error {
	cmd := exec.Command("npm", "i", "tailwindcss", "postcss", "autoprefixer")
	fmt.Println("[#] Installing Tailwindcss...")
	err := cmd.Run()
	if err != nil {
		panic("[!] Tailwind installation failed, Exiting.")
	} else {
		fmt.Println("[#] Tailwind installed.")
	}
	return err
}

func generateConfig() error {
	cmd := exec.Command("npx", "tailwindcss", "init", "-p")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("[!] Configuration file could not be generated. Exiting.")
	}
	return err
}

func cC(srcs, exts []string) []string {
	var content []string
	if len(exts) > 1 {
		for _, src := range srcs {
			content = append(content, fmt.Sprintf("`./%s/**/*.{%s}`", src, strings.Join(exts, ",")))
		}
	} else {
		for _, src := range srcs {
			content = append(
				content, fmt.Sprintf("`./%s/**/*.%s`", src, strings.Join(exts, ",")))
		}
	}

	return content
}

func genCSS() {
	d := []string{"@tailwind base;", "@tailwind components;", "@tailwind utilities;"}
	f, err := os.Create("RENAME_ME.css")
	if err != nil {
		panic("[!] Could not generate file. Exiting")
	}
	defer f.Close()
	_, error := f.Write([]byte(strings.Join(d, "\n")))
	if error != nil {
		panic("[!] Could not generate file. Exiting")

	}
}

func main() {
	// Valid Exts array
	validExts := []string{"vue", "svelte", "jsx", "tsx", "js", "ts"}
	var preExts string

	fmt.Printf("[?] File extension/s of frontend framework  (No dots and seperated by commas) : ")
	fmt.Scan(&preExts)
	// Check if Exts are valid
	Exts := strings.Split(preExts, ",")
	isEqual := Equal(Exts, validExts)
	if !isEqual {
		fmt.Println("[!] One or more of the file extensions you've input are invalid ") // Err
		return
	} else {
		err := InstallTW()
		if err != nil {
			return
		}
	}

	fmt.Println("[#] Generating configuration file..")
	err := generateConfig()
	if err != nil {
		panic("[!] Something went wrong. Exiting.")

	}
	fmt.Printf("[?] src/app directories (Seperated by commas) : ")
	var preSourceFiles string
	fmt.Scan(&preSourceFiles)
	SourceFiles := strings.Split(preSourceFiles, ",")
	preContent := cC(SourceFiles, Exts)

	content := fmt.Sprintf(`/*@type {import('tailwindcss').Config}*/ 
module.exports = {
  content: [%s],
  theme: {
    extend: {},
  },
  plugins: [],
}`, strings.Join(preContent, ","))

	var _x error = os.WriteFile("./tailwind.config.cjs", []byte(content), 0644)
	if _x != nil {
		panic("[!] Couldnt overwrite Configuration file. Exiting.")
	}

	fmt.Println("[#] Generating a CSS file with tailwind directives....")
	genCSS()

	fmt.Println("[#] Done !")

}
