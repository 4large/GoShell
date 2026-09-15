package cat

import (
	"bufio"
	"fmt"
	"io"
)

// Core function, readers input stream, applies flags, writes to writer,
// maintains structure IE line number or line state
func reader(r io.Reader, writer *bufio.Writer, flags []bool, lineNum *int,
	newline *bool) error {
	bufioReader := bufio.NewReader(r)
	squeezed := false
	printingChar := true
	for {
		//read single byte at a time
		b, err := bufioReader.ReadByte()
		if err != nil {
			if err == io.EOF {
				//end of the stream
				return nil
			} else {
				return err
			}
		}
		//apply flags

		if *newline {
			//-s
			if flags[2] && b == '\n' && !squeezed {
				squeezed = true
			} else if b != '\n' && squeezed {
				squeezed = false
			} else if squeezed {
				continue
			}

			//-b or -n
			if flags[1] && b != '\n' {
				fmt.Fprintf(writer, "%6d\t", *lineNum)
				*lineNum++
			} else if flags[0] && !flags[1] {
				fmt.Fprintf(writer, "%6d\t", *lineNum)
				*lineNum++
			}

		}

		//-E
		if flags[3] && b == '\n' {
			writer.WriteByte('$')
		}

		//-T
		if flags[4] && b == '\t' {
			writer.WriteByte('^')
			writer.WriteByte('I')
			continue
		}

		//-v
		if flags[5] {
			if b != '\t' && b != '\n' {
				printingChar = v(writer, b)
				if !printingChar {
					continue
				}
			}
		}

		//-e, equivalent to -vE
		if flags[6] && !flags[3] {
			if b == '\n' {
				writer.WriteByte('$')
			} else if b != '\t' {
				printingChar = v(writer, b)
				if !printingChar {
					continue
				}
			}
		}

		//-t, equivalent to -vT
		if flags[7] {
			if b == '\t' {
				writer.WriteByte('^')
				writer.WriteByte('I')
				continue
			} else if b != '\n' {
				printingChar = v(writer, b)
				if !printingChar {
					continue
				}
			}
		}

		//-A, equivalent to -vET
		if flags[8] && !flags[3] {
			switch b {
			case '\t':
				writer.WriteByte('^')
				writer.WriteByte('I')
				continue
			case '\n':
				writer.WriteByte('$')
			default:
				printingChar = v(writer, b)
				if !printingChar {
					continue
				}
			}
		}

		writer.WriteByte(b)
		*newline = (b == '\n')

		if b == '\n' {
			writer.Flush()
		}
	}
}

func v(writer *bufio.Writer, b byte) bool {
	if b <= 31 && b != 9 && b != 10 {
		writer.WriteByte('^')
		writer.WriteByte(b + 64)
		return false
	} else if b == 127 {
		writer.WriteByte('^')
		writer.WriteByte('?')
		return false
	} else if b >= 128 {
		writer.WriteByte('M')
		writer.WriteByte('-')
		return v(writer, b&0x7F)
	}
	return true
}
