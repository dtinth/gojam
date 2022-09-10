package jitterbuffer

/*
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn jitter_buffer_works() {
        let mut buffer = JitterBuffer::new(3);
        assert_eq!(buffer.put_in("A", 20), None);
        assert_eq!(buffer.put_in("B", 21), None);
        assert_eq!(buffer.put_in("C", 22), None);
        assert_eq!(buffer.put_in("D", 23), Some("A"));
        assert_eq!(buffer.put_in("E", 24), Some("B"));
        assert_eq!(buffer.put_in("F", 25), Some("C"));
    }

    #[test]
    fn jitter_buffer_can_handle_jitter() {
        let mut buffer = JitterBuffer::new(3);
        assert_eq!(buffer.put_in("C", 22), None);
        assert_eq!(buffer.put_in("B", 21), None);
        assert_eq!(buffer.put_in("A", 20), None);
        assert_eq!(buffer.put_in("E", 24), Some("A"));
        assert_eq!(buffer.put_in("F", 25), Some("B"));
        assert_eq!(buffer.put_in("D", 23), Some("C"));
    }

    #[test]
    fn jitter_buffer_works_at_u8_boundary() {
        let mut buffer = JitterBuffer::new(3);
        assert_eq!(buffer.put_in("A", 253), None);
        assert_eq!(buffer.put_in("D", 0), None);
        assert_eq!(buffer.put_in("C", 255), None);
        assert_eq!(buffer.put_in("B", 254), Some("A"));
        assert_eq!(buffer.put_in("F", 2), Some("B"));
        assert_eq!(buffer.put_in("E", 1), Some("C"));
    }
}
*/

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func test(t *testing.T, buffer *JitterBuffer, input string, sequenceNum int, expected string) {
	inputBytes := []byte(input)
	outputBytes := buffer.PutIn(inputBytes, uint8(sequenceNum))
	if expected == "" {
		assert.Nil(t, outputBytes)
	} else {
		assert.Equal(t, expected, string(outputBytes))
	}
}

func TestJitterBufferWorks(t *testing.T) {
	buffer := NewJitterBuffer(3)
	test(t, buffer, "A", 20, "")
	test(t, buffer, "B", 21, "")
	test(t, buffer, "C", 22, "")
	test(t, buffer, "D", 23, "A")
	test(t, buffer, "E", 24, "B")
	test(t, buffer, "F", 25, "C")
}

func TestJitterBufferCanHandleJitter(t *testing.T) {
	buffer := NewJitterBuffer(3)
	test(t, buffer, "C", 22, "")
	test(t, buffer, "B", 21, "")
	test(t, buffer, "A", 20, "")
	test(t, buffer, "E", 24, "A")
	test(t, buffer, "F", 25, "B")
	test(t, buffer, "D", 23, "C")
}

func TestJitterBufferWorksAtU8Boundary(t *testing.T) {
	buffer := NewJitterBuffer(3)
	test(t, buffer, "A", 253, "")
	test(t, buffer, "D", 0, "")
	test(t, buffer, "C", 255, "")
	test(t, buffer, "B", 254, "A")
	test(t, buffer, "F", 2, "B")
	test(t, buffer, "E", 1, "C")
}
