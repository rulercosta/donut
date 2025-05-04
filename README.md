<div align="center">
  <a href="https://golang.org/">
    <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white&style=for-the-badge" alt="Go Version"/>
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge&logo=opensourceinitiative&logoColor=white" alt="License"/>
  </a>
</div>

# 🍩 Donut - ASCII Donut Animation in Go

This project is an idiomatic Go implementation that maintains the spirit and logic of the original [donut.c](https://www.a1k0n.net/2006/09/15/obfuscated-c-donut.html), with improvements for Go’s syntax, structure, and best practices. The project delivers the same mesmerizing ASCII donut animation in your terminal, while benefiting from Go’s modern language features and code clarity.

Building upon the original C implementation, this version offers improvements through Go's strong concurrency support and cleaner code structure, while maintaining the same mesmerizing visual effect with no external dependencies.

><center><i>Inspired by <a href="https://www.a1k0n.net/2011/07/20/donut-math.html">Andy Sloane's donut.c math</a>.</i></center>

## How To Run

Run the following to see the donut in action:

```
go run donut.go
```

<div align="center">
  <img src="donut.gif" alt="Donut Animation" width="100%"/>
</div>

### How It Works

- Uses trigonometry and matrix math to project a 3D torus onto a 2D ASCII grid
- Calculates luminance for each point to simulate shading
- Double buffers output for smooth animation
- Tweak animation parameters by editing the constants at the top of `donut.go`

## License

MIT License. See [LICENSE](LICENSE) for details.
