#include <linux/module.h>
#include <linux/export-internal.h>
#include <linux/compiler.h>

MODULE_INFO(name, KBUILD_MODNAME);

__visible struct module __this_module
__section(".gnu.linkonce.this_module") = {
	.name = KBUILD_MODNAME,
	.init = init_module,
#ifdef CONFIG_MODULE_UNLOAD
	.exit = cleanup_module,
#endif
	.arch = MODULE_ARCH_INIT,
};



static const struct modversion_info ____versions[]
__used __section("__versions") = {
	{ 0x7b4041d1, "single_open" },
	{ 0x418b8bdc, "seq_putc" },
	{ 0x24d09486, "seq_write" },
	{ 0xeaee2b79, "seq_printf" },
	{ 0xbd03ed67, "__ref_stack_chk_guard" },
	{ 0xc7ffe1aa, "si_meminfo" },
	{ 0x151d4c65, "init_task" },
	{ 0x35283b01, "mmput" },
	{ 0x976c74cb, "get_task_mm" },
	{ 0xe838feb3, "access_process_vm" },
	{ 0x476167a9, "task_cputime_adjusted" },
	{ 0x97acb853, "ktime_get" },
	{ 0x90a48d82, "__ubsan_handle_out_of_bounds" },
	{ 0xd272d446, "__stack_chk_fail" },
	{ 0x46968a43, "proc_remove" },
	{ 0x41fd17a1, "seq_read" },
	{ 0x68140f95, "seq_lseek" },
	{ 0x1fde1a17, "single_release" },
	{ 0xd272d446, "__fentry__" },
	{ 0x6b01a33d, "proc_create" },
	{ 0xe8213e80, "_printk" },
	{ 0xd272d446, "__x86_return_thunk" },
	{ 0xd954c786, "module_layout" },
};

static const u32 ____version_ext_crcs[]
__used __section("__version_ext_crcs") = {
	0x7b4041d1,
	0x418b8bdc,
	0x24d09486,
	0xeaee2b79,
	0xbd03ed67,
	0xc7ffe1aa,
	0x151d4c65,
	0x35283b01,
	0x976c74cb,
	0xe838feb3,
	0x476167a9,
	0x97acb853,
	0x90a48d82,
	0xd272d446,
	0x46968a43,
	0x41fd17a1,
	0x68140f95,
	0x1fde1a17,
	0xd272d446,
	0x6b01a33d,
	0xe8213e80,
	0xd272d446,
	0xd954c786,
};
static const char ____version_ext_names[]
__used __section("__version_ext_names") =
	"single_open\0"
	"seq_putc\0"
	"seq_write\0"
	"seq_printf\0"
	"__ref_stack_chk_guard\0"
	"si_meminfo\0"
	"init_task\0"
	"mmput\0"
	"get_task_mm\0"
	"access_process_vm\0"
	"task_cputime_adjusted\0"
	"ktime_get\0"
	"__ubsan_handle_out_of_bounds\0"
	"__stack_chk_fail\0"
	"proc_remove\0"
	"seq_read\0"
	"seq_lseek\0"
	"single_release\0"
	"__fentry__\0"
	"proc_create\0"
	"_printk\0"
	"__x86_return_thunk\0"
	"module_layout\0"
;

MODULE_INFO(depends, "");


MODULE_INFO(srcversion, "6FD5EBE318BE8A34B0DF5A5");
