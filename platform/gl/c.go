package gl

// #include <stdint.h>
// void* toPointer(intptr_t offset) {
//     return (void*)(offset);
// }
import "C"
import "unsafe"

func ptrOffset(offset int) unsafe.Pointer {
	return C.toPointer((C.intptr_t)(offset))
}
