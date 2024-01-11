package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Structure pour représenter l'image PBM
type PBM struct {
	Data          [][]bool
	Width, Height int
	MagicNumber   string
}

// Fonction pour lire une image PBM à partir d'un fichier
func ReadPBM(filename string) (*PBM, error) {
	// Ouvrir le fichier en mode lecture
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Lire la magie number (P1 ou P4 pour PBM)
	if !scanner.Scan() {
		return nil, fmt.Errorf("Erreur de lecture du Magic Number")
	}
	magicNumber := scanner.Text()
	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, fmt.Errorf("Format d'image non pris en charge. Magic Number : %s", magicNumber)
	}

	// Lire la largeur et la hauteur
	if !scanner.Scan() {
		return nil, fmt.Errorf("Erreur de lecture des dimensions : aucune ligne de dimensions")
	}
	dimensionsLine := scanner.Text()

	var width, height int
	_, err = fmt.Sscanf(dimensionsLine, "%d %d", &width, &height)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la lecture des dimensions : %v", err)
	}

	// Initialiser la matrice pour stocker les pixels
	data := make([][]bool, height)
	for i := range data {
		data[i] = make([]bool, width)
	}

	// Remplir la matrice avec les pixels de l'image
	for i := 0; i < height && scanner.Scan(); i++ {
		line := strings.Fields(scanner.Text())
		for j := range line {
			if j < width {
				if magicNumber == "P1" {
					if line[j] != "0" && line[j] != "1" {
						return nil, fmt.Errorf("Valeur de pixel invalide à la ligne %d, colonne %d", i, j)
					}
					data[i][j] = line[j] == "1"
				} else if magicNumber == "P4" {
					// Format binaire P4
					data[i][j] = line[j] == "1"
				}
			} else {
				return nil, fmt.Errorf("Trop de valeurs dans la ligne %d", i)
			}
		}
	}

	// Créer une structure pour représenter l'image
	image := &PBM{
		Data:        data,
		Width:       width,
		Height:      height,
		MagicNumber: magicNumber,
	}

	return image, nil
}

func (pbm *PBM) Taille() (int, int) {
	return pbm.Width, pbm.Height
}
func (pbm *PBM) At(x, y int) bool {
	return pbm.Data[y][x]
}
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.Data[y][x] = value
}

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write the Magic Number to the file
	fmt.Fprintf(writer, "%s\n", pbm.MagicNumber)

	// Write the dimensions to the file
	fmt.Fprintf(writer, "%d %d\n", pbm.Width, pbm.Height)

	// Write the image data to the file
	for i := 0; i < pbm.Height; i++ {
		for j := 0; j < pbm.Width; j++ {
			if pbm.MagicNumber == "P1" {
				// Write 0 or 1 based on the boolean value
				if pbm.Data[i][j] {
					fmt.Fprint(writer, "1 ")
				} else {
					fmt.Fprint(writer, "0 ")
				}
			} else if pbm.MagicNumber == "P4" {
				// Write 0 or 1 based on the boolean value
				if pbm.Data[i][j] {
					fmt.Fprint(writer, "1")
				} else {
					fmt.Fprint(writer, "0")
				}
			}
		}
		fmt.Fprintln(writer)
	}

	return writer.Flush()
}
func (pbm *PBM) Invert() {
	for i := 0; i < pbm.Height; i++ {
		for j := 0; j < pbm.Width; j++ {
			// Invert the color by flipping the boolean value
			pbm.Data[i][j] = !pbm.Data[i][j]
		}
	}
}

func main() {
	filename := "test1.pbm" // Specify the correct path to the PBM file

	image, err := ReadPBM(filename)
	if err != nil {
		fmt.Println("Error reading the image:", err)
		return
	}

	fmt.Println("Original Image:")
	fmt.Printf("Magic Number: %s\n", image.MagicNumber)
	width, height := image.Taille()
	fmt.Printf("Width: %d\n", width)
	fmt.Printf("Height: %d\n", height)
	fmt.Printf("Data: %v\n", image.Data)

	// Invert the colors of the image
	image.Invert()

	fmt.Println("Inverted Image:")
	fmt.Printf("Magic Number: %s\n", image.MagicNumber)
	width, height = image.Taille()
	fmt.Printf("Width: %d\n", width)
	fmt.Printf("Height: %d\n", height)
	fmt.Printf("Data: %v\n", image.Data)

	saveFilename := "inverted_test1.pbm"

	// Save the inverted image to a new file
	err = image.Save(saveFilename)
	if err != nil {
		fmt.Println("Error saving the inverted image:", err)
		return
	}

	fmt.Println("Inverted image saved successfully.")
}
