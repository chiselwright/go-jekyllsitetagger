package jekyllsitetagger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// WalkMatch will recursively walk through a directory and return the paths to all files whose name matches the given pattern
// https://stackoverflow.com/questions/55300117/how-do-i-find-all-files-that-have-a-certain-extension-in-go-regardless-of-depth
func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// Unique takes a list of strings and returns the list with duplicated removed
func Unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// FindTags takes a path to a file and scans the file for the "tags:" entry in
// the header block
func FindTags(path string) ([]string, error) {
	var tags []string

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// have we seen the first divider
	var seenFirstDivider bool = false

	scanner := bufio.NewScanner(file)

	tagRE := regexp.MustCompile("^tags:")
	dividerRE := regexp.MustCompile("^---$")

	for scanner.Scan() {
		if dividerRE.MatchString(scanner.Text()) {
			if seenFirstDivider {
				// it must be our second one, meaning No More Content
				break
			}

			// that was our first and only allowed divider
			seenFirstDivider = true

		} else {
			// if we haven't seen a divider, then the file is wonky!
			if !seenFirstDivider {
				file.Close()
				return nil, fmt.Errorf("file should start with '---' marker")
			}
			if tagRE.MatchString(scanner.Text()) {
				fmt.Println(scanner.Text())
				tags = append(tags, ParseTagLine(scanner.Text())...)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		file.Close()
		return nil, err
	}

	return tags, nil
}

// ParseTagLine takes a string of "tags: t1 t2 t2" and returns a slice of the
// parsed tags
func ParseTagLine(tagLine string) []string {
	// if we split the line as whitespace we get our tags
	// nothing complicated about the parsing in our simple little world
	split := strings.Split(tagLine, " ")
	// the first element should be tags: which we don't care about
	return split[1:]
}

// GenerateTagFiles takes a "post" directory to scan, and outputs tag files into
// the "tag" output directory
func GenerateTagFiles(postDir, tagDir string) {
	files, err := WalkMatch("testdata", "*.md")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%+v\n", files)

	var alltags []string

	for _, file := range files {
		tags, err := FindTags(file)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%+v\n", tags)
		alltags = append(alltags, tags...)
	}
	fmt.Printf("%+v\n", alltags)
	fmt.Printf("%+v\n", Unique(alltags))

	for _, tag := range alltags {
		err := WriteTagFile(tagDir, tag)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

// WriteTagFile outputs "tag".md into tagDir/
func WriteTagFile(tagDir, tag string) error {
	var err error
	// make sure we have our output directory
	if _, err = os.Stat(tagDir); os.IsNotExist(err) {
		err = os.Mkdir(tagDir, 0755)
		if err != nil {
			// log.Fatal(err.Error())
			return err
		}
	}

	// create the file
	f, err := os.Create(
		fmt.Sprintf("%s/%s.md", tagDir, tag),
	)
	if err != nil {
		// log.Fatal(err.Error())
		return err
	}

	// It’s idiomatic to defer a Close immediately after opening a file.
	defer f.Close()

	var lines []string = []string{
		"---\n",
		"layout: tagpage\n",
		fmt.Sprintf("title: Tag: %s\n", tag),
		fmt.Sprintf("tag: Tag: %s\n", tag),
		"robots: noindex\n",
		fmt.Sprintf("permalink: /tag/%s/\n", tag),
		"exclude: true\n",
		"---\n",
	}

	for _, line := range lines {
		_, err = f.WriteString(line)
		if err != nil {
			// log.Fatal(err.Error())
			f.Close()
			return err
		}

	}

	err = f.Sync()
	if err != nil {
		// log.Fatal(err.Error())
		f.Close()
		return err
	}
	f.Close()

	return err
}
