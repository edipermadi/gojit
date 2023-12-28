package arm7tdmi

import (
	"fmt"
	"io"
	"strings"
)

func Emit(opcode uint16) string {
	var builder strings.Builder
	emit(&builder, opcode)
	return builder.String()
}

func emit(writer io.Writer, opcode uint16) {
	op := opcode >> 11
	if op < 3 {
		emitMovs(writer, opcode)
	} else if op == 3 {
		emitAddSub(writer, opcode)
	}
}

func EmitFunc(opcode uint16) string {
	var builder strings.Builder
	_, _ = fmt.Fprintf(&builder, "static void op_%04x(arm7tdmi_thumb_ctx_t *ctx){\n", opcode)

	emit(&builder, opcode)

	_, _ = fmt.Fprintf(&builder, "}\n")
	return builder.String()
}

func bitval(opcode uint16, pos int, bits int) uint16 {
	return (opcode >> pos) & ((1 << bits) - 1)
}

func emitMovs(writer io.Writer, opcode uint16) {
	op := bitval(opcode, 11, 2)
	if op == 3 {
		return
	}

	offset := bitval(opcode, 6, 5)
	rs := bitval(opcode, 3, 3)
	rd := bitval(opcode, 0, 3)

	if offset == 0 {
		_, _ = fmt.Fprintf(writer, "\tctx->r%d = ctx->r%d;\n", rd, rs)
	} else if op == 0 {
		_, _ = fmt.Fprintf(writer, "\tctx->r%d = op_lsl(ctx->r%d, %d);\n", rd, rs, offset)
	} else if op == 1 {
		_, _ = fmt.Fprintf(writer, "\tctx->r%d = op_lsr(ctx->r%d, %d);\n", rd, rs, offset)
	} else {
		_, _ = fmt.Fprintf(writer, "\tctx->r%d = op_asr(ctx->r%d, %d);\n", rd, rs, offset)
	}
}

func emitAddSub(writer io.Writer, opcode uint16) {
	imm := bitval(opcode, 10, 1)
	sub := bitval(opcode, 9, 1)
	rn := bitval(opcode, 6, 3)
	rs := bitval(opcode, 3, 3)
	rd := bitval(opcode, 0, 3)

	if sub > 0 {
		if imm > 0 {
			_, _ = fmt.Fprintf(writer, "\tctx->r%d = ctx->r%d - %d;\n", rd, rs, rn)
		} else {
			_, _ = fmt.Fprintf(writer, "\tctx->r%d = ctx->r%d - ctx->r%d;\n", rd, rs, rn)
		}
	} else if imm > 0 {
		_, _ = fmt.Fprintf(writer, "\tctx->r%d = ctx->r%d + %d;\n", rd, rs, rn)
	} else {
		_, _ = fmt.Fprintf(writer, "\tctx->r%d = ctx->r%d + ctx->r%d;\n", rd, rs, rn)
	}
}
