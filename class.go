package main

import (
	"errors"

	"github.com/mrnavastar/assist/bytes"
)

const (
	JAVA_1   = 45
	JAVA_1_2 = 46
	JAVA_1_3 = 47
	JAVA_1_4 = 48
	JAVA_5   = 49
	JAVA_6   = 50
	JAVA_7   = 51
	JAVA_8   = 52
	JAVA_9   = 53
	JAVA_10  = 54
	JAVA_11  = 55
	JAVA_12  = 56
	JAVA_13  = 57
	JAVA_14  = 58
	JAVA_15  = 59
	JAVA_16  = 60
	JAVA_17  = 61
	JAVA_18  = 62
	JAVA_19  = 63
	JAVA_20  = 64
	JAVA_21  = 65
	JAVA_22  = 66
	JAVA_23  = 67
	JAVA_24  = 68

	CONSTANT_Class              = 7
	CONSTANT_Fieldref           = 9
	CONSTANT_Methodref          = 10
	CONSTANT_InterfaceMethodref = 11
	CONSTANT_String             = 8
	CONSTANT_Integer            = 3
	CONSTANT_Float              = 4
	CONSTANT_Long               = 5
	CONSTANT_Double             = 6
	CONSTANT_NameAndType        = 12
	CONSTANT_Utf8               = 1
	CONSTANT_MethodHandle       = 15
	CONSTANT_MethodType         = 16
	CONSTANT_Dynamic            = 17
	CONSTANT_InvokeDynamic      = 18
	CONSTANT_Module             = 19
	CONSTANT_Package            = 20

	ACC_PUBLIC     = 0x0001
	ACC_PRIVATE    = 0x0002
	ACC_PROTECTED  = 0x0004
	ACC_STATIC     = 0x0008
	ACC_FINAL      = 0x0010
	ACC_SUPER      = 0x0020
	ACC_VOLATILE   = 0x0040
	ACC_TRANSIENT  = 0x0080
	ACC_NATIVE     = 0x0100
	ACC_INTERFACE  = 0x0200
	ACC_ABSTRACT   = 0x0400
	ACC_STRICT     = 0x0800 //In a class file whose major version number is at least 46 and at most 60: Declared strictfp.
	ACC_SYNTHETIC  = 0x1000
	ACC_ANNOTATION = 0x2000
	ACC_ENUM       = 0x4000
	ACC_MODULE     = 0x8000
)

type InfoConstructor func() Info

var infoConstructors = map[byte]InfoConstructor{
	CONSTANT_Class:              func() Info { return &ClassInfo{} },
	CONSTANT_Fieldref:           func() Info { return &FieldRefInfo{} },
	CONSTANT_Methodref:          func() Info { return &MethodRefInfo{} },
	CONSTANT_InterfaceMethodref: func() Info { return &InterfaceMethodRefInfo{} },
	CONSTANT_String:             func() Info { return &StringInfo{} },
	CONSTANT_Integer:            func() Info { return &IntegerInfo{} },
	CONSTANT_Float:              func() Info { return &FloatInfo{} },
	CONSTANT_Long:               func() Info { return &LongInfo{} },
	CONSTANT_Double:             func() Info { return &DoubleInfo{} },
	CONSTANT_NameAndType:        func() Info { return &NameAndTypeInfo{} },
	CONSTANT_Utf8:               func() Info { return &Utf8Info{} },
	CONSTANT_MethodHandle:       func() Info { return &MethodHandleInfo{} },
	CONSTANT_MethodType:         func() Info { return &MethodTypeInfo{} },
	CONSTANT_Dynamic:            func() Info { return &DynamicInfo{} },
	CONSTANT_InvokeDynamic:      func() Info { return &InvokeDynamicInfo{} },
	CONSTANT_Module:             func() Info { return &ModuleInfo{} },
	CONSTANT_Package:            func() Info { return &PackageInfo{} },
}

var ErrNotClass = errors.New("jarhax: not a jvm class file")
var ErrInvalidClass = errors.New("jarhax: invalid class file")

type Info interface {
	Read(buf *bytes.Buffer)
	Write(buf *bytes.Buffer)
}

type ClassInfo struct {
	NameIndex uint16
}

func (info *ClassInfo) Read(buf *bytes.Buffer) {
	info.NameIndex = buf.ReadU16()
}

func (info *ClassInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_Class)
	buf.WriteU16(info.NameIndex)
}

type FieldRefInfo struct {
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (info *FieldRefInfo) Read(buf *bytes.Buffer) {
	info.ClassIndex = buf.ReadU16()
	info.NameAndTypeIndex = buf.ReadU16()
}

func (info *FieldRefInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_Fieldref)
	buf.WriteU16(info.ClassIndex)
	buf.WriteU16(info.NameAndTypeIndex)
}

type MethodRefInfo struct {
	FieldRefInfo
}

func (info *MethodRefInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_Methodref)
	buf.WriteU16(info.ClassIndex)
	buf.WriteU16(info.NameAndTypeIndex)
}

type InterfaceMethodRefInfo struct {
	FieldRefInfo
}

func (info InterfaceMethodRefInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_InterfaceMethodref)
	buf.WriteU16(info.ClassIndex)
	buf.WriteU16(info.NameAndTypeIndex)
}

type StringInfo struct {
	StringIndex uint16
}

func (info *StringInfo) Read(buf *bytes.Buffer) {
	info.StringIndex = buf.ReadU16()
}

func (info *StringInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_String)
	buf.WriteU16(info.StringIndex)
}

type IntegerInfo struct {
	Bytes uint32
}

func (info *IntegerInfo) Read(buf *bytes.Buffer) {
	info.Bytes = buf.ReadU32()
}

func (info *IntegerInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_Integer)
	buf.WriteU32(info.Bytes)
}

func (info *IntegerInfo) GetInt() int32 {
	return int32(info.Bytes)
}

type FloatInfo struct {
	IntegerInfo
}

func (info *FloatInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_Float)
	buf.WriteU32(info.Bytes)
}

type LongInfo struct {
	HighBytes uint32
	LowBytes  uint32
}

func (info *LongInfo) Read(buf *bytes.Buffer) {
	info.HighBytes = buf.ReadU32()
	info.LowBytes = buf.ReadU32()
}

func (info *LongInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_Long)
	buf.WriteU32(info.HighBytes)
	buf.WriteU32(info.LowBytes)
}

func (info *LongInfo) GetLong() int64 {
	return (int64(info.HighBytes) << int64(32)) + int64(info.LowBytes)
}

type DoubleInfo struct {
	LongInfo
}

func (info *DoubleInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_Double)
	buf.WriteU32(info.HighBytes)
	buf.WriteU32(info.LowBytes)
}

type NameAndTypeInfo struct {
	NameIndex       uint16
	DescriptorIndex uint16
}

func (info *NameAndTypeInfo) Read(buf *bytes.Buffer) {
	info.NameIndex = buf.ReadU16()
	info.DescriptorIndex = buf.ReadU16()
}

func (info *NameAndTypeInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_NameAndType)
	buf.WriteU16(info.NameIndex)
	buf.WriteU16(info.DescriptorIndex)
}

type Utf8Info struct {
	Length uint16
	Bytes  []byte
}

func (info *Utf8Info) Read(buf *bytes.Buffer) {
	info.Length = buf.ReadU16()
	info.Bytes = buf.ReadBytes(int(info.Length))
}

func (info *Utf8Info) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_Utf8)
	buf.WriteU16(info.Length)
	buf.Write(info.Bytes)
}

func (info *Utf8Info) Set(string string) {
	info.Length = uint16(len(string))
	info.Bytes = []byte(string)
}

// TODO: Implement FUll string spec
func (info Utf8Info) String() string {
	return string(info.Bytes)
}

type MethodHandleInfo struct {
	ReferenceKind  byte
	ReferenceIndex uint16
}

func (info *MethodHandleInfo) Read(buf *bytes.Buffer) {
	info.ReferenceKind = buf.ReadByte()
	info.ReferenceIndex = buf.ReadU16()
}

func (info *MethodHandleInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_MethodHandle)
	buf.WriteByte(info.ReferenceKind)
	buf.WriteU16(info.ReferenceIndex)
}

type MethodTypeInfo struct {
	DescriptorIndex uint16
}

func (info *MethodTypeInfo) Read(buf *bytes.Buffer) {
	info.DescriptorIndex = buf.ReadU16()
}

func (info *MethodTypeInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_MethodType)
	buf.WriteU16(info.DescriptorIndex)
}

type DynamicInfo struct {
	BootstrapMethodAttrIndex uint16
	NameAndTypeIndex         uint16
}

func (info *DynamicInfo) Read(buf *bytes.Buffer) {
	info.BootstrapMethodAttrIndex = buf.ReadU16()
	info.NameAndTypeIndex = buf.ReadU16()
}

func (info *DynamicInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_Dynamic)
	buf.WriteU16(info.BootstrapMethodAttrIndex)
	buf.WriteU16(info.NameAndTypeIndex)
}

type InvokeDynamicInfo struct {
	DynamicInfo
}

func (info *InvokeDynamicInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_InvokeDynamic)
	buf.WriteU16(info.BootstrapMethodAttrIndex)
	buf.WriteU16(info.NameAndTypeIndex)
}

type ModuleInfo struct {
	ClassInfo
}

func (info ModuleInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_Module)
	buf.WriteU16(info.NameIndex)
}

type PackageInfo struct {
	ClassInfo
}

func (info *PackageInfo) Write(buf *bytes.Buffer) {
	buf.WriteByte(CONSTANT_Package)
	buf.WriteU16(info.NameIndex)
}

type AttributeInfo struct {
	AttributeNameIndex uint16
	AttributeLength    uint32
	Data               []byte
}

func (info *AttributeInfo) Read(buf *bytes.Buffer) {
	info.AttributeNameIndex = buf.ReadU16()
	info.AttributeLength = buf.ReadU32()
	info.Data = buf.ReadBytes(int(info.AttributeLength))
}

func (info *AttributeInfo) Write(buf *bytes.Buffer) {
	buf.WriteU16(info.AttributeNameIndex)
	buf.WriteU32(info.AttributeLength)
	buf.Write(info.Data)
}

func ReadAttributes(buf *bytes.Buffer, count int) []AttributeInfo {
	var attributes []AttributeInfo
	for i := 0; i < count; i++ {
		var attribute AttributeInfo
		attribute.Read(buf)
		attributes = append(attributes, attribute)
	}
	return attributes
}

func WriteAttributes(buf *bytes.Buffer, attributes []AttributeInfo) {
	for _, attribute := range attributes {
		attribute.Write(buf)
	}
}

type FieldInfo struct {
	class           *Class
	AccessFlags     uint16
	NameIndex       uint16
	DescriptorIndex uint16
	AttributesCount uint16
	Attributes      []AttributeInfo
}

func (info *FieldInfo) Read(class *Class, buf *bytes.Buffer) {
	info.class = class
	info.AccessFlags = buf.ReadU16()
	info.NameIndex = buf.ReadU16()
	info.DescriptorIndex = buf.ReadU16()
	info.AttributesCount = buf.ReadU16()
	info.Attributes = ReadAttributes(buf, int(info.AttributesCount))
}

func (info *FieldInfo) Write(buf *bytes.Buffer) {
	buf.WriteU16(info.AccessFlags)
	buf.WriteU16(info.NameIndex)
	buf.WriteU16(info.DescriptorIndex)
	buf.WriteU16(info.AttributesCount)
	WriteAttributes(buf, info.Attributes)
}

func (info *FieldInfo) HasModifier(mod int) bool {
	return (info.AccessFlags & uint16(mod)) != 0
}

func (info *FieldInfo) GetName() string {
	return info.class.GetConstant(info.NameIndex).(*Utf8Info).String()
}

func (info *FieldInfo) GetDescriptor() string {
	return info.class.GetConstant(info.DescriptorIndex).(*Utf8Info).String()
}

type MethodInfo struct {
	FieldInfo
}

type Class struct {
	Magic             uint32
	MinorVersion      uint16
	MajorVersion      uint16
	ConstantPoolCount uint16
	ConstantPool      []Info
	AccessFlags       uint16
	ThisClass         uint16
	SuperClass        uint16
	InterfacesCount   uint16
	Interfaces        []uint16
	FieldsCount       uint16
	Fields            []FieldInfo
	MethodCount       uint16
	Methods           []MethodInfo
	AttributesCount   uint16
	Attributes        []AttributeInfo
}

func (class *Class) Read(b []byte) error {
	buf := bytes.Buffer{Data: &b, Index: 0}

	class.Magic = buf.ReadU32()
	if class.Magic != 0xCAFEBABE {
		return ErrNotClass
	}
	class.MinorVersion = buf.ReadU16()
	class.MajorVersion = buf.ReadU16()

	class.ConstantPoolCount = buf.ReadU16()
	for i := uint16(0); i < class.ConstantPoolCount-1; i++ {
		tag := buf.ReadByte()
		infoConstructor, ok := infoConstructors[tag]
		if !ok {
			return ErrInvalidClass
		}

		info := infoConstructor()
		info.Read(&buf)
		class.ConstantPool = append(class.ConstantPool, info)

		if tag == CONSTANT_Long || tag == CONSTANT_Double {
			i++ // Java specification: long and double take up two entries
		}
	}

	class.AccessFlags = buf.ReadU16()
	class.ThisClass = buf.ReadU16()
	class.SuperClass = buf.ReadU16()

	class.InterfacesCount = buf.ReadU16()
	for i := uint16(0); i < class.InterfacesCount; i++ {
		class.Interfaces = append(class.Interfaces, buf.ReadU16())
	}
	class.FieldsCount = buf.ReadU16()
	for i := uint16(0); i < class.FieldsCount; i++ {
		var fieldInfo FieldInfo
		fieldInfo.Read(class, &buf)
		class.Fields = append(class.Fields, fieldInfo)
	}
	class.MethodCount = buf.ReadU16()
	for i := uint16(0); i < class.MethodCount; i++ {
		var methodInfo MethodInfo
		methodInfo.Read(class, &buf)
		class.Methods = append(class.Methods, methodInfo)
	}

	class.AttributesCount = buf.ReadU16()
	class.Attributes = ReadAttributes(&buf, int(class.AttributesCount))
	return nil
}

func (class *Class) Write(b *[]byte) {
	buf := bytes.Buffer{Data: b, Index: 0}

	buf.WriteU32(class.Magic)
	buf.WriteU16(class.MinorVersion)
	buf.WriteU16(class.MajorVersion)

	buf.WriteU16(class.ConstantPoolCount)
	for _, constant := range class.ConstantPool {
		constant.(Info).Write(&buf)
	}

	buf.WriteU16(class.AccessFlags)
	buf.WriteU16(class.ThisClass)
	buf.WriteU16(class.SuperClass)

	buf.WriteU16(class.InterfacesCount)
	for _, i := range class.Interfaces {
		buf.WriteU16(i)
	}
	buf.WriteU16(class.FieldsCount)
	for _, field := range class.Fields {
		field.Write(&buf)
	}
	buf.WriteU16(class.MethodCount)
	for _, method := range class.Methods {
		method.Write(&buf)
	}

	buf.WriteU16(class.AttributesCount)
	WriteAttributes(&buf, class.Attributes)
}

func (class *Class) Supports(version int) bool {
	return class.MajorVersion >= uint16(version)
}

func (class *Class) GetConstant(index uint16) any {
	return class.ConstantPool[index-1]
}

func (class *Class) SetConstant(index uint16, constant Info) {
	class.ConstantPool[index-1] = constant
}

func (class *Class) GetClassName() string {
	return class.GetConstant(class.GetConstant(class.ThisClass).(ClassInfo).NameIndex).(Utf8Info).String()
}

func (class *Class) SetClassName(name string) {
	info := class.GetConstant(class.GetConstant(class.ThisClass).(ClassInfo).NameIndex).(Utf8Info)
	info.Set(name)
}

func (class *Class) GetSuperClassName() string {
	return class.GetConstant(class.GetConstant(class.SuperClass).(ClassInfo).NameIndex).(Utf8Info).String()
}

func (class *Class) SetSuperClassName(name string) {
	info := class.GetConstant(class.GetConstant(class.SuperClass).(ClassInfo).NameIndex).(Utf8Info)
	info.Set(name)
}

func (class *Class) GetInterfaceNames() []string {
	var interfaces []string
	for _, i := range class.Interfaces {
		interfaces = append(interfaces, class.GetConstant(class.GetConstant(i).(ClassInfo).NameIndex).(Utf8Info).String())
	}
	return interfaces
}

func (class *Class) HasModifier(mod int) bool {
	return (class.AccessFlags & uint16(mod)) != 0
}
 
func (class *Class) HasMainMethod() bool {
	for _, method := range class.Methods {
		if method.GetDescriptor() == "([Ljava/lang/String;)V" {
			return true
		}
	}
	return false
}
