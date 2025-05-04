// donut.go - ASCII donut animation in Go
// Inspired by the classic donut.c
package main

import (
	"fmt"
	"math"
	"time"
)

const (
	// Screen and buffer dimensions
	screenWidth 		= 		80
	screenHeight		= 		22
	bufferSize  		= 		screenWidth * screenHeight

	// Donut geometry and animation
	phiStep           	= 		0.07
	thetaStep         	= 		0.02
	frameDelay        	= 		30 * time.Millisecond
	radiusMajor       	= 		2.0
	radiusMinor       	= 		1.0
	distanceOffset    	= 		5.0
	xScale            	= 		30.0
	yScale            	= 		15.0
	rotationAIncrement	= 		0.04
	rotationBIncrement	= 		0.02

	// Luminance
	luminanceLevels 	=		8
	luminanceChars  	=		".,-~:;=!*#$@"
	
	// Terminal escape sequences
	hideCursor 			=		"\x1b[?25l"
	showCursor 			=		"\x1b[?25h"
	clearScreen			=		"\x1b[2J"
	homeCursor 			=		"\x1b[H"
)

func main() {
	fmt.Print(hideCursor)
	defer fmt.Print(showCursor)
	fmt.Print(clearScreen)
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
		rotationA += rotationAIncrement
		rotationB += rotationBIncrement
	}
}

func clearBuffers(outputBuffer []byte, zBuffer []float64) {
	for i := range outputBuffer {
		outputBuffer[i] = ' '
		zBuffer[i] = 0
	}
}

func computeFrame(outputBuffer []byte, zBuffer []float64, rotationA, rotationB float64) {
	for anglePhi := 0.0; anglePhi < 2*math.Pi; anglePhi += phiStep {
		for angleTheta := 0.0; angleTheta < 2*math.Pi; angleTheta += thetaStep {
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
	sinTheta, cosTheta := getSinCos(angleTheta)
	sinPhi, cosPhi := getSinCos(anglePhi)
	sinRotationA, cosRotationA := getSinCos(rotationA)
	sinRotationB, cosRotationB := getSinCos(rotationB)
	h := getMajorRadiusOffset(cosPhi)
	invDepth := getInverseDepth(sinTheta, h, sinRotationA, sinPhi, cosRotationA)
	temp := getTemp(sinTheta, h, cosRotationA, sinPhi, sinRotationA)
	x := getScreenX(cosTheta, h, cosRotationB, temp, sinRotationB, invDepth)
	y := getScreenY(cosTheta, h, sinRotationB, temp, cosRotationB, invDepth)
	bufferOffset := x + screenWidth*y
	return x, y, bufferOffset, invDepth
}

func calculateLuminanceIndex(angleTheta, anglePhi, rotationA, rotationB float64) int {
	sinTheta, cosTheta := getSinCos(angleTheta)
	sinPhi, cosPhi := getSinCos(anglePhi)
	sinRotationA, cosRotationA := getSinCos(rotationA)
	sinRotationB, cosRotationB := getSinCos(rotationB)
	luminance := getLuminance(sinTheta, cosTheta, sinPhi, cosPhi, sinRotationA, cosRotationA, sinRotationB, cosRotationB)
	return int(float64(luminanceLevels) * luminance)
}

func getSinCos(angle float64) (float64, float64) {
	return math.Sin(angle), math.Cos(angle)
}

func getMajorRadiusOffset(cosPhi float64) float64 {
	return cosPhi + radiusMajor
}

func getInverseDepth(sinTheta, h, sinRotationA, sinPhi, cosRotationA float64) float64 {
	return 1 / (sinTheta*h*sinRotationA + sinPhi*cosRotationA + distanceOffset)
}

func getTemp(sinTheta, h, cosRotationA, sinPhi, sinRotationA float64) float64 {
	return sinTheta*h*cosRotationA - sinPhi*sinRotationA
}

func getScreenX(cosTheta, h, cosRotationB, temp, sinRotationB, invDepth float64) int {
	return int(float64(screenWidth/2) + xScale*invDepth*(cosTheta*h*cosRotationB-temp*sinRotationB))
}

func getScreenY(cosTheta, h, sinRotationB, temp, cosRotationB, invDepth float64) int {
	return int(float64(screenHeight/2) + yScale*invDepth*(cosTheta*h*sinRotationB+temp*cosRotationB))
}

func getLuminance(sinTheta, cosTheta, sinPhi, cosPhi, sinRotationA, cosRotationA, sinRotationB, cosRotationB float64) float64 {
	return (sinPhi*sinRotationA-sinTheta*cosPhi*cosRotationA)*cosRotationB - sinTheta*cosPhi*sinRotationA - sinPhi*cosRotationA - cosTheta*cosPhi*sinRotationB
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
	fmt.Print(homeCursor)
	for bufferIndex := range bufferSize {
		if bufferIndex%screenWidth == screenWidth-1 {
			fmt.Printf("%c\n", outputBuffer[bufferIndex])
		} else {
			fmt.Printf("%c", outputBuffer[bufferIndex])
		}
	}
}
