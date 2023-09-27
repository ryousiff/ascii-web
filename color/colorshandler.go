package color

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	hexRegex string = `^[A-Fa-f0-9]{1,6}$`
	hslRegex string = `^[0-9]{1,3}`
)

type RGB struct {
	R, G, B int64
}

func (rgb *RGB) String() string {
	if rgb == nil {
		return "\033[0m"
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", rgb.R, rgb.G, rgb.B)
}

func Factory(w http.ResponseWriter, r *http.Request) *RGB {
	input := r.FormValue("inputpicker")

	// switch input {
	// case "red":
	// 	return &RGB{255, 0, 0}
	// case "green":
	// 	return &RGB{0, 255, 0}
	// case "yellow":
	// 	return &RGB{255, 255, 0}
	// case "blue":
	// 	return &RGB{0, 0, 255}
	// case "orange":
	// 	return &RGB{255, 165, 0}
	// case "purple":
	// 	return &RGB{128, 0, 128}
	// case "pink":
	// 	return &RGB{255, 192, 203}
	// case "cyan":
	// 	return &RGB{0, 255, 255}
	// case "brown":
	// 	return &RGB{165, 42, 42}
	// case "gray":
	// 	return &RGB{128, 128, 128}
	// }
	// if strings.HasPrefix(input, "hsl") {
	// 	return parseHsl(input)
	// }
	// if strings.HasPrefix(input, "#") {
	// 	return parseHex(input)
	// }
	if strings.HasPrefix(input, "rgb") {
		return parseRgb(input)
	}
	os.Exit(0)
	return nil
}

func parseRgb(rgb string) *RGB {
	rgb = strings.TrimLeft(rgb, "rgb(")
	rgb = strings.TrimRight(rgb, ")")

	rgbSlice := strings.Split(rgb, ",")
	// invalid
	if len(rgbSlice) != 3 {
		return nil
	}
	// clean
	rgbRe := regexp.MustCompile(hslRegex)

	for i := range rgbSlice {
		rgbSlice[i] = strings.TrimSpace(rgbSlice[i])
		if !rgbRe.MatchString(rgbSlice[i]) {
			return nil
		}
	}

	r_dec, _ := strconv.Atoi(rgbSlice[0])
	g_dec, _ := strconv.Atoi(rgbSlice[1])
	b_dec, _ := strconv.Atoi(rgbSlice[2])

	return &RGB{int64(r_dec), int64(g_dec), int64(b_dec)}

}

func parseHex(hex string) *RGB {
	hex = strings.TrimPrefix(hex, "#")
	hexRe := regexp.MustCompile(hexRegex)
	isHex := hexRe.MatchString(hex)
	if !isHex {
		return nil
	}
	// if the hex is not in full form, add padding to expand
	if len(hex) < 6 {
		padding := strings.Repeat("0", 6-len(hex))
		hex += padding
	}
	r_hex := hex[:2]
	g_hex := hex[2:4]
	b_hex := hex[4:]

	r_dec, _ := strconv.ParseInt(r_hex, 16, 64)
	g_dec, _ := strconv.ParseInt(g_hex, 16, 64)
	b_dec, _ := strconv.ParseInt(b_hex, 16, 64)

	return &RGB{r_dec, g_dec, b_dec}
}

// Anything below this is ChatGPT
func parseHsl(hsl string) *RGB {
	hsl = strings.TrimLeft(hsl, "hsl(")
	hsl = strings.TrimRight(hsl, ")")

	hslSlice := strings.Split(hsl, ",")
	// invalid
	if len(hslSlice) != 3 {
		return nil
	}
	// clean
	hslRe := regexp.MustCompile(hslRegex)
	for i := range hslSlice {
		hslSlice[i] = strings.TrimSuffix(hslSlice[i], "%")
		hslSlice[i] = strings.TrimSpace(hslSlice[i])
		if !hslRe.MatchString(hslSlice[i]) {
			return nil
		}
	}

	h, _ := strconv.ParseFloat(hslSlice[0], 64)
	s, _ := strconv.ParseFloat(hslSlice[1], 64)
	l, _ := strconv.ParseFloat(hslSlice[2], 64)

	return hslToRGB(h, s, l)
}

func hslToRGB(h, s, l float64) *RGB {
	h = h / 360.0
	s = s / 100.0
	l = l / 100.0

	var r, g, b int
	if s == 0 {
		r = int(l * 255.0)
		g = int(l * 255.0)
		b = int(l * 255.0)
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - (l * s)
		}

		p := 2*l - q
		r = int(round(hueToRGB(p, q, h+1.0/3.0) * 255.0))
		g = int(round(hueToRGB(p, q, h) * 255.0))
		b = int(round(hueToRGB(p, q, h-1.0/3.0) * 255.0))
	}

	return &RGB{R: int64(r), G: int64(g), B: int64(b)}
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}

	if t < 1.0/6.0 {
		return p + (q-p)*6.0*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6.0
	}

	return p
}

// to avoid rounding errors
func round(x float64) float64 {
	if math.Abs(x-math.Floor(x+0.5)) < 1e-9 {
		return math.Floor(x + 0.5)
	}
	return math.Round(x)
}
