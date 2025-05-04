// donut.go - ASCII donut animation in Go
// Inspired by the classic donut.c
package main

import (
	"fmt"
	"math"
	"time"
)

const (
	screenWidth    = 80
	screenHeight   = 22
	bufferSize     = 1760
	phiStep        = 0.07
	thetaStep      = 0.02
	frameDelay     = 30 * time.Millisecond
	luminanceChars = ".,-~:;=!*#$@"
)

func main() {
	fmt.Print("\x1b[?25l")
	defer fmt.Print("\x1b[?25h")
	fmt.Print("\x1b[2J")
	animateDonut()
}

func animateDonut() {
	rotationA, rotationB := 0.0, 0.0
	zBuffer := make([]float64, bufferSize)
	outputBuffer := make([]byte, bufferSize)
	for {
		clearBuffers(outputBuffer, zBuffer)
		computeFrame(outputBuffer, zBuffer, rotationA, rotationB)
		drawFrame(outputBuffer)
		time.Sleep(frameDelay)
		rotationA += 0.04
		rotationB += 0.02
	}
}

func clearBuffers(outputBuffer []byte, zBuffer []float64) {
	for i := range outputBuffer {
		outputBuffer[i] = ' '
		zBuffer[i] = 0
	}
}

func computeFrame(outputBuffer []byte, zBuffer []float64, rotationA, rotationB float64) {
	for anglePhi := 0.0; anglePhi < 6.28; anglePhi += phiStep {
		for angleTheta := 0.0; angleTheta < 6.28; angleTheta += thetaStep {
			x, y, bufferOffset, invDepth := project3DTo2D(angleTheta, anglePhi, rotationA, rotationB)
			luminanceIndex := calculateLuminanceIndex(angleTheta, anglePhi, rotationA, rotationB)
			if isPointVisible(x, y, invDepth, zBuffer, bufferOffset) {
				zBuffer[bufferOffset] = invDepth
				outputBuffer[bufferOffset] = getLuminanceChar(luminanceIndex)
			}
		}
	}
}

func project3DTo2D(angleTheta, anglePhi, rotationA, rotationB float64) (int, int, int, float64) {
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
	x := int(float64(screenWidth/2) + 30*invDepth*(cosTheta*h*cosRotationB-temp*sinRotationB))
	y := int(float64(screenHeight/2) + 15*invDepth*(cosTheta*h*sinRotationB+temp*cosRotationB))
	bufferOffset := x + screenWidth*y
	return x, y, bufferOffset, invDepth
}

func calculateLuminanceIndex(angleTheta, anglePhi, rotationA, rotationB float64) int {
	sinTheta := math.Sin(angleTheta)
	cosPhi := math.Cos(anglePhi)
	sinRotationA := math.Sin(rotationA)
	sinPhi := math.Sin(anglePhi)
	cosRotationA := math.Cos(rotationA)
	cosRotationB := math.Cos(rotationB)
	luminance := (sinPhi*sinRotationA-sinTheta*cosPhi*cosRotationA)*cosRotationB - sinTheta*cosPhi*sinRotationA - sinPhi*cosRotationA - math.Cos(angleTheta)*cosPhi*math.Sin(rotationB)
	return int(8 * luminance)
}

func getLuminanceChar(luminanceIndex int) byte {
	if luminanceIndex > 0 && luminanceIndex < len(luminanceChars) {
		return luminanceChars[luminanceIndex]
	}
	return luminanceChars[0]
}

func isPointVisible(x, y int, invDepth float64, zBuffer []float64, bufferOffset int) bool {
	return screenHeight > y && y > 0 && x > 0 && screenWidth > x && invDepth > zBuffer[bufferOffset]
}

func drawFrame(outputBuffer []byte) {
	fmt.Print("\x1b[H")
	for bufferIndex := range bufferSize {
		if bufferIndex%screenWidth == screenWidth-1 {
			fmt.Printf("%c\n", outputBuffer[bufferIndex])
		} else {
			fmt.Printf("%c", outputBuffer[bufferIndex])
		}
	}
}
