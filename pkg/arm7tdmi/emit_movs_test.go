package arm7tdmi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/edipermadi/gojit/pkg/arm7tdmi"
	"github.com/stretchr/testify/assert"
)

func TestEmit_Movs(t *testing.T) {
	t.Run("NoPreprocess", func(t *testing.T) {
		start := 0
		span := 1 << 6
		end := start + span
		for i := start; i < end; i++ {
			rd := i & 0x07
			rs := (i >> 3) & 0x07
			actual := strings.TrimSpace(arm7tdmi.Emit(uint16(i)))
			expected := fmt.Sprintf("ctx->r%d = ctx->r%d;", rd, rs)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("LslPreprocess", func(t *testing.T) {
		start := (0 << 11) + (1 << 6)
		span := 1 << 10
		end := start + span
		for i := start; i < end; i++ {
			offset := (i >> 6) & 0x1f
			rd := i & 0x07
			rs := (i >> 3) & 0x07
			actual := strings.TrimSpace(arm7tdmi.Emit(uint16(i)))
			expected := fmt.Sprintf("ctx->r%d = op_lsl(ctx->r%d, %d);", rd, rs, offset)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("LsrPreprocess", func(t *testing.T) {
		start := (1 << 11) + (1 << 6)
		span := 1 << 10
		end := start + span
		for i := start; i < end; i++ {
			offset := (i >> 6) & 0x1f
			rd := i & 0x07
			rs := (i >> 3) & 0x07
			actual := strings.TrimSpace(arm7tdmi.Emit(uint16(i)))
			expected := fmt.Sprintf("ctx->r%d = op_lsr(ctx->r%d, %d);", rd, rs, offset)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("AsrPreprocess", func(t *testing.T) {
		start := (2 << 11) + (1 << 6)
		span := 1 << 10
		end := start + span
		for i := start; i < end; i++ {
			offset := (i >> 6) & 0x1f
			rd := i & 0x07
			rs := (i >> 3) & 0x07
			actual := strings.TrimSpace(arm7tdmi.Emit(uint16(i)))
			expected := fmt.Sprintf("ctx->r%d = op_asr(ctx->r%d, %d);", rd, rs, offset)
			assert.Equal(t, expected, actual)
		}
	})
}

func TestEmit_AddSub(t *testing.T) {
	t.Run("AddReg", func(t *testing.T) {
		start := (3 << 11) + (0 << 10) + (0 << 9)
		span := 1 << 8
		end := start + span
		for i := start; i < end; i++ {
			rd := i & 0x07
			rs := (i >> 3) & 0x07
			rn := (i >> 6) & 0x07
			actual := strings.TrimSpace(arm7tdmi.Emit(uint16(i)))
			expected := fmt.Sprintf("ctx->r%d = ctx->r%d + ctx->r%d;", rd, rs, rn)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("AddImm", func(t *testing.T) {
		start := (3 << 11) + (1 << 10) + (0 << 9)
		span := 1 << 8
		end := start + span
		for i := start; i < end; i++ {
			rd := i & 0x07
			rs := (i >> 3) & 0x07
			rn := (i >> 6) & 0x07
			actual := strings.TrimSpace(arm7tdmi.Emit(uint16(i)))
			expected := fmt.Sprintf("ctx->r%d = ctx->r%d + %d;", rd, rs, rn)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("SubReg", func(t *testing.T) {
		start := (3 << 11) + (0 << 10) + (1 << 9)
		span := 1 << 8
		end := start + span
		for i := start; i < end; i++ {
			rd := i & 0x07
			rs := (i >> 3) & 0x07
			rn := (i >> 6) & 0x07
			actual := strings.TrimSpace(arm7tdmi.Emit(uint16(i)))
			expected := fmt.Sprintf("ctx->r%d = ctx->r%d - ctx->r%d;", rd, rs, rn)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("SubImm", func(t *testing.T) {
		start := (3 << 11) + (1 << 10) + (1 << 9)
		span := 1 << 8
		end := start + span
		for i := start; i < end; i++ {
			rd := i & 0x07
			rs := (i >> 3) & 0x07
			rn := (i >> 6) & 0x07
			actual := strings.TrimSpace(arm7tdmi.Emit(uint16(i)))
			expected := fmt.Sprintf("ctx->r%d = ctx->r%d - %d;", rd, rs, rn)
			assert.Equal(t, expected, actual)
		}
	})
}
