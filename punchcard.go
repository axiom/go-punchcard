package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"antiklimax.se/go-colproximate"
	"github.com/lucasb-eyer/go-colorful"
)

// This table contains the "keypoints" of the colorgradient you want to generate.
// The position of each keypoint has to live in the range [0,1]
type GradientTable []struct {
	Col colorful.Color
	Pos float64
}

// This is the meat of the gradient computation. It returns a HCL-blend between
// the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (self GradientTable) GetInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(self)-1; i++ {
		c1 := self[i]
		c2 := self[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}
	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return self[len(self)-1].Col
}

// This is a very nice thing Golang forces you to do!
// It is necessary so that we can write out the literal of the colortable below.
func MustParseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("MustParseHex: " + err.Error())
	}
	return c
}

var gradients = map[string]GradientTable{
	"fire": {
		{MustParseHex("#ffffcc"), 0.0},
		{MustParseHex("#ffeda0"), 0.111111},
		{MustParseHex("#fed976"), 0.333333},
		{MustParseHex("#feb24c"), 0.444444},
		{MustParseHex("#fd8d3c"), 0.555556},
		{MustParseHex("#fc4e2a"), 0.666667},
		{MustParseHex("#e31a1c"), 0.777778},
		{MustParseHex("#bd0026"), 0.888889},
		{MustParseHex("#800026"), 1.0},
	},

	"blackred": {
		{MustParseHex("#b2182b"), 0.0},
		{MustParseHex("#d6604d"), 0.111111},
		{MustParseHex("#f4a582"), 0.333333},
		{MustParseHex("#fddbc7"), 0.444444},
		{MustParseHex("#ffffff"), 0.555556},
		{MustParseHex("#e0e0e0"), 0.666667},
		{MustParseHex("#bababa"), 0.777778},
		{MustParseHex("#878787"), 0.888889},
		{MustParseHex("#4d4d4d"), 1.0},
	},

	"rainbow": {
		{MustParseHex("#9e0142"), 0.0},
		{MustParseHex("#d53e4f"), 0.1},
		{MustParseHex("#f46d43"), 0.2},
		{MustParseHex("#fdae61"), 0.3},
		{MustParseHex("#fee090"), 0.4},
		{MustParseHex("#ffffbf"), 0.5},
		{MustParseHex("#e6f598"), 0.6},
		{MustParseHex("#abdda4"), 0.7},
		{MustParseHex("#66c2a5"), 0.8},
		{MustParseHex("#3288bd"), 0.9},
		{MustParseHex("#5e4fa2"), 1.0},
	},

	"pink": {
		{MustParseHex("#f7f4f9"), 0.0},
		{MustParseHex("#e7e1ef"), 0.111111},
		{MustParseHex("#d4b9da"), 0.333333},
		{MustParseHex("#c994c7"), 0.444444},
		{MustParseHex("#df65b0"), 0.555556},
		{MustParseHex("#e7298a"), 0.666667},
		{MustParseHex("#ce1256"), 0.777778},
		{MustParseHex("#980043"), 0.888889},
		{MustParseHex("#67001f"), 1.0},
	},

	"rainbow2": {
		{MustParseHex("#5e4fa2"), 0},
		{MustParseHex("#3288bd"), 0.181818},
		{MustParseHex("#66c2a5"), 0.272727},
		{MustParseHex("#abdda4"), 0.363636},
		{MustParseHex("#e6f598"), 0.454545},
		{MustParseHex("#ffffbf"), 0.545455},
		{MustParseHex("#fee08b"), 0.636364},
		{MustParseHex("#fdae61"), 0.727273},
		{MustParseHex("#f46d43"), 0.818182},
		{MustParseHex("#d53e4f"), 0.909091},
		{MustParseHex("#9e0142"), 1.0},
	},

	"blackwhite": {
		{MustParseHex("#000000"), 0},
		{MustParseHex("#ffffff"), 1.0},
	},

	"orangepurple": {
		{MustParseHex("#7f3b08"), 0},
		{MustParseHex("#b35806"), 0.181818},
		{MustParseHex("#e08214"), 0.272727},
		{MustParseHex("#fdb863"), 0.363636},
		{MustParseHex("#fee0b6"), 0.454545},
		{MustParseHex("#f7f7f7"), 0.545455},
		{MustParseHex("#d8daeb"), 0.636364},
		{MustParseHex("#b2abd2"), 0.727273},
		{MustParseHex("#8073ac"), 0.818182},
		{MustParseHex("#542788"), 0.909091},
		{MustParseHex("#2d004b"), 1.0},
	},

	"greenpink": {
		{MustParseHex("#276419"), 0},
		{MustParseHex("#4d9221"), 0.181818},
		{MustParseHex("#7fbc41"), 0.272727},
		{MustParseHex("#b8e186"), 0.363636},
		{MustParseHex("#e6f5d0"), 0.454545},
		{MustParseHex("#f7f7f7"), 0.545455},
		{MustParseHex("#fde0ef"), 0.636364},
		{MustParseHex("#f1b6da"), 0.727273},
		{MustParseHex("#de77ae"), 0.818182},
		{MustParseHex("#c51b7d"), 0.909091},
		{MustParseHex("#8e0152"), 1.0},
	},
}

type when struct {
	day  time.Weekday
	hour int
}

type buckets map[when]int

func (b buckets) Sum() float64 {
	var sum float64
	for _, v := range b {
		sum += float64(v)
	}
	return sum
}

func (b buckets) Avg() float64 {
	return b.Sum() / float64(len(b))
}

func (b buckets) Max() float64 {
	var max float64

	for _, v := range b {
		if float64(v) > max {
			max = float64(v)
		}
	}
	return max
}

func (b buckets) Print() {
	pal := colproximate.XTerm256[16:len(colproximate.XTerm256)]

	keypoints := gradients["blackwhite"]

	if pal, ok := gradients[*palette]; ok {
		keypoints = pal
	} else if *palette != "" {
		fmt.Println("Valid palettes:")
		for name, _ := range gradients {
			fmt.Printf("- %v\n", name)
		}
		os.Exit(1)
	}

	for _, day := range []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday} {
		fmt.Printf("%9v ", day)

		for hour := 0; hour < 24; hour++ {
			v := b[when{day: day, hour: hour}]
			f := float64(v) / b.Max()
			c := keypoints.GetInterpolatedColorFor(f)
			i := 16 + pal.Index(c)

			if v == 0 && *transparent {
				fmt.Print("  ")
			} else {
				fmt.Printf("\x1b[48;5;%vm  \x1b[0m", i)
			}
		}
		fmt.Println()
	}

	if *scale {
		fmt.Println()
		fmt.Printf("%9v", "")
		for i := 0; i < 48; i++ {
			f := float64(i) / float64(48)
			c := keypoints.GetInterpolatedColorFor(f)
			i := 16 + pal.Index(c)
			fmt.Printf("\x1b[48;5;%vm \x1b[0m", i)
		}
	}
}

var (
	palette     = flag.String("palette", "", "")
	scale       = flag.Bool("scale", false, "Show color scale")
	transparent = flag.Bool("transparent", false, "Missing data is not colored")
)

func main() {
	flag.Parse()

	r := bufio.NewReader(os.Stdin)
	buckets := make(buckets)

	for {
		line, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Fatal(err)
			os.Exit(1)
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)

		t, err := time.Parse("2006-01-02 15:04:05 -0700", line)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}

		buckets[when{day: t.Weekday(), hour: t.Hour()}] += 1
	}

	buckets.Print()
}
