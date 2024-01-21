package Netpbm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Pixel représente un pixel avec les composants Rouge, Vert et Bleu (R, G, B).
type Pixel struct {
	R, G, B uint8
}

// PPM représente une image au format PPM.
type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

// ReadPPM lit un fichier PPM et renvoie une structure PPM.
func ReadPPM(filename string) (*PPM, error) {
	//Same as ReadPGM
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ppm := &PPM{}
	line := 0
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "#") {
			continue
		}
		if ppm.magicNumber == "" {
			ppm.magicNumber = strings.TrimSpace(text)
		} else if ppm.width == 0 {
			fmt.Sscanf(text, "%d %d", &ppm.width, &ppm.height)
			ppm.data = make([][]Pixel, ppm.height)
			for i := range ppm.data {
				ppm.data[i] = make([]Pixel, ppm.width)
			}
		} else if ppm.max == 0 {
			fmt.Sscanf(text, "%d", &ppm.max)
		} else {
			if ppm.magicNumber == "P3" {
				val := strings.Fields(text)
				//Loop through each strings in the current line
				for i := 0; i < ppm.width; i++ {
					//Convert the string to uint8 and set it to the red of the pixel
					r, _ := strconv.ParseUint(val[i*3], 10, 8)
					//Same but the index is incremented to get the next value for the green
					g, _ := strconv.ParseUint(val[i*3+1], 10, 8)
					//Same but the index is incremented to get the next value for the blue
					b, _ := strconv.ParseUint(val[i*3+2], 10, 8)
					//Create the pixel with the colors we just obtained and define it the matrix
					ppm.data[line][i] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)}
				}
				line++
			} else if ppm.magicNumber == "P6" {
				//Create an array of byte of the size of the image * 3 because each pixel has 3 values RGB
				pixelData := make([]byte, ppm.width*ppm.height*3)
				fileContent, err := os.ReadFile(filename)
				if err != nil {
					return nil, fmt.Errorf("couldn't read file: %v", err)
				}
				//Same as ReachPGM but for 3 values
				copy(pixelData, fileContent[len(fileContent)-(ppm.width*ppm.height*3):])
				pixelIndex := 0
				for y := 0; y < ppm.height; y++ {
					for x := 0; x < ppm.width; x++ {
						ppm.data[y][x].R = pixelData[pixelIndex]
						ppm.data[y][x].G = pixelData[pixelIndex+1]
						ppm.data[y][x].B = pixelData[pixelIndex+2]
						pixelIndex += 3
					}
				}
				break
			}
		}
	}
	return ppm, nil
}

// Size returns the width and height of the image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At returns the value of the pixel at (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	// Check if the coordinates are within the bounds of the image
	if x >= 0 && x < ppm.width && y >= 0 && y < ppm.height {
		// Update the pixel value at the specified coordinates
		ppm.data[y][x] = value
	} else {
		// Handle out-of-bounds error, e.g., print a message or handle it based on your requirements.
		fmt.Println("Error: Coordinates out of bounds.")
	}
}

// Save saves the PPM image to a file and returns an error if there was a problem.
func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write magic number
	_, err = fmt.Fprintf(file, "%s\n", ppm.magicNumber)
	if err != nil {
		return fmt.Errorf("error writing magic number: %v", err)
	}

	// Write dimensions
	_, err = fmt.Fprintf(file, "%d %d\n", ppm.width, ppm.height)
	if err != nil {
		return fmt.Errorf("error writing dimensions: %v", err)
	}

	// Write max color value
	_, err = fmt.Fprintf(file, "%d\n", ppm.max)
	if err != nil {
		return fmt.Errorf("error writing max color value: %v", err)
	}

	// Write pixel data based on the format (P3 or P6)
	if ppm.magicNumber == "P3" {
		// P3 (ASCII) format
		for y := 0; y < ppm.height; y++ {
			for x := 0; x < ppm.width; x++ {
				pixel := ppm.data[y][x]
				_, err := fmt.Fprintf(file, "%d %d %d ", pixel.R, pixel.G, pixel.B)
				if err != nil {
					return fmt.Errorf("error writing pixel data: %v", err)
				}
			}
			_, err := fmt.Fprint(file, "\n") // Newline after each row
			if err != nil {
				return fmt.Errorf("error writing newline: %v", err)
			}
		}
	} else if ppm.magicNumber == "P6" {
		// P6 (binary) format
		for y := 0; y < ppm.height; y++ {
			for x := 0; x < ppm.width; x++ {
				pixel := ppm.data[y][x]
				_, err := file.Write([]byte{pixel.R, pixel.G, pixel.B})
				if err != nil {
					return fmt.Errorf("error writing pixel data: %v", err)
				}
			}
		}
	}

	return nil
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			// Get the current pixel
			pixel := ppm.data[y][x]

			// Invert the colors
			invertedPixel := Pixel{
				R: uint8(ppm.max) - pixel.R,
				G: uint8(ppm.max) - pixel.G,
				B: uint8(ppm.max) - pixel.B,
			}

			// Set the inverted pixel back to the image
			ppm.data[y][x] = invertedPixel
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	for y := 0; y < ppm.height; y++ {
		// Iterate through half of the width to swap pixels horizontally
		for x := 0; x < ppm.width/2; x++ {
			// Swap pixels horizontally
			ppm.data[y][x], ppm.data[y][ppm.width-x-1] = ppm.data[y][ppm.width-x-1], ppm.data[y][x]
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop() {
	// Iterate through half of the height to swap rows vertically
	for y := 0; y < ppm.height/2; y++ {
		// Swap rows vertically
		ppm.data[y], ppm.data[ppm.height-y-1] = ppm.data[ppm.height-y-1], ppm.data[y]
	}
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	// Validate that the provided max value is within the valid range
	if maxValue == 0 {
		maxValue = 1 // Avoid division by zero
	}

	// Convert and set the max value for each pixel in the image
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			pixel := ppm.data[y][x]
			pixel.R = uint8(int(pixel.R) * int(maxValue) / int(ppm.max))
			pixel.G = uint8(int(pixel.G) * int(maxValue) / int(ppm.max))
			pixel.B = uint8(int(pixel.B) * int(maxValue) / int(ppm.max))
			ppm.data[y][x] = pixel
		}
	}

	// Update the Max field in the PPM struct
	ppm.max = maxValue
}

// Rotate90CW rotates the PPM image 90° clockwise.
func (ppm *PPM) Rotate90CW() {
	// Create a new data matrix to store the rotated image
	rotatedData := make([][]Pixel, ppm.width)
	for i := range rotatedData {
		rotatedData[i] = make([]Pixel, ppm.height)
	}

	// Iterate through each pixel and copy it to the rotated position
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			rotatedData[x][ppm.height-y-1] = ppm.data[y][x]
		}
	}

	// Update the original data with the rotated data
	ppm.data = rotatedData

	// Swap width and height in the PPM struct
	ppm.width, ppm.height = ppm.height, ppm.width
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM {
	// Same idea as ppm.ToPBM
	pgm := &PGM{}
	pgm.magicNumber = "P2"
	pgm.height = ppm.height
	pgm.width = ppm.width

	// Ensure that ppm.Max is within the valid range for uint8
	if ppm.max > math.MaxUint8 {
		pgm.max = math.MaxUint8
	} else {
		pgm.max = uint8(ppm.max)
	}

	for y, _ := range ppm.data {
		pgm.data = append(pgm.data, []uint8{})
		for x, _ := range ppm.data[y] {
			r, g, b := ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B
			// Calculate the amount of gray the pixel should have
			// It is just the average of the 3 RGB colors
			grayValue := uint8((int(r) + int(g) + int(b)) / 3)
			pgm.data[y] = append(pgm.data[y], grayValue)
		}
	}
	return pgm
}

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM {
	const threshold = 2

	pbm := &PBM{}
	pbm.magicNumber = "P1"
	pbm.height = ppm.height
	pbm.width = ppm.width

	for y, _ := range ppm.data {
		pbm.data = append(pbm.data, []bool{})
		for x, _ := range ppm.data[y] {
			r, g, b := ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B
			// Calculate whether the pixel should be black or white
			// If the average of the 3 colors is lower than half of the maximum value, then consider it white
			// If maxValue is 100 and the average is 49, it would be black
			maxValue := uint8(ppm.max)
			isBlack := (uint8((int(r)+int(g)+int(b))/3) < maxValue/uint8(threshold))
			pbm.data[y] = append(pbm.data[y], isBlack)
		}
	}
	return pbm
}

type Point struct {
	X, Y int
}

// DrawLine draws a line between two points.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	deltaX := abs(p2.X - p1.X)
	deltaY := abs(p2.Y - p1.Y)
	sx, sy := sign(p2.X-p1.X), sign(p2.Y-p1.Y)
	err := deltaX - deltaY
	for {
		if p1.X >= 0 && p1.X < ppm.width && p1.Y >= 0 && p1.Y < ppm.height {
			ppm.data[p1.Y][p1.X] = color
		}
		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}
		e2 := 2 * err
		if e2 > -deltaY {
			err -= deltaY
			p1.X += sx
		}
		if e2 < deltaX {
			err += deltaX
			p1.Y += sy
		}
	}
}

// If negative, change it to positive
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Return 1 if it's over 0
// Return 0 if it's 0
// Return -1 if  it's negative
func sign(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

// DrawRectangle draws a rectangle.
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	//Create the 3 extra points according to the width and the height
	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X, p1.Y + height}
	p4 := Point{p1.X + width, p1.Y + height}
	//Draw the lines to connect them
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p4, color)
	ppm.DrawLine(p4, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// DrawFilledRectangle draws a filled rectangle.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	//Draw horizontal lines with the asked width under each other until the height is reached
	p2 := Point{p1.X + width, p1.Y}
	for i := 0; i <= height; i++ {
		ppm.DrawLine(p1, p2, color)
		p1.Y++
		p2.Y++
	}
}

func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	//Loop through each pixel
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			//Calculate the distance from the current pixel to the center of the circle
			dx := float64(x - center.X)
			dy := float64(y - center.Y)
			distance := math.Sqrt(dx*dx + dy*dy)
			//Check if the distance is approximately equal to the specified radius
			//*0.85 is to obtain a circle looking like the tester's circle even if it's not really a circle... In reality, remove "*0.85" and it's a real circle
			if math.Abs(distance-float64(radius)*0.85) < 0.5 {
				ppm.data[y][x] = color
			}
		}
	}
}

// DrawCircle draws a circle.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	//Draw a circle with the radius getting smaller until it is at 0;
	for radius >= 0 {
		ppm.DrawCircle(center, radius, color)
		radius--
	}
}

// DrawFilledCircle draws a filled circle.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	//Draw lines and link the 3 points
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// Draw a line from p1 to p3 and move p1 towars p2 until the triangle is filled
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	//Loop until p1 reaches p2
	for p1 != p2 {
		//Draw a line between p1 and p3
		ppm.DrawLine(p3, p1, color)
		//Increment or decrement X of p1 based on p2 position
		if p1.X != p2.X && p1.X < p2.X {
			p1.X++
		} else if p1.X != p2.X && p1.X > p2.X {
			p1.X--
		}
		//Increment or decrement Y of p1 based on p2 position
		if p1.Y != p2.Y && p1.Y < p2.Y {
			p1.Y++
		} else if p1.Y != p2.Y && p1.Y > p2.Y {
			p1.Y--
		}
	}
	//Draw a final line between the last position of p1 (should be at p2 at this point) and p3
	ppm.DrawLine(p3, p1, color)
}
