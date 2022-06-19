module github.com/aschey/bubbleprompt

go 1.18

require (
	github.com/alecthomas/participle/v2 v2.0.0-alpha8
	github.com/aschey/tui-tester v0.0.0-20220619042203-388d67e69f05
	github.com/charmbracelet/bubbles v0.10.4-0.20220301123521-e349920524a2
	github.com/charmbracelet/bubbletea v0.20.0
	github.com/charmbracelet/lipgloss v0.5.0
	github.com/onsi/ginkgo/v2 v2.1.4
	github.com/onsi/gomega v1.19.0
)

require github.com/aschey/termtest v0.7.2-0.20220618051905-02b060238256 // indirect

require (
	github.com/ActiveState/vt10x v1.3.1 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Netflix/go-expect v0.0.0-20220104043353-73e0943537d2 // indirect
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/creack/pty v1.1.18 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/muesli/ansi v0.0.0-20211031195517-c9f0611b6c70 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.11.1-0.20220212125758-44cd13922739 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c // indirect
	golang.org/x/term v0.0.0-20220411215600-e5f449aeb171 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.11 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/aschey/tui-tester => ../tui-tester
