package babe

import (
	"bytes"
	"encoding/binary"
	"errors"
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

var ErrNotClass = errors.New("not a jvm class file")
var ErrInvalidClass = errors.New("invalid class file")

type ConstantInfo interface {}

type ClassInfo struct {
	ConstantInfo
	Tag		  byte
	NameIndex uint16
}

type FieldRefInfo struct {
	ConstantInfo
	Tag				 byte
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

type MethodRefInfo struct {
	FieldRefInfo
}

type InterfaceMethodRefInfo struct {
	FieldRefInfo
}

type StringInfo struct {
	ConstantInfo
	Tag			byte
	StringIndex uint16
}

type IntegerInfo struct {
	ConstantInfo
	Tag	  byte
	Bytes uint32
}

func (info *IntegerInfo) GetInt() int32 {
	return int32(info.Bytes)
}

type FloatInfo struct {
	IntegerInfo
}

type LongInfo struct {
	ConstantInfo
	Tag		  byte
	HighBytes uint32
	LowBytes  uint32
}

func (info *LongInfo) GetLong() int64 {
	return (int64(info.HighBytes) << int64(32)) + int64(info.LowBytes)
}

type DoubleInfo struct {
	LongInfo
}

type NameAndTypeInfo struct {
	ConstantInfo
	Tag				byte
	NameIndex       uint16
	DescriptorIndex uint16
}

type Utf8Info struct {
	ConstantInfo
	Tag	   byte
	Length uint16
	Bytes  []byte
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
	ConstantInfo
	Tag			   byte
	ReferenceKind  byte
	ReferenceIndex uint16
}

type MethodTypeInfo struct {
	ConstantInfo
	Tag				byte
	DescriptorIndex uint16
}

type DynamicInfo struct {
	ConstantInfo
	Tag						 byte
	BootstrapMethodAttrIndex uint16
	NameAndTypeIndex         uint16
}

type InvokeDynamicInfo struct {
	DynamicInfo
}

type ModuleInfo struct {
	ClassInfo
}

type PackageInfo struct {
	ClassInfo
}

type AttributeInfo struct {
	AttributeNameIndex uint16
	AttributeLength    uint32
	Data               []byte
}

type FieldInfo struct {
	class           *Class
	AccessFlags     uint16
	NameIndex       uint16
	DescriptorIndex uint16
	AttributesCount uint16
	Attributes      []AttributeInfo
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
	ConstantPool      []ConstantInfo
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

func (class *Class) Read(b *bytes.Buffer) error {
	err := binary.Read(b, binary.BigEndian, &class)
	if err != nil {
		return err
	}
	if class.Magic != 0xCAFEBABE {
		return ErrNotClass
	}
	return nil
}

func (class *Class) Write(b *bytes.Buffer) error {
	return binary.Write(b, binary.BigEndian, class)
}

func (class *Class) Supports(version int) bool {
	return class.MajorVersion >= uint16(version)
}

func (class *Class) GetConstant(index uint16) any {
	return class.ConstantPool[index-1]
}

func (class *Class) SetConstant(index uint16, constant ConstantInfo) {
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
