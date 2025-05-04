// donut.go - ASCII donut animation in Go
// Inspired by the classic donut.c
package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	fmt.Print("\x1b[?25l")
	defer fmt.Print("\x1b[?25h")

	rotationA, rotationB := 0.0, 0.0
	var angleTheta, anglePhi float64
	var bufferIndex int
	zBuffer := make([]float64, 1760)
	outputBuffer := make([]byte, 1760)
	fmt.Print("\x1b[2J")
	for {
		for i := range outputBuffer {
			outputBuffer[i] = ' '
			zBuffer[i] = 0
		}
		for anglePhi = 0; anglePhi < 6.28; anglePhi += 0.07 {
			for angleTheta = 0; angleTheta < 6.28; angleTheta += 0.02 {
				sinTheta := math.Sin(angleTheta)
				cosPhi := math.Cos(anglePhi)
				sinRotationA := math.Sin(rotationA)
				sinPhi := math.Sin(anglePhi)
				cosRotationA := math.Cos(rotationA)
				h := cosPhi + 2
				invDepth := 1 / (sinTheta*h*sinRotationA + sinPhi*cosRotationA + 5)
				cosTheta := math.Cos(angleTheta)
				cosRotationB := math.Cos(rotationB)
				sinRotationB := math.Sin(rotationB)
				temp := sinTheta*h*cosRotationA - sinPhi*sinRotationA
				x := int(40 + 30*invDepth*(cosTheta*h*cosRotationB-temp*sinRotationB))
				y := int(12 + 15*invDepth*(cosTheta*h*sinRotationB+temp*cosRotationB))
				bufferOffset := x + 80*y
				luminanceIndex := int(8 * ((sinPhi*sinRotationA-sinTheta*cosPhi*cosRotationA)*cosRotationB - sinTheta*cosPhi*sinRotationA - sinPhi*cosRotationA - cosTheta*cosPhi*sinRotationB))
				if 22 > y && y > 0 && x > 0 && 80 > x && invDepth > zBuffer[bufferOffset] {
					zBuffer[bufferOffset] = invDepth
					luminanceChars := ".,-~:;=!*#$@"
					if luminanceIndex > 0 && luminanceIndex < len(luminanceChars) {
						outputBuffer[bufferOffset] = luminanceChars[luminanceIndex]
					} else {
						outputBuffer[bufferOffset] = luminanceChars[0]
					}
				}
			}
		}
		fmt.Print("\x1b[H")
		for bufferIndex = range 1760 {
			if bufferIndex%80 == 79 {
				fmt.Printf("%c\n", outputBuffer[bufferIndex])
			} else {
				fmt.Printf("%c", outputBuffer[bufferIndex])
			}
		}
		time.Sleep(30 * time.Millisecond)
		rotationA += 0.04
		rotationB += 0.02
	}
}
