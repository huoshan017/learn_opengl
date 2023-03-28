package gl

import (
	"log"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

func Init() error {
	return gl.Init()
}

func NeedVao() bool {
	return true
}

func GetError() uint32 {
	return gl.GetError()
}

func Viewport(x, y, width, height int32) {
	gl.Viewport(x, y, width, height)
}

func ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func Clear(flag uint32) {
	gl.Clear(flag)
}

func Disable(flag uint32) {
	gl.Disable(flag)
}

func Enable(flag uint32) {
	gl.Enable(flag)
}

func Scissor(x, y, w, h int32) {
	gl.Scissor(x, y, w, h)
}

func DepthMask(flag bool) {
	gl.DepthMask(flag)
}

func ColorMask(r, g, b, a bool) {
	gl.ColorMask(r, g, b, a)
}

func BlendFunc(src, dst uint32) {
	gl.BlendFunc(src, dst)
}

func DepthFunc(fn uint32) {
	gl.DepthFunc(fn)
}

func StencilMask(mask uint32) {
	gl.StencilMask(mask)
}

func StencilFunc(xfunc uint32, ref int32, mask uint32) {
	gl.StencilFunc(xfunc, ref, mask)
}

func StencilOp(fail, zfail, zpass uint32) {
	gl.StencilOp(fail, zfail, zpass)
}

// vao

func GenVertexArrays(n int32, array *uint32) {
	gl.GenVertexArrays(n, array)
}

func BindVertexArray(array uint32) {
	gl.BindVertexArray(array)
}

func DeleteVertexArrays(n int32, array *uint32) {
	gl.DeleteVertexArrays(n, array)
}

// program & shader

func CreateProgram() uint32 {
	return gl.CreateProgram()
}

func DeleteProgram(program uint32) {
	gl.DeleteProgram(program)
}

func AttachShader(program, shader uint32) {
	gl.AttachShader(program, shader)
}

func LinkProgram(program uint32) {
	gl.LinkProgram(program)
}

func UseProgram(program uint32) {
	gl.UseProgram(program)
}

func GetProgramiv(program, pname uint32, params *int32) {
	gl.GetProgramiv(program, pname, params)
}

// TODO 原来的视线中 buf = logLength + 1，需要测试这种情况
func GetProgramInfoLog(program uint32) string {
	var logLength int32
	GetProgramiv(program, INFO_LOG_LENGTH, &logLength)

	if logLength == 0 {
		return ""
	}

	buf := make([]uint8, logLength)
	gl.GetProgramInfoLog(program, logLength, nil, &buf[0])
	return string(buf)
}

func CreateShader(xtype uint32) uint32 {
	return gl.CreateShader(xtype)
}

func ShaderSource(shader uint32, src string) {
	cstr, free := gl.Strs(src)
	gl.ShaderSource(shader, 1, cstr, nil)
	free()
}

func CompileShader(shader uint32) {
	gl.CompileShader(shader)
}

func GetShaderiv(shader uint32, pname uint32, params *int32) {
	gl.GetShaderiv(shader, pname, params)
}

func GetShaderInfoLog(shader uint32) string {
	var logLength int32
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		logLength = 128
	}
	buf := make([]uint8, logLength)
	log.Printf("logLength = %v", logLength)
	gl.GetShaderInfoLog(shader, logLength, nil, &buf[0])
	return string(buf)
}

func DeleteShader(shader uint32) {
	gl.DeleteShader(shader)
}

// buffer & draw

func GenBuffers(n int32, buffers *uint32) {
	gl.GenBuffers(n, buffers)
}

func BufferData(target uint32, size int, data unsafe.Pointer, usage uint32) {
	gl.BufferData(target, size, data, usage)
}

func BufferSubData(target uint32, offset, size int, data unsafe.Pointer) {
	gl.BufferSubData(target, offset, size, data)
}

func BindBuffer(target, buffer uint32) {
	gl.BindBuffer(target, buffer)
}

func DeleteBuffers(n int32, buffers *uint32) {
	gl.DeleteBuffers(n, buffers)
}

func BindBufferRange(target, index, buffer uint32, offset, size int) {
	gl.BindBufferRange(target, index, buffer, offset, size)
}

func GenFramebuffers(n int32, buffers *uint32) {
	gl.GenFramebuffers(n, buffers)
}

func BindFramebuffer(target, buffer uint32) {
	gl.BindFramebuffer(target, buffer)
}

func FramebufferTexture1D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32) {
	gl.FramebufferTexture1D(target, attachment, textarget, texture, level)
}

func FramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32) {
	gl.FramebufferTexture2D(target, attachment, textarget, texture, level)
}

func FramebufferTexture3D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32, zoffset int32) {
	gl.FramebufferTexture3D(target, attachment, textarget, texture, level, zoffset)
}

func FramebufferRenderbuffer(target uint32, attachment uint32, renderbuffertarget uint32, renderbuffer uint32) {
	gl.FramebufferRenderbuffer(target, attachment, renderbuffertarget, renderbuffer)
}

func GenRenderbuffers(n int32, buffers *uint32) {
	gl.GenRenderbuffers(n, buffers)
}

func BindRenderbuffer(target, buffer uint32) {
	gl.BindRenderbuffer(target, buffer)
}

func RenderbufferStorage(target, internalformat uint32, width, height int32) {
	gl.RenderbufferStorage(target, internalformat, width, height)
}

func RenderbufferStorageMultisample(target uint32, samples int32, internalformat uint32, width int32, height int32) {
	gl.RenderbufferStorageMultisample(target, samples, internalformat, width, height)
}

func CheckFramebufferStatus(target uint32) uint32 {
	return gl.CheckFramebufferStatus(target)
}

func TexImage2DMultisample(target uint32, samples int32, internalformat uint32, width int32, height int32, fixedsamplelocations bool) {
	gl.TexImage2DMultisample(target, samples, internalformat, width, height, fixedsamplelocations)
}

func DrawElements(mode uint32, count int32, typ uint32, offset int) {
	gl.DrawElements(mode, count, typ, gl.PtrOffset(offset))
}

func DrawArrays(mode uint32, first, count int32) {
	gl.DrawArrays(mode, first, count)
}

func DrawArraysInstanced(mode uint32, first, count, instancecount int32) {
	gl.DrawArraysInstanced(mode, first, count, instancecount)
}

func DrawElementsInstanced(mode uint32, count int32, xtype uint32, indices unsafe.Pointer, instancecount int32) {
	gl.DrawElementsInstanced(mode, count, xtype, indices, instancecount)
}

func DrawBuffer(buf uint32) {
	gl.DrawBuffer(buf)
}

func DrawBuffers(n int32, bufs *uint32) {
	gl.DrawBuffers(n, bufs)
}

func ReadBuffer(src uint32) {
	gl.ReadBuffer(src)
}

// uniform

func GetUniformLocation(program uint32, name string) int32 {
	return gl.GetUniformLocation(program, gl.Str(name))
}

func Uniform1i(loc, v int32) {
	gl.Uniform1i(loc, v)
}

func Uniform1iv(loc, num int32, v *int32) {
	gl.Uniform1iv(loc, num, v)
}

func Uniform1f(location int32, v0 float32) {
	gl.Uniform1f(location, v0)
}

func Uniform2f(location int32, v0, v1 float32) {
	gl.Uniform2f(location, v0, v1)
}

func Uniform3f(location int32, v0, v1, v2 float32) {
	gl.Uniform3f(location, v0, v1, v2)
}

func Uniform4f(location int32, v0, v1, v2, v3 float32) {
	gl.Uniform4f(location, v0, v1, v2, v3)
}

func Uniform1fv(loc, num int32, v *float32) {
	gl.Uniform1fv(loc, num, v)
}

func Uniform2fv(loc, num int32, v *float32) {
	gl.Uniform2fv(loc, num, v)
}

func Uniform3fv(loc, num int32, v *float32) {
	gl.Uniform3fv(loc, num, v)
}

func Uniform4fv(loc, num int32, v *float32) {
	gl.Uniform4fv(loc, num, v)
}

func UniformMatrix2fv(loc, num int32, t bool, v *float32) {
	gl.UniformMatrix2fv(loc, num, t, v)
}

func UniformMatrix3fv(loc, num int32, t bool, v *float32) {
	gl.UniformMatrix3fv(loc, num, t, v)
}

func UniformMatrix4fv(loc, num int32, t bool, v *float32) {
	gl.UniformMatrix4fv(loc, num, t, v)
}

func GetUniformBlockIndex(program uint32, uniformBlockName *uint8) uint32 {
	return gl.GetUniformBlockIndex(program, uniformBlockName)
}

func UniformBlockBinding(program uint32, uniformBlockIndex uint32, uniformBlockBinding uint32) {
	gl.UniformBlockBinding(program, uniformBlockIndex, uniformBlockBinding)
}

// attribute

func EnableVertexAttribArray(index uint32) {
	gl.EnableVertexAttribArray(index)
}

func VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset int) {
	gl.VertexAttribPointer(index, size, xtype, normalized, stride, gl.PtrOffset(offset))
}

func VertexAttribIPointer(index uint32, size int32, xtype uint32, stride int32, offset int) {
	gl.VertexAttribIPointer(index, size, xtype, stride, gl.PtrOffset(offset))
}

func VertexAttribDivisor(index, divisor uint32) {
	gl.VertexAttribDivisor(index, divisor)
}

func DisableVertexAttribArray(index uint32) {
	gl.DisableVertexAttribArray(index)
}

func GetAttribLocation(program uint32, name string) int32 {
	return gl.GetAttribLocation(program, gl.Str(name))
}

func BindFragDataLocation(program uint32, color uint32, name string) {
	gl.BindFragDataLocation(program, color, gl.Str(name))
}

// texture

func ActiveTexture(texture uint32) {
	gl.ActiveTexture(texture)
}

func BindTexture(target uint32, texture uint32) {
	gl.BindTexture(target, texture)
}

func TexSubImage2D(target uint32, level int32, xOffset, yOffset, width, height int32, format, xtype uint32, pixels unsafe.Pointer) {
	gl.TexSubImage2D(target, level, xOffset, yOffset, width, height, format, xtype, pixels)
}

func TexImage2D(target uint32, level int32, internalFormat int32, width, height, border int32, format, xtype uint32, pixels unsafe.Pointer) {
	gl.TexImage2D(target, level, internalFormat, width, height, border, format, xtype, pixels)
}

func GenTextures(n int32, textures *uint32) {
	gl.GenTextures(n, textures)
}

func DeleteTextures(n int32, textures *uint32) {
	gl.DeleteTextures(n, textures)
}

func TexParameteri(texture, pname uint32, param int32) {
	gl.TexParameteri(texture, pname, param)
}

func TexParameterfv(target uint32, pname uint32, params *float32) {
	gl.TexParameterfv(target, pname, params)
}

func GenerateMipmap(target uint32) {
	gl.GenerateMipmap(target)
}

func BlitFramebuffer(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1 int32, mask, filter uint32) {
	gl.BlitFramebuffer(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1, mask, filter)
}
