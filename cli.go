package jekyllsitetagger

// CLI is the (kong) structure definition for command-line arguments
var CLI struct {
	Source   string `default:"posts" help:"the directory to scan for *.md files" short:"s"`
	OutputTo string `default:"tags"  help:"the directory to output 'tag'.md files into" short:"o"`
}
