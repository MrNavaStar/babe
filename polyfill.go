package main

func Polyfill(class *Class, version int) {
	if class.Supports(version) {
		return
	}
	class.MajorVersion = uint16(version)

	for i, constant := range class.ConstantPool {
		switch info := constant.(type) {
		case *MethodHandleInfo:
			if !class.Supports(JAVA_7) {

				class.ConstantPool[i] = &Utf8Info{1, []byte{0}}
			}
		case *MethodTypeInfo:
			if !class.Supports(JAVA_7) {
				class.ConstantPool[i] = &ClassInfo{info.DescriptorIndex}
			}
		case *InvokeDynamicInfo:
			if !class.Supports(JAVA_7) {
				class.ConstantPool[i] = &FieldRefInfo{info.BootstrapMethodAttrIndex, info.NameAndTypeIndex}
			}
		case *ModuleInfo:
			if !class.Supports(JAVA_9) {
				class.ConstantPool[i] = &ClassInfo{info.NameIndex}
			}
		case *PackageInfo:
			if !class.Supports(JAVA_9) {
				class.ConstantPool[i] = &ClassInfo{info.NameIndex}
			}
		case *DynamicInfo:
			if !class.Supports(JAVA_11) {
				class.ConstantPool[i] = &FieldRefInfo{info.BootstrapMethodAttrIndex, info.NameAndTypeIndex}
			}
		}
	}
}
